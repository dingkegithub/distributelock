package zk

import (
	"testing"
	"time"

	"github.com/dingkegithub/distrubutelock/lock"
)

func TestNewZkLock1(t *testing.T) {

	zkLock := NewZkLock("mylock")

	start := time.Now()
	ok, err := zkLock.Lock()
	elapse := time.Since(start)

	t.Log("get lock: ", ok, " elapse:", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}
}

func TestNewZkLock2(t *testing.T) {

	zkLock := NewZkLock("mylock")

	start := time.Now()
	ok, err := zkLock.Lock(lock.WithBlock())
	elapse := time.Since(start)

	t.Log("get lock: ", ok, " elapse:", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}
}

func TestNewZkLock3(t *testing.T) {

	zkLock := NewZkLock("mylock")

	start := time.Now()
	ok, err := zkLock.Lock(lock.WithTimeout(20000))
	elapse := time.Since(start)

	t.Log("get lock: ", ok, " elapse:", elapse)

	if err != nil {
		t.Log("lock error: ", err)
		t.FailNow()
	}
}
