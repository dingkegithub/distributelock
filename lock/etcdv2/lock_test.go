package etcdv2

import (
	"fmt"
	"testing"
	"time"

	"github.com/dingkegithub/distrubutelock/lock"
)

func TestNewEtcdLock1(t *testing.T) {
	etcdCli, err := NewEtcdClient(WithAddr([]string{"127.0.0.1:2379"}))
	if err != nil {
		t.Log("new client: ", err)
		t.FailNow()
	}

	if etcdCli == nil {
		t.Log("create etcd client")
		t.FailNow()
	}

	l := NewEtcdLock("etcdlocker").SetClient(etcdCli)

	s := time.Now()
	ok, err := l.Lock()
	elapse := time.Since(s)
	t.Log("lock elapse: ", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}

	if !ok {
		t.Log("lock failed")
		t.FailNow()
	}
}

func TestNewEtcdLock2(t *testing.T) {
	fmt.Println()
	etcdCli, err := NewEtcdClient(WithAddr([]string{"127.0.0.1:2379"}))
	if err != nil {
		t.Log("new client: ", err)
		t.FailNow()
	}

	if etcdCli == nil {
		t.Log("create etcd client")
		t.FailNow()
	}

	l := NewEtcdLock("etcdlocker").SetClient(etcdCli)

	s := time.Now()
	ok, err := l.Lock()
	elapse := time.Since(s)
	t.Log("lock elapse: ", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}

	if !ok {
		t.Log("lock failed")
		t.FailNow()
	}

	ok, err = l.UnLock()

	if err != nil {
		t.Log("unlock error: ", err)
		t.FailNow()
	}

	if !ok {
		t.Log("unlock failed")
		t.FailNow()
	}
}

func TestNewEtcdLock3(t *testing.T) {
	etcdCli, err := NewEtcdClient(WithAddr([]string{"127.0.0.1:2379"}))
	if err != nil {
		t.Log("new client: ", err)
		t.FailNow()
	}

	if etcdCli == nil {
		t.Log("create etcd client")
		t.FailNow()
	}

	l := NewEtcdLock("etcdlocker").SetClient(etcdCli)

	s := time.Now()
	ok, err := l.Lock(lock.WithTimeout(10000))
	elapse := time.Since(s)
	t.Log("lock elapse: ", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}

	if !ok {
		t.Log("lock failed")
		t.FailNow()
	}
}
