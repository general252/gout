package ucontainer

import (
	"container/list"
	"sync"
)

type List[T any] struct {
	mu   sync.RWMutex
	vals list.List
}

func NewList[T any]() *List[T] {
	return &List[T]{}
}

func (tis *List[T]) PushBack(v T) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	tis.vals.PushBack(v)
}

func (tis *List[T]) PushFront(v T) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	tis.vals.PushFront(v)
}

func (tis *List[T]) Len() int {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	return tis.vals.Len()
}

func (tis *List[T]) Pop() {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	if e := tis.vals.Front(); e != nil {
		tis.vals.Remove(e)
	}
}

func (tis *List[T]) Range(fn func(v T) bool) {
	tis.mu.RLock()
	defer tis.mu.RUnlock()

	for e := tis.vals.Front(); e != nil; e = e.Next() {
		if !fn(e.Value.(T)) {
			return
		}
	}
}

func (tis *List[T]) Delete(value T, cmp func(a, b T) bool) {
	tis.mu.Lock()
	defer tis.mu.Unlock()

	for e := tis.vals.Front(); e != nil; e = e.Next() {
		if cmp(value, e.Value.(T)) {
			prev := e.Prev()
			if e == tis.vals.Front() {
				prev = e
			}
			tis.vals.Remove(e)
			e = prev
		}
	}
}
