package redis

import (
	"testing"
	"time"
)

type LocalLog struct {
	t *testing.T
}

func (l LocalLog) Log(kv ...interface{}) {
	l.t.Log(kv...)
}

const (
	redisKey string = "redis.key"
)

func TestNewRedisClient(t *testing.T) {
	testLog := &LocalLog{t: t}
	redisCli := NewRedisClient(testLog, WithRdb("127.0.0.1", 16379))
	redisCli.SetNxPx(redisKey, "1", time.Second)

	res, err := redisCli.Get(redisKey)
	if err != nil {
		t.Log("get key err: ", err)
		t.FailNow()
	}

	if res != "1" {
		t.Log("store error, wat 1, actual", res)
		t.FailNow()
	}

	err = redisCli.Delete(redisKey)
	if err != nil {
		t.Log("delete key err: ", err)
		t.FailNow()
	}

	r, err := redisCli.Exist(redisKey)
	if err != nil {
		t.Log("get key err: ", err)
		t.FailNow()
	}

	if r != 0 {
		t.Log("store error, wat 0, actual", r)
		t.FailNow()
	}

}
