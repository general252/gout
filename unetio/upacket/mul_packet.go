package upacket

import (
	"bytes"
	"github.com/general252/gout/uerror"
	"github.com/general252/gout/ulog"
	"time"
)

const (
	MTU = 1500 // 576
	// Internet上的标准MTU值为576字节
	// unix网络编程第一卷里说：ipv4协议规定ip层的最小重组缓冲区大小为576
	UdpPacketMaxSize        = MTU - 8 - 20
	UdpPacketPayloadMaxSize = UdpPacketMaxSize - UdpPacketHeadSize

	MaxUdpSize  = 20 * 102 * 1024 // 20MB
	MaxUdpCount = MaxUdpSize / MTU
)

// 分包的udp
type MulUdpPacket struct {
	firstRecvTime time.Time
	packetArray   []*UdpPacket
}

func NewMulUdpPacket(count int) *MulUdpPacket {
	return &MulUdpPacket{
		firstRecvTime: time.Now(),
		packetArray:   make([]*UdpPacket, count),
	}
}

// 分包接收完整
func (c *MulUdpPacket) IsRecvAllPacket() bool {
	if c.packetArray == nil {
		return false
	}

	for i := 0; i < len(c.packetArray); i++ {
		if c.packetArray[i] == nil {
			return false
		}
	}

	return true
}

// seq分包接收超时
func (c *MulUdpPacket) IsTimeout(now time.Time, duration time.Duration) bool {
	if now.Sub(c.firstRecvTime) > duration {
		return true
	}
	return false
}

// 接收到一个包
func (c *MulUdpPacket) AddPacket(pkt *UdpPacket) error {
	var count = len(c.packetArray)
	if pkt.Index() >= count || pkt.Index() < 0 {
		return uerror.WithMessageF("packet index error, index: %v count: %v", pkt.Index(), count)
	}

	c.packetArray[pkt.Index()] = pkt
	return nil
}

func (c *MulUdpPacket) GetSeq() (uint16, error) {
	var count = len(c.packetArray)
	for i := 0; i < count; i++ {
		if c.packetArray[i] != nil {
			return c.packetArray[i].PktSeq, nil
		}
	}

	return 0, uerror.WithMessage("mul packet no packet")
}

func (c *MulUdpPacket) CheckPacket(ptk *UdpPacket) bool {
	var count = len(c.packetArray)
	for i := 0; i < count; i++ {
		if c.packetArray[i] != nil {
			if c.packetArray[i].PktSeq == ptk.PktSeq && c.packetArray[i].Magic == ptk.Magic {
				return true
			} else {
				return false
			}
		}
	}

	ulog.WarnF("inner error")
	return false
}

// 拼接分包
func (c *MulUdpPacket) GetPacketData() ([]byte, error) {
	var buf bytes.Buffer

	var count = len(c.packetArray)
	for i := 0; i < count; i++ {
		if c.packetArray[i] == nil {
			return nil, uerror.WithMessage("count: %v, index %v is nil", count, i)
		}
		buf.Write(c.packetArray[i].Payload.Bytes())
	}

	return buf.Bytes(), nil
}
