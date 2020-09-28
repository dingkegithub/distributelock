package zk

import (
	"fmt"
	"sync"

	"github.com/dingkegithub/distrubutelock/lock"
	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
	"github.com/dingkegithub/distrubutelock/utils/log"
)

const (
	lockPrefix string = "lock.zookeeper"
)

type ZkLock struct {
	// lock key
	lockKey string

	// zk client
	cli *ZkClient

	lockVal string

	//
	opts *zkOptions

	logger log.Logging

	curEmpheralNode string

	once sync.Once
}

func NewZkLock(key string) *ZkLock {
	opts := defaultzkOptions()

	return &ZkLock{
		opts:    opts,
		logger:  opts.logger,
		lockKey: fmt.Sprintf("/%s.%s", lockPrefix, key),
		lockVal: fmt.Sprintf("%s", lockPrefix),
	}
}

func (l *ZkLock) SetOption(opts ...commoninterface.ParamOption) {
	l.opts.Apply(opts...)
	l.logger = l.opts.logger
}

func (l *ZkLock) Lock(opts ...lock.Option) (bool, error) {

	l.once.Do(func() {
		if l.cli == nil {
			l.cli = NewZkClient(l.opts)
		}
	})

	o := &lock.LockOptions{}

	if len(opts) > 0 {
		o.Apply(opts...)
	}

	err := l.cli.Open(3)
	if err != nil {
		l.logger.Log("file", "lock.go", "func", "Lock", "msg", "open zk client error", "error", err)
		return false, err
	}

	ok, name, err := l.cli.CreateFirstEphemeral(l.lockKey, l.lockVal, o.Block, o.Tmoutable, o.Timeout)
	if err != nil {
		l.cli.DeleteEphemeralNode(l.lockKey, name)
		l.cli.Close()
		return false, err
	}

	if !ok {
		l.cli.DeleteEphemeralNode(l.lockKey, name)
		l.cli.Close()
		return false, nil
	}

	l.curEmpheralNode = name
	return true, nil
}

// UnLock, release lock
// delete lockKey and stop keeplive routine
func (l *ZkLock) UnLock() (bool, error) {
	l.cli.DeleteEphemeralNode(l.lockKey, l.curEmpheralNode)
	l.cli.Close()
	return true, nil
}
