package distributelock

import (
	"github.com/dingkegithub/distrubutelock/lock"
	"github.com/dingkegithub/distrubutelock/lock/etcdv2"
	"github.com/dingkegithub/distrubutelock/lock/redis"
	"github.com/dingkegithub/distrubutelock/lock/zk"
	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
)

const (
	LockTypeRedis = 1 << iota

	LockTypeEtcd

	LockTypeZk
)

type LockType int

func (lt LockType) String() string {
	switch lt {
	case LockTypeRedis:
		return "RedisLock"
	case LockTypeEtcd:
		return "EtcdLock"
	case LockTypeZk:
		return "ZkLock"
	default:
		return "UnknownLock"
	}
}

type LockClient struct {
	proxy lock.Locker
}

func NewLockClient(t LockType, lockName string, options ...commoninterface.ParamOption) *LockClient {
	switch t {
	case LockTypeRedis:
		rdCli := redis.NewRdLock(lockName)
		rdCli.SetOption(options...)
		return &LockClient{
			proxy: rdCli,
		}

	case LockTypeEtcd:
		etcdCli := etcdv2.NewEtcdLock(lockName)
		etcdCli.SetOption(options...)
		return &LockClient{
			proxy: etcdCli,
		}

	case LockTypeZk:
		zkCli := zk.NewZkLock(lockName)
		zkCli.SetOption(options...)
		return &LockClient{
			proxy: zkCli,
		}

	default:
		return nil
	}
}

func (l *LockClient) Lock(o ...lock.Option) (bool, error) {
	return l.proxy.Lock(o...)
}

func (l *LockClient) UnLock() (bool, error) {
	return l.proxy.UnLock()
}
