package redis

import "fmt"

var (
	ErrParaUnknownRedisMode error = fmt.Errorf("unknow redis mode")
)
