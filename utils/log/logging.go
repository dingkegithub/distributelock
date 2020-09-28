package log

import "fmt"

type Logging interface {
	Log(kvargs ...interface{})
}

type DefaultLogging struct {
}

func (DefaultLogging) Log(kvargs ...interface{}) {
	fmt.Println(kvargs...)
}
