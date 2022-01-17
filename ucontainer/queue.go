package ucontainer

import (
	"container/list"
	"sync"
)

type Queue struct {
	mux sync.Mutex

	l *list.List
}

func NewQueue() *Queue {
	return &Queue{
		l: list.New(),
	}
}

func (c *Queue) Push(v interface{}) {
	c.mux.Lock()
	defer c.mux.Unlock()

	c.l.PushBack(v)
}

func (c *Queue) Pop() (interface{}, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.l.Len() == 0 {
		return nil, false
	}

	f := c.l.Front()
	if f == nil {
		return nil, false
	}

	c.l.Remove(f)

	return f.Value, true
}

func (c *Queue) Len() int {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.l.Len()
}

func (c *Queue) Peek() (interface{}, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if c.l.Len() == 0 {
		return nil, false
	}

	f := c.l.Front()
	if f == nil {
		return nil, false
	}

	return f.Value, true
}
