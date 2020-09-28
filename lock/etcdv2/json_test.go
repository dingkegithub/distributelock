package etcdv2

import (
	"encoding/json"
	"testing"
)

func TestStruct(t *testing.T) {
	create := &AtomicCreateResp{}
	s := "{\"action\":\"update\",\"node\":{\"key\":\"/foo\",\"value\":\"\",\"expiration\":\"2020-09-23T08:32:05.0389929Z\",\"ttl\":30,\"modifiedIndex\":15,\"createdIndex\":7},\"prevNode\":{\"key\":\"/foo\",\"value\":\"\",\"expiration\":\"2020-09-23T08:32:03.8015632Z\",\"ttl\":29,\"modifiedIndex\":14,\"createdIndex\":7}}"
	s2 := "{\"action\":\"create\",\"node\":{\"key\":\"/foo\",\"value\":\"three\",\"expiration\":\"2020-09-23T08:31:12.3119462Z\",\"ttl\":30,\"modifiedIndex\":7,\"createdIndex\":7}}"
	s3 := " {\"action\":\"compareAndDelete\",\"node\":{\"key\":\"/foo\",\"modifiedIndex\":4,\"createdIndex\":3},\"prevNode\":{\"key\":\"/foo\",\"value\":\"one\",\"modifiedIndex\":3,\"createdIndex\":3}}"
	s4 := "{\"errorCode\":105,\"message\":\"Key already exists\",\"cause\":\"/etcd.lock:etcdlock\",\"index\":88}"
	err := json.Unmarshal([]byte(s), &create)
	if err != nil {
		t.FailNow()
	}

	if create.Action != "update" {
		t.FailNow()
	}

	err = json.Unmarshal([]byte(s2), &create)
	if err != nil {
		t.FailNow()
	}

	err = json.Unmarshal([]byte(s3), &create)
	if err != nil {
		t.FailNow()
	}

	var create4 *AtomicCreateResp
	err = json.Unmarshal([]byte(s4), &create4)
	if err != nil {
		t.FailNow()
	}
}
