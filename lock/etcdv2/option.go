package etcdv2

import (
	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
	"github.com/dingkegithub/distrubutelock/utils/log"
)

type EtcdOptions struct {
	retry         int
	heartInterval uint64
	Addrs         []string
	logger        log.Logging
}

func (e *EtcdOptions) Apply(opts ...commoninterface.ParamOption) {
	for _, o := range opts {
		o(e)
	}
}

func defaultClientOptions() *EtcdOptions {
	return &EtcdOptions{
		Addrs:         []string{"127.0.0.1:2380"},
		logger:        &log.DefaultLogging{},
		heartInterval: 1,
		retry:         3,
	}
}

func WithLogger(logger log.Logging) commoninterface.ParamOption {
	return func(p commoninterface.ParamOptions) {
		e := p.(*EtcdOptions)
		e.logger = logger
	}
}

func WithAddr(addrs []string) commoninterface.ParamOption {
	return func(p commoninterface.ParamOptions) {
		e := p.(*EtcdOptions)
		e.Addrs = make([]string, len(addrs))
		copy(e.Addrs, addrs)
	}
}

func WithRetry(c int) commoninterface.ParamOption {
	return func(p commoninterface.ParamOptions) {
		e := p.(*EtcdOptions)
		e.retry = c
	}
}

func WithHeartbeatInterval(s uint64) commoninterface.ParamOption {
	return func(p commoninterface.ParamOptions) {
		e := p.(*EtcdOptions)
		e.heartInterval = s
	}
}
