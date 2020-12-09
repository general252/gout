package udp_packet

import (
	"bytes"
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
	pktSeq   uint16        // 65536
	pktCount uint16        // 分包数
	pktIndex uint16        // 分包索引 0, 1, 2, ..., (pktCount-1)
	pktType  unetio.PktType // 包类型
	pktCRC   uint8         // 包校验

	// UdpPacketHeadSize = 8

	payload    bytes.Buffer // 包(分包)的负载数据
	remoteAddr net.UDPAddr
}

// 序列化头
func (c *UdpPacket) serHead() []byte {
	var result []byte

	result = append(result, unetio.Uint16ToBytes(c.pktSeq)...)
	result = append(result, unetio.Uint16ToBytes(c.pktCount)...)
	result = append(result, unetio.Uint16ToBytes(c.pktIndex)...)
	result = append(result, uint8(c.pktType))
	result = append(result, uint8(crc32.ChecksumIEEE(result)%math.MaxUint8))

	return result
}

// 反序列化头
func (c *UdpPacket) unSerHead(head []byte) error {
	if len(head) < unetio.UdpPacketHeadSize {
		return uerror.WithMessageF("head >= %v", unetio.UdpPacketHeadSize)
	}

	c.pktSeq = unetio.BytesToUint16(head[0:2])
	c.pktCount = unetio.BytesToUint16(head[2:4])
	c.pktIndex = unetio.BytesToUint16(head[4:6])
	c.pktType = unetio.PktType(head[6])
	c.pktCRC = head[7]

	calcCRC := uint8(crc32.ChecksumIEEE(head[:7]) % math.MaxUint8)
	if calcCRC != c.pktCRC {
		return uerror.WithMessageF("crc check fail. %v != %v head len: %v\n%v",
			calcCRC, c.pktCRC, len(head), c)
	}

	if c.pktIndex >= c.pktCount {
		return uerror.WithMessageF("pktIndex error. pktIndex: %v pktCount: %v", c.pktIndex, c.pktCount)
	}
	if c.pktCount > 1024 {
		// 540*1024 = 54KB
		ulog.WarnF("maybe wrong packet. pktSeq: %v pktIndex: %v pktCount: %v", c.pktSeq, c.pktIndex, c.pktCount)
	}
	return nil
}

// 打包
func (c *UdpPacket) Packet() []byte {
	return append(c.serHead(), c.payload.Bytes()...)
}

// 解包
func (c *UdpPacket) unPacket(data []byte) error {
	if err := c.unSerHead(data); err != nil {
		return err
	} else {
		_, err := c.payload.Write(data[unetio.UdpPacketHeadSize:])
		if err != nil {
			return uerror.WithErrorAndMessage(err, "c.payload write fail.")
		}
		return nil
	}
}

// 包的负载数据
func (c *UdpPacket) Payload() []byte {
	return c.payload.Bytes()
}

// 包的序列号
func (c *UdpPacket) Seq() uint16 {
	return c.pktSeq
}

// 数据包来源
func (c *UdpPacket) SetRemoteAddr(addr net.UDPAddr) {
	c.remoteAddr = addr
}

func (c *UdpPacket) RemoteAddr() net.UDPAddr {
	return c.remoteAddr
}

func (c *UdpPacket) ConnSeqId() string {
	return fmt.Sprintf("%v_%v", c.remoteAddr, c.pktSeq)
}

// 是否分包
func (c *UdpPacket) IsMulPacket() bool {
	if c.pktCount > 1 {
		return true
	}

	return false
}

// 分包个数
func (c *UdpPacket) Count() int {
	return int(c.pktCount)
}

// 此包索引0,1,2...
func (c *UdpPacket) Index() int {
	return int(c.pktIndex)
}

func (c *UdpPacket) String() string {
	return fmt.Sprintf("pktSeq: %v, pktCount: %v, pktIndex: %v", c.pktSeq, c.pktCount, c.pktIndex)
}
