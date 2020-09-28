package zk

import (
	"github.com/dingkegithub/distrubutelock/pkg/commoninterface"
	"github.com/dingkegithub/distrubutelock/utils/log"
)

type zkOptions struct {
	logger log.Logging
	addrs  []string
}

func defaultzkOptions() *zkOptions {
	opts := &zkOptions{
		logger: &log.DefaultLogging{},
		addrs:  []string{"127.0.0.1:2181"},
	}

	return opts
}

func (zo *zkOptions) Apply(opts ...commoninterface.ParamOption) {
	for _, o := range opts {
		o(zo)
	}
}

func WithAddrs(addrs []string) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		zo := options.(*zkOptions)
		zo.addrs = make([]string, len(addrs))
		copy(zo.addrs, addrs)
	}
}

func WithLogger(logger log.Logging) commoninterface.ParamOption {
	return func(options commoninterface.ParamOptions) {
		zo := options.(*zkOptions)
		zo.logger = logger
	}
}
