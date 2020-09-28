package redis

import (
	"time"

	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
	"github.com/dingkegithub/distrubutelock/utils/log"
)

type Mode int

func (m Mode) String() string {
	switch m {
	case ModeRdb:
		return "rdb.mode"
	case ModeCluster:
		return "cluster.mode"
	default:
		return "unknown.mode"
	}
}

const (
	ModeRdb = 1 << iota
	ModeCluster
)

type DbOption struct {
	RedisHost string
	RedisPort uint16
	UserName  string
	UserAuth  string
	PoolSize  int
	WTimeout  time.Duration
	RTimeout  time.Duration
}

type ClusterOption struct {
	Addrs []string
}

type RedisOptions struct {
	mode    Mode
	rdb     *DbOption
	cluster *ClusterOption
	logger  log.Logging
}

func defaultRedisOptions() *RedisOptions {
	opts := &RedisOptions{}
	opts.mode = ModeRdb
	opts.rdb = &DbOption{
		RedisHost: "127.0.0.1",
		RedisPort: 6379,
		PoolSize:  10,
		RTimeout:  3 * time.Second,
		WTimeout:  3 * time.Second,
	}

	return opts
}

func (opt *RedisOptions) Apply(opts ...commoninterface.ParamOption) {
	for _, o := range opts {
		o(opt)
	}
}

func WithMode(m Mode) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.mode = m
	}
}

func WithRdb(host string, port uint16) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.rdb.RedisHost = host
		rdsOptions.rdb.RedisPort = port
	}
}

func WithRdbAuth(uname string, auth string) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.rdb.UserName = uname
		rdsOptions.rdb.UserAuth = auth
	}
}

func WithPoolSize(size int) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.rdb.PoolSize = size
	}
}

func WithRdbTimeout(wt uint64, rt uint64) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.rdb.WTimeout = time.Duration(wt) * time.Millisecond
		rdsOptions.rdb.RTimeout = time.Duration(rt) * time.Millisecond
	}
}

func WithClusterAddrs(addrs []string) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.cluster.Addrs = make([]string, len(addrs))
		copy(rdsOptions.cluster.Addrs, addrs)
	}
}

func WithLogger(logger log.Logging) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		rdsOptions := options.(*RedisOptions)
		rdsOptions.logger = logger
	}
}
