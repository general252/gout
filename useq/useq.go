package useq

import (
	"sync"
)

var (
	mux    sync.Mutex
	seq64  int64  = 0
	seq32  int32  = 0
	seq16  int16  = 0
	seqU64 uint64 = 0
	seqU32 uint32 = 0
	seqU16 uint16 = 0
)

// Int64 获取唯一int64值
func Int64() int64 {
	mux.Lock()
	defer mux.Unlock()

	seq64++
	return seq64
}

func Uint64() uint64 {
	mux.Lock()
	defer mux.Unlock()

	seqU64++
	return seqU64
}

func Int32() int32 {
	mux.Lock()
	defer mux.Unlock()

	seq32++
	return seq32
}

func Uint32() uint32 {
	mux.Lock()
	defer mux.Unlock()

	seqU32++
	return seqU32
}

func Int16() int16 {
	mux.Lock()
	defer mux.Unlock()

	seq16++
	return seq16
}

func Uint16() uint16 {
	mux.Lock()
	defer mux.Unlock()

	seqU16++
	return seqU16
}
