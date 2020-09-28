package etcdv2

import "fmt"

var (
	ErrorParseParam     = fmt.Errorf("parse param error")
	ErrorCommit         = fmt.Errorf("commit request failed after retry")
	ErrorExist          = fmt.Errorf("key has exist")
	ErrorUnknownEtcdMsg = fmt.Errorf("unknown etcd msg")
)
