package unetio

import (
	"bytes"
	"encoding/binary"
	"github.com/general252/gout/ulog"
	"time"
)

const (
	MTU               = 1500 // 576
	UdpPacketHeadSize = 12
	// Internet上的标准MTU值为576字节
	// unix网络编程第一卷里说：ipv4协议规定ip层的最小重组缓冲区大小为576
	UdpPacketMaxSize        = MTU - 8 - 20
	UdpPacketPayloadMaxSize = UdpPacketMaxSize - UdpPacketHeadSize

	/**
	uint32:uint16:uint16
	包的序号:分包数:分包序号  // (512-4-2-2)*65536=33,030,144

	uint32:uint8:uint8
	包的序号:分包数:分包序号  // (512-4-1-1)*256=129,536
	*/

	MulUdpPacketTimeout = time.Second * time.Duration(10)
)

type PktType uint8

const (
	PktTypeData PktType = 1 // 数据包
	PktTypeReq  PktType = 2 // 请求包, 请求重发丢失的包
)

//整形转换成字节
func Uint16ToBytes(x uint16) []byte {
	//var out = make([]byte, 2)
	//binary.BigEndian.PutUint16(out, x)
	//return out

	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, x); err != nil {
		ulog.ErrorF("binary write fail. %v", err)
		return nil
	}

	return bytesBuffer.Bytes()
}

// 字节转换成整形
func BytesToUint16(b []byte) uint16 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint16
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		ulog.ErrorF("binary read fail. %v", err)
		return 0
	}

	return x
}

func BytesToUint64(b []byte) uint64 {
	bytesBuffer := bytes.NewBuffer(b)

	var x uint64
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); err != nil {
		ulog.ErrorF("binary read fail. %v", err)
		return 0
	}

	return x
}

func Uint64ToBytes(x uint64) []byte {
	//var out = make([]byte, 2)
	//binary.BigEndian.PutUint16(out, x)
	//return out

	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, x); err != nil {
		ulog.ErrorF("binary write fail. %v", err)
		return nil
	}

	return bytesBuffer.Bytes()
}
