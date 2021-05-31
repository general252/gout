package upool

import "time"

const (
	defaultPoolSize      = 1000
	defaultSyncWriteIdle = time.Second * 2
)

type poolOptions struct {
	poolSize      int64
	syncWriteIdle time.Duration
}

type PoolOption interface {
	apply(*poolOptions)
}

// funcOption 返回PoolOption对象较繁琐, 使用函数封装一层
type funcOption struct {
	f func(*poolOptions)
}

func (fdo *funcOption) apply(do *poolOptions) {
	fdo.f(do)
}

func newFuncOption(f func(*poolOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithPoolSize(s int64) PoolOption {
	return newFuncOption(func(o *poolOptions) {
		o.poolSize = s
	})
}

func WithSyncWaitTime(duration time.Duration) PoolOption {
	return newFuncOption(func(o *poolOptions) {
		o.syncWriteIdle = duration
	})
}
