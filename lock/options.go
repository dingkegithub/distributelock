package lock

import "time"

type Option func(*LockOptions)

//
// Lock option
// Block: true block until get lock, false return immediately wheather get lock
// Timeout: wait lock time
//
type LockOptions struct {
	Block     bool
	Tmoutable bool
	Timeout   time.Duration
}

func (opt *LockOptions) Apply(opts ...Option) {
	for _, o := range opts {
		o(opt)
	}
}

func WithBlock() Option {
	return func(lo *LockOptions) {
		lo.Block = true
		lo.Tmoutable = false
		lo.Timeout = 0
	}
}

func WithTimeout(timeout uint64) Option {
	return func(lo *LockOptions) {
		lo.Tmoutable = true
		lo.Timeout = time.Duration(timeout) * time.Millisecond
		lo.Block = false
	}
}
