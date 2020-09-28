package etcdv2

import (
	"testing"
)

func TestNewEtcdClient(t *testing.T) {

	cli, err := NewEtcdClient(WithAddr([]string{"127.0.0.1:2379"}))
	if err != nil {
		t.Log("new client error: ", err)
		t.FailNow()
	}

	_, err = cli.AtomicCreate("you", "are", 60)
	if err != nil {
		t.Log("atomic create error: ", err)
		t.FailNow()
	}
}
