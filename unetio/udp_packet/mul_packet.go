package udp_packet

import (
	"fmt"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/unetio"
	"sort"
	"time"
)

// 分包的udp
type mulUdpPacket struct {
	addTime     time.Time
	packetArray []*UdpPacket
}

// 分包接收完整
func (c *mulUdpPacket) IsRecvAllPacket() bool {
	if c.packetArray == nil {
		return false
	}
	if len(c.packetArray) == 0 {
		return false
	}
	if c.packetArray[0].Count() != len(c.packetArray) {
		return false
	}
	return true
}

// seq分包接收超时
func (c *mulUdpPacket) IsTimeout(now time.Time) bool {
	if now.Sub(c.addTime) > unetio.MulUdpPacketTimeout {
		return true
	}
	return false
}

// 接收到一个包
func (c *mulUdpPacket) RecvPacket(pkt *UdpPacket) {
	for _, packet := range c.packetArray {
		if pkt.Index() == packet.Index() {
			ulog.WarnF("packet PktIndex repeat. addPacket: %v. %v", pkt, packet)
			return
		}
	}

	c.packetArray = append(c.packetArray, pkt)
}

// 拼接分包
func (c *mulUdpPacket) GetPacketData() []byte {
	sort.Slice(c.packetArray, func(i, j int) bool {
		return c.packetArray[i].Index() < c.packetArray[j].Index()
	})

	var payload []byte
	for i := 0; i < len(c.packetArray); i++ {
		payload = append(payload, c.packetArray[i].PayloadData()...)
	}
	return payload
}

// 分包的seq
func (c *mulUdpPacket) GetSeq() uint16 {
	if c.packetArray == nil || len(c.packetArray) <= 0 {
		return 0
	}

	return c.packetArray[0].Seq()
}

// 序列化
func (c *mulUdpPacket) String() string {
	if c.packetArray == nil || len(c.packetArray) <= 0 {
		return "no packet"
	}

	var str string
	for i, packet := range c.packetArray {
		str += fmt.Sprintf(" %v: %v\n", i, packet)
	}
	return str
}
