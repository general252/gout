package upacket

import (
	"bytes"
	"encoding/binary"
	"github.com/general252/gout/uerror"
	"hash/crc32"
	"math"
)

const (
	UdpPacketHeadSize = 8
)

// UdpPacket UDP数据包
type UdpPacket struct {
	PktSeq   uint16 // 2
	PktCount uint16 // 2 分包数
	PktIndex uint16 // 2 分包索引 0, 1, 2, ..., (PktCount-1)
	Magic    uint8  // 1 魔法数
	PktCRC   uint8  // 1 包校验

	Payload bytes.Buffer // 包(分包)的负载数据
}

func (c *UdpPacket) Count() int {
	return int(c.PktCount)
}
func (c *UdpPacket) Index() int {
	return int(c.PktIndex)
}

// 打包
func (c *UdpPacket) Packet() ([]byte, error) {
	if c.PktCount > MaxUdpCount {
		return nil, uerror.WithMessageF("packet count error %v (MaxUdpCount: %v)", c.PktCount, MaxUdpCount)
	}

	var buf = &bytes.Buffer{}

	_ = binary.Write(buf, binary.BigEndian, c.PktSeq)
	_ = binary.Write(buf, binary.BigEndian, c.PktCount)
	_ = binary.Write(buf, binary.BigEndian, c.PktIndex)
	_ = binary.Write(buf, binary.BigEndian, c.Magic)
	c.PktCRC = uint8(crc32.ChecksumIEEE(buf.Bytes()) % math.MaxUint8)
	_ = binary.Write(buf, binary.BigEndian, c.PktCRC)

	buf.Write(c.Payload.Bytes())

	return buf.Bytes(), nil
}

// 解包
func (c *UdpPacket) UnPacket(pktData []byte) error {
	if len(pktData) < UdpPacketHeadSize {
		return uerror.WithMessageF("head >= %v", UdpPacketHeadSize)
	}

	var buf = bytes.NewReader(pktData)

	_ = binary.Read(buf, binary.BigEndian, &c.PktSeq)
	_ = binary.Read(buf, binary.BigEndian, &c.PktCount)
	_ = binary.Read(buf, binary.BigEndian, &c.PktIndex)
	_ = binary.Read(buf, binary.BigEndian, &c.Magic)
	_ = binary.Read(buf, binary.BigEndian, &c.PktCRC)

	// 检查参数
	if c.PktIndex >= c.PktCount || c.PktIndex < 0 {
		return uerror.WithMessageF("PktIndex error. PktIndex: %v PktCount: %v", c.PktIndex, c.PktCount)
	}

	if c.PktCount > MaxUdpCount {
		return uerror.WithMessageF("packet count error %v (MaxUdpCount: %v)", c.PktCount, MaxUdpCount)
	}

	// 计算校验和
	calcCRC := uint8(crc32.ChecksumIEEE(pktData[:UdpPacketHeadSize-1]) % math.MaxUint8)
	if calcCRC != c.PktCRC {
		return uerror.WithMessageF("crc check fail. %v != %v head len: %v\n%v",
			calcCRC, c.PktCRC, len(pktData), c)
	}

	c.Payload.Write(pktData[UdpPacketHeadSize:])

	return nil
}
