package etcdv2

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/dingkegithub/distributelock/lock"
	"github.com/dingkegithub/distributelock/pkg/commoninterface"
	"github.com/dingkegithub/distributelock/utils/log"
)

const (
	LockPrefix string = "lock.etcd"
)

//
// RdLock: redis distribute lock
//
type EtcdLock struct {
	// lock key
	lockKey string

	// lease time, second
	ttl uint64

	// etcd client
	cli *EtcdClient

	// control keeplive routine
	ttlChan chan struct{}

	lockVal string

	cliOpts *EtcdOptions

	logger log.Logging

	once sync.Once
}

func NewEtcdLock(key string) *EtcdLock {
	opts := defaultClientOptions()
	return &EtcdLock{
		ttl:     5,
		cliOpts: opts,
		logger:  opts.logger,
		ttlChan: make(chan struct{}),
		lockKey: fmt.Sprintf("%s.%s", LockPrefix, key),
		lockVal: fmt.Sprintf("%s.%d", LockPrefix, rand.Int63()),
	}
}

func (l *EtcdLock) SetOption(opts ...commoninterface.ParamOption) {
	l.cliOpts.Apply(opts...)
	l.logger = l.cliOpts.logger
}

// Lock:
// exec commond "SET resource_name my_random_value NX PX 30000"
// available option:
// Block: true, wait until lockKey deleted and SETEX success
//        false, Timeout=0 immediately return else see Timeout
// Timeout: wait lock time. 0 means no need wait
func (l *EtcdLock) Lock(opts ...lock.Option) (resOk bool, resErr error) {
	defer func() {
		if r := recover(); r != nil {
			resOk, resErr = false, fmt.Errorf("%v", r)
		}
	}()

	l.once.Do(func() {
		cli, err := NewEtcdClient(l.cliOpts)
		if err != nil {
			panic(err)
		}
		l.cli = cli
	})

	o := &lock.LockOptions{}
	o.Apply(opts...)

	if l.cli == nil {
		return false, nil
	}

	tm := time.Now().Add(o.Timeout)

	for {
		_, err := l.cli.AtomicCreate(l.lockKey, l.lockVal, l.ttl)
		if err != nil {
			if err == ErrorExist {
				if o.Block || (o.Tmoutable && time.Now().Before(tm)) {
					time.Sleep(10 * time.Millisecond)
					continue
				}
				return false, nil
			} else {
				return false, err
			}
		}
		go l.keeplive()
		return true, nil
	}
}

// UnLock, release lock
// delete lockKey and stop keeplive routine
func (l *EtcdLock) UnLock() (bool, error) {
	l.ttlChan <- struct{}{}
	_, err := l.cli.AtomicDelete(l.lockKey, l.lockVal)
	if err != nil {
		fmt.Println("unlock failed with error: ", err)
		return false, nil
	}
	<-l.ttlChan
	return true, nil
}

// keeplive: rent lockKey by set expire
func (l *EtcdLock) keeplive() {
	circle := l.ttl / 3
	interval := time.Tick(time.Duration(circle))
	for {
		select {
		case <-interval:
			l.cli.AtomicLease(l.lockKey, l.lockVal, l.ttl)

		case <-l.ttlChan:
			l.ttlChan <- struct{}{}
			break
		}
	}
}
