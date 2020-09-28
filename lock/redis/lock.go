package redis

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/dingkegithub/distrubutelock/lock"
	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
	"github.com/dingkegithub/distrubutelock/utils/log"
)

const (
	lockPrefix string = "lock.redis"
)

//
// RdLock: redis distribute lock
//
type RdLock struct {
	// lock key
	lockKey string

	// rent time
	ttl time.Duration

	// redis client
	cli *RedisClient

	cliOpt *RedisOptions

	logger log.Logging

	// control keeplive routine
	ttlChan chan struct{}

	lockVal string

	once sync.Once
}

func NewRdLock(key string) *RdLock {
	opts := defaultRedisOptions()

	return &RdLock{
		cliOpt:  opts,
		logger:  opts.logger,
		ttl:     time.Second,
		ttlChan: make(chan struct{}),
		lockKey: fmt.Sprintf("%s.%s", lockPrefix, key),
		lockVal: fmt.Sprintf("%s.%d", lockPrefix, rand.Int63()),
	}
}

func (l *RdLock) SetOption(opts ...commoninterface.ParamOption) {
	l.cliOpt.Apply(opts...)
	l.logger = l.cliOpt.logger
}

func (l *RdLock) Close() {
	if l.cli != nil {
		l.cli.Close()
	}
}

// Lock:
// exec commond "SET resource_name my_random_value NX PX 30000"
// available option:
// Block: true, wait until lockKey deleted and SETEX success
//        false, Timeout=0 immediately return else see Timeout
// Timeout: wait lock time. 0 means no need wait
func (l *RdLock) Lock(opts ...lock.Option) (bool, error) {
	l.once.Do(func() {
		if l.cli == nil {
			l.cli = NewRedisClient(l.cliOpt)
			l.cli.Open()
		}
	})

	o := &lock.LockOptions{}

	o.Apply(opts...)

	if l.cli == nil {
		return false, nil
	}

	tm := time.Now().Add(o.Timeout)

	for {
		ok, err := l.cli.SetNxPx(l.lockKey, l.lockVal, l.ttl)
		if err != nil {
			return false, err
		}

		if !ok {
			if o.Block || (o.Tmoutable && time.Now().Before(tm)) {
				time.Sleep(5 * time.Millisecond)
				continue
			} else {
				return false, nil
			}
		}
		go l.keeplive()
		return true, nil
	}
	return false, nil
}

// UnLock, release lock
// delete lockKey and stop keeplive routine
func (l *RdLock) UnLock() (bool, error) {
	l.ttlChan <- struct{}{}
	err := l.cli.Delete(l.lockKey)
	if err != nil {
		return false, nil
	}
	<-l.ttlChan
	return true, nil
}

// keeplive: rent lockKey by set expire
func (l *RdLock) keeplive() {
	circle := l.ttl / 3
	interval := time.Tick(circle)
	for {
		select {
		case <-interval:
			l.cli.Expire(l.lockKey, l.ttl)

		case <-l.ttlChan:
			l.ttlChan <- struct{}{}
			break
		}
	}
}
