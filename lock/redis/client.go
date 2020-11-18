package redis

import (
	"fmt"
	"time"

	"github.com/dingkegithub/distributelock/utils/log"
	goredis "github.com/go-redis/redis"
)

type RedisClient struct {
	rds     *goredis.Client
	opts    *RedisOptions
	logger  log.Logging
	cluster *goredis.ClusterClient
}

func NewRedisClient(opts *RedisOptions) *RedisClient {
	return &RedisClient{
		opts:   opts,
		logger: opts.logger,
	}
}

func (rc *RedisClient) Close() {
	if rc.rds != nil {
		rc.rds.Close()
		rc.rds = nil
	}

	if rc.cluster != nil {
		rc.cluster.Close()
	}
}

func (rc *RedisClient) Open() {
	if rc.opts.mode == ModeRdb {
		rdbCli := goredis.NewClient(&goredis.Options{
			Addr:         fmt.Sprintf("%s:%d", rc.opts.rdb.RedisHost, rc.opts.rdb.RedisPort),
			PoolSize:     rc.opts.rdb.PoolSize,
			ReadTimeout:  rc.opts.rdb.RTimeout,
			WriteTimeout: rc.opts.rdb.WTimeout,
			Password:     rc.opts.rdb.UserAuth,
		})

		resp, err := rdbCli.Ping().Result()
		if err != nil {
			rc.logger.Log("file", "client.go", "func", "NewRedisClient", "msg", "ping redis server error", "error", err)
			panic(err)
		}

		if resp != "PONG" {
			rc.logger.Log("file", "client.go", "func", "NewRedisClient", "msg", "ping redis server", "resp", resp)
		}

		rc.rds = rdbCli

	} else if rc.opts.mode == ModeCluster {

		cluster := goredis.NewClusterClient(&goredis.ClusterOptions{
			Addrs: rc.opts.cluster.Addrs,
		})

		resp, err := cluster.Ping().Result()
		if err != nil {
			rc.logger.Log("file", "client.go", "func", "NewRedisClient", "msg", "ping redis cluster error", "error", err)
			panic(err)
		}

		if resp != "PONG" {
			rc.logger.Log("file", "client.go", "func", "NewRedisClient", "msg", "ping redis cluster", "resp", resp)
		}
		rc.cluster = cluster

	} else {
		rc.logger.Log("file", "client.go", "func", "NewRedisClient", "msg", "parameter unknown")
		panic(ErrParaUnknownRedisMode)
	}
}

//
// Commond: SETNX key value
//
func (rc *RedisClient) SetNxPx(key, value string, expire time.Duration) (bool, error) {
	if rc.opts.mode == ModeRdb {
		return rc.rds.SetNX(key, value, expire).Result()
	} else if rc.opts.mode == ModeCluster {
		return rc.cluster.SetNX(key, value, expire).Result()
	} else {
		return false, ErrParaUnknownRedisMode
	}
}

func (rc *RedisClient) Get(key string) (string, error) {
	if rc.opts.mode == ModeRdb {
		return rc.rds.Get(key).Result()
	} else if rc.opts.mode == ModeCluster {
		return rc.cluster.Get(key).Result()
	} else {
		return "", ErrParaUnknownRedisMode
	}
}

func (rc *RedisClient) Expire(key string, expire time.Duration) (bool, error) {
	if rc.opts.mode == ModeRdb {
		return rc.rds.Expire(key, expire).Result()
	} else if rc.opts.mode == ModeCluster {
		return rc.cluster.Expire(key, expire).Result()
	} else {
		return false, ErrParaUnknownRedisMode
	}
}

func (rc *RedisClient) Exist(key string) (int64, error) {
	if rc.opts.mode == ModeRdb {
		return rc.rds.Exists(key).Result()
	} else if rc.opts.mode == ModeCluster {
		return rc.cluster.Exists(key).Result()
	} else {
		return 0, ErrParaUnknownRedisMode
	}
}

func (rc *RedisClient) Delete(key string) error {
	if rc.opts.mode == ModeRdb {
		return rc.rds.Del(key).Err()
	} else if rc.opts.mode == ModeCluster {
		return rc.cluster.Del(key).Err()
	} else {
		return ErrParaUnknownRedisMode
	}
}
