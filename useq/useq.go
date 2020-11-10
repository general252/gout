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

// GetSeq64 获取唯一int64值
func GetSeq64() int64 {
	mux.Lock()
	defer mux.Unlock()

	seq64++
	return seq64
}

func GetSeqU64() uint64 {
	mux.Lock()
	defer mux.Unlock()

	seqU64++
	return seqU64
}

func GetSeq32() int32 {
	mux.Lock()
	defer mux.Unlock()

	seq32++
	return seq32
}

func GetSeqU32() uint32 {
	mux.Lock()
	defer mux.Unlock()

	seqU32++
	return seqU32
}

func GetSeq16() int16 {
	mux.Lock()
	defer mux.Unlock()

	seq16++
	return seq16
}

func GetSeqU16() uint16 {
	mux.Lock()
	defer mux.Unlock()

	seqU16++
	return seqU16
}
