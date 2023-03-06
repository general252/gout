package ustring

import (
	"testing"
)

func BenchmarkStr2bytes(b *testing.B) {
	s := "testString"
	var bs []byte
	for n := 0; n < b.N; n++ {
		bs = String2Bytes(s)
	}
	_ = bs
}

// 不推荐：用类型转换实现string到bytes
func BenchmarkStr2bytes2(b *testing.B) {
	s := "testString"
	var bs []byte
	for n := 0; n < b.N; n++ {
		bs = []byte(s)
	}
	_ = bs
}

// 推荐：用unsafe.Pointer实现bytes到string
func BenchmarkBytes2str(b *testing.B) {
	bs := String2Bytes("testString")
	var s string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s = Bytes2String(bs)
	}
	_ = s
}

// 不推荐：用类型转换实现bytes到string
func BenchmarkBytes2str2(b *testing.B) {
	bs := String2Bytes("testString")
	var s string
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s = string(bs)
	}
	_ = s
}
