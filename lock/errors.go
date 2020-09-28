package lock

import "fmt"

var (
	LockTypeUnknown error = fmt.Errorf("unknown type lock, check param LockType")
)
