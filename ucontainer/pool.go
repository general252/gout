package ucontainer

import (
	"errors"
	"sync"
)

const DefaultMaxPoolSize = 512 * 1024 // 512KB

var SizeTooLarge = errors.New("size large")

type Pool struct {
	mux sync.Mutex

	maxPoolSize int
	pos         int
	buf         []byte
}

func NewPool() *Pool {
	return &Pool{
		maxPoolSize: DefaultMaxPoolSize,
		buf:         make([]byte, DefaultMaxPoolSize),
	}
}

func NewPoolWithSize(poolSize int) *Pool {
	if poolSize <= 0 {
		panic("error pool size")
	}

	return &Pool{
		maxPoolSize: poolSize,
		buf:         make([]byte, poolSize),
	}
}

func (c *Pool) Get(size int) ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if size > c.maxPoolSize {
		return nil, SizeTooLarge
	}

	if c.maxPoolSize-c.pos < size {
		c.pos = 0
		c.buf = make([]byte, c.maxPoolSize)
	}

	b := c.buf[c.pos : c.pos+size]
	c.pos += size

	return b, nil
}
