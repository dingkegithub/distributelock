package distributelock

import (
	"testing"
	"time"

	"github.com/dingkegithub/distributelock/lock"
	"github.com/dingkegithub/distributelock/lock/etcdv2"
	"github.com/dingkegithub/distributelock/lock/redis"
)

type tlog struct {
	t *testing.T
}

func (t *tlog) Log(kvargs ...interface{}) {
	t.t.Log(kvargs...)
}

func TestLockFactory1_NonBlock(t *testing.T) {

	dlock := NewLockClient(LockTypeRedis, "mylock", redis.WithMode(redis.ModeRdb), redis.WithRdb("127.0.0.1", 16379))

	s := time.Now()
	ok, err := dlock.Lock()
	elapse := time.Since(s)
	t.Log("NonBlock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactor1_Block(t *testing.T) {

	dlock := NewLockClient(LockTypeRedis, "mylock", redis.WithMode(redis.ModeRdb), redis.WithRdb("127.0.0.1", 16379))

	s := time.Now()
	ok, err := dlock.Lock(lock.WithBlock())
	elapse := time.Since(s)
	t.Log("BlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactor1_Timeout(t *testing.T) {

	dlock := NewLockClient(LockTypeRedis, "mylock", redis.WithMode(redis.ModeRdb), redis.WithRdb("127.0.0.1", 16379))

	s := time.Now()
	ok, err := dlock.Lock(lock.WithTimeout(30000))
	elapse := time.Since(s)
	t.Log("TimeoutLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}
func TestLockFactory2_NonBlock(t *testing.T) {
	dlock := NewLockClient(LockTypeEtcd, "mylock", etcdv2.WithAddr([]string{"127.0.0.1:2379"}))
	s := time.Now()
	ok, err := dlock.Lock()
	elapse := time.Since(s)
	t.Log("EtcdNonBlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactory2_Block(t *testing.T) {
	dlock := NewLockClient(LockTypeEtcd, "mylock", etcdv2.WithAddr([]string{"127.0.0.1:2379"}))
	s := time.Now()
	ok, err := dlock.Lock(lock.WithBlock())
	elapse := time.Since(s)
	t.Log("EtcdBlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactory2_Timeout(t *testing.T) {
	dlock := NewLockClient(LockTypeEtcd, "mylock", etcdv2.WithAddr([]string{"127.0.0.1:2379"}))
	s := time.Now()
	ok, err := dlock.Lock(lock.WithTimeout(20000))
	elapse := time.Since(s)
	t.Log("EtcdTimeoutLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactory3_NonBlock(t *testing.T) {
	dlock := NewLockClient(LockTypeZk, "mylock")
	s := time.Now()
	ok, err := dlock.Lock()
	elapse := time.Since(s)
	t.Log("RdsNonBlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactory3_Block(t *testing.T) {
	dlock := NewLockClient(LockTypeZk, "mylock")
	s := time.Now()
	ok, err := dlock.Lock(lock.WithBlock())
	elapse := time.Since(s)
	t.Log("RdsNonBlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}

func TestLockFactory3_Timeout(t *testing.T) {
	dlock := NewLockClient(LockTypeZk, "mylock")
	s := time.Now()
	ok, err := dlock.Lock(lock.WithTimeout(20000))
	elapse := time.Since(s)
	t.Log("RdsNonBlockLock elapse: ", elapse, " ok:", ok, " err:", err)
	//	dlock.UnLock()
}
