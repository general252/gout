package udp_packet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/general252/gout/uerror"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/unetio"
	"hash/crc32"
	"math"
	"net"
)

// UdpPacket UDP数据包
type UdpPacket struct {
	PktSeq     uint32         // 4
	PktCount   uint16         // 2 分包数
	PktIndex   uint16         // 2 分包索引 0, 1, 2, ..., (PktCount-1)
	PktReserve [2]uint8       // 2 预留
	PktType    unetio.PktType // 1 包类型
	PktCRC     uint8          // 1 包校验

	Payload    bytes.Buffer // 包(分包)的负载数据
	RemoteAddr net.UDPAddr
}

// 打包
func (c *UdpPacket) Packet() []byte {
	var buf = &bytes.Buffer{}

	_ = binary.Write(buf, binary.BigEndian, c.PktSeq)
	_ = binary.Write(buf, binary.BigEndian, c.PktCount)
	_ = binary.Write(buf, binary.BigEndian, c.PktIndex)
	_ = binary.Write(buf, binary.BigEndian, c.PktReserve)
	_ = binary.Write(buf, binary.BigEndian, c.PktType)
	c.PktCRC = uint8(crc32.ChecksumIEEE(buf.Bytes()) % math.MaxUint8)
	_ = binary.Write(buf, binary.BigEndian, c.PktCRC)

	buf.Write(c.Payload.Bytes())

	return buf.Bytes()
}

// 解包
func (c *UdpPacket) UnPacket(pktData []byte) error {
	if len(pktData) < unetio.UdpPacketHeadSize {
		return uerror.WithMessageF("head >= %v", unetio.UdpPacketHeadSize)
	}

	var buf = bytes.NewReader(pktData)

	_ = binary.Read(buf, binary.BigEndian, &c.PktSeq)
	_ = binary.Read(buf, binary.BigEndian, &c.PktCount)
	_ = binary.Read(buf, binary.BigEndian, &c.PktIndex)
	_ = binary.Read(buf, binary.BigEndian, &c.PktReserve)
	_ = binary.Read(buf, binary.BigEndian, &c.PktType)
	_ = binary.Read(buf, binary.BigEndian, &c.PktCRC)

	calcCRC := uint8(crc32.ChecksumIEEE(pktData[:unetio.UdpPacketHeadSize-1]) % math.MaxUint8)
	if calcCRC != c.PktCRC {
		return uerror.WithMessageF("crc check fail. %v != %v head len: %v\n%v",
			calcCRC, c.PktCRC, len(pktData), c)
	}

	if c.PktIndex >= c.PktCount {
		return uerror.WithMessageF("PktIndex error. PktIndex: %v PktCount: %v", c.PktIndex, c.PktCount)
	}
	if c.PktCount > 1024 {
		// 540*1024 = 54KB
		ulog.WarnF("maybe wrong packet. PktSeq: %v PktIndex: %v PktCount: %v", c.PktSeq, c.PktIndex, c.PktCount)
	}

	c.Payload.Write(pktData[unetio.UdpPacketHeadSize:])

	return nil
}

// 包的负载数据
func (c *UdpPacket) PayloadData() []byte {
	return c.Payload.Bytes()
}

// 包的序列号
func (c *UdpPacket) Seq() uint32 {
	return c.PktSeq
}

// 数据包来源
func (c *UdpPacket) SetRemoteAddr(addr net.UDPAddr) {
	c.RemoteAddr = addr
}

func (c *UdpPacket) RemoteUdpAddr() net.UDPAddr {
	return c.RemoteAddr
}

func (c *UdpPacket) ConnSeqId() string {
	return fmt.Sprintf("%v_%v", c.RemoteAddr, c.PktSeq)
}

// 是否分包
func (c *UdpPacket) IsMulPacket() bool {
	if c.PktCount > 1 {
		return true
	}

	return false
}

// 分包个数
func (c *UdpPacket) Count() int {
	return int(c.PktCount)
}

// 此包索引0,1,2...
func (c *UdpPacket) Index() int {
	return int(c.PktIndex)
}

func (c *UdpPacket) String() string {
	return fmt.Sprintf("PktSeq: %v, PktCount: %v, PktIndex: %v", c.PktSeq, c.PktCount, c.PktIndex)
}
