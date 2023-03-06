package ustring

import (
	"reflect"
	"unsafe"
)

func String2Bytes(s string) []byte {
	if true {
		x := (*reflect.StringHeader)(unsafe.Pointer(&s))
		h := reflect.SliceHeader{
			Data: x.Data,
			Len:  x.Len,
			Cap:  x.Len,
		}
		return *(*[]byte)(unsafe.Pointer(&h))
	} else {
		x := (*[2]uintptr)(unsafe.Pointer(&s))
		h := [3]uintptr{x[0], x[1], x[1]}
		return *(*[]byte)(unsafe.Pointer(&h))
	}
}

func Bytes2String(b []byte) string {
	if true {
		x := (*reflect.SliceHeader)(unsafe.Pointer(&b))
		h := reflect.StringHeader{
			Data: x.Data,
			Len:  x.Len,
		}
		return *(*string)(unsafe.Pointer(&h))
	} else {
		return *(*string)(unsafe.Pointer(&b))
	}
}
