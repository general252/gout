// ref https://github.com/siddontang/go/tree/master/sync2
package usync

import (
	"sync"
	"sync/atomic"
	"time"
)

type AtomicInt32 int32

func (i *AtomicInt32) Add(n int32) int32 {
	return atomic.AddInt32((*int32)(i), n)
}

func (i *AtomicInt32) Set(n int32) {
	atomic.StoreInt32((*int32)(i), n)
}

func (i *AtomicInt32) Get() int32 {
	return atomic.LoadInt32((*int32)(i))
}

func (i *AtomicInt32) CompareAndSwap(oldVal, newVal int32) (swapped bool) {
	return atomic.CompareAndSwapInt32((*int32)(i), oldVal, newVal)
}

type AtomicUint32 uint32

func (i *AtomicUint32) Add(n uint32) uint32 {
	return atomic.AddUint32((*uint32)(i), n)
}

func (i *AtomicUint32) Set(n uint32) {
	atomic.StoreUint32((*uint32)(i), n)
}

func (i *AtomicUint32) Get() uint32 {
	return atomic.LoadUint32((*uint32)(i))
}

func (i *AtomicUint32) CompareAndSwap(oldVal, newVal uint32) (swapped bool) {
	return atomic.CompareAndSwapUint32((*uint32)(i), oldVal, newVal)
}

type AtomicInt64 int64

func (i *AtomicInt64) Add(n int64) int64 {
	return atomic.AddInt64((*int64)(i), n)
}

func (i *AtomicInt64) Set(n int64) {
	atomic.StoreInt64((*int64)(i), n)
}

func (i *AtomicInt64) Get() int64 {
	return atomic.LoadInt64((*int64)(i))
}

func (i *AtomicInt64) CompareAndSwap(oldVal, newVal int64) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(i), oldVal, newVal)
}

type AtomicUint64 uint64

func (i *AtomicUint64) Add(n uint64) uint64 {
	return atomic.AddUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Set(n uint64) {
	atomic.StoreUint64((*uint64)(i), n)
}

func (i *AtomicUint64) Get() uint64 {
	return atomic.LoadUint64((*uint64)(i))
}

func (i *AtomicUint64) CompareAndSwap(oldVal, newVal uint64) (swapped bool) {
	return atomic.CompareAndSwapUint64((*uint64)(i), oldVal, newVal)
}

type AtomicDuration int64

func (d *AtomicDuration) Add(duration time.Duration) time.Duration {
	return time.Duration(atomic.AddInt64((*int64)(d), int64(duration)))
}

func (d *AtomicDuration) Set(duration time.Duration) {
	atomic.StoreInt64((*int64)(d), int64(duration))
}

func (d *AtomicDuration) Get() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(d)))
}

func (d *AtomicDuration) CompareAndSwap(oldVal, newVal time.Duration) (swapped bool) {
	return atomic.CompareAndSwapInt64((*int64)(d), int64(oldVal), int64(newVal))
}

// AtomicString gives you atomic-style APIs for string, but
// it's only a convenience wrapper that uses a mutex. So, it's
// not as efficient as the rest of the atomic types.
type AtomicString struct {
	mu  sync.Mutex
	str string
}

func (s *AtomicString) Set(str string) {
	s.mu.Lock()
	s.str = str
	s.mu.Unlock()
}

func (s *AtomicString) Get() string {
	s.mu.Lock()
	str := s.str
	s.mu.Unlock()
	return str
}

func (s *AtomicString) CompareAndSwap(oldVal, newVal string) (swapped bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.str == oldVal {
		s.str = newVal
		return true
	}
	return false
}

type AtomicBool int32

func (b *AtomicBool) Set(v bool) {
	if v {
		atomic.StoreInt32((*int32)(b), 1)
	} else {
		atomic.StoreInt32((*int32)(b), 0)
	}
}

func (b *AtomicBool) Get() bool {
	return atomic.LoadInt32((*int32)(b)) == 1
}
