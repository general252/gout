package udp_packet

import (
	"github.com/general252/gout/unetio"
	"math/rand"
	"sync"
	"time"
)

type UdpPacketClient struct {
	seq  uint32
	lock sync.Mutex
}

func NewUdpPacketClient() *UdpPacketClient {
	rand.Seed(time.Now().UnixNano())

	return &UdpPacketClient{
		seq: rand.Uint32(),
	}
}

func (c *UdpPacketClient) Serialization(data []byte) []UdpPacket {
	var seq = c.GetSeq()

	if len(data) <= unetio.UdpPacketPayloadMaxSize {
		var pkt = UdpPacket{
			PktSeq:   seq,
			PktCount: 1,
			PktIndex: 0,
			PktType:  unetio.PktTypeData,
		}
		_, _ = pkt.Payload.Write(data)

		return []UdpPacket{pkt}
	} else {
		var dataLength = len(data)

		var n = dataLength / unetio.UdpPacketPayloadMaxSize
		if dataLength%unetio.UdpPacketPayloadMaxSize > 0 {
			n += 1
		}

		var result = make([]UdpPacket, 0, n)
		for i := 0; i < n; i++ {
			var offset = i * unetio.UdpPacketPayloadMaxSize
			var end = offset + unetio.UdpPacketPayloadMaxSize
			if end > dataLength {
				end = dataLength
			}

			var pkt = UdpPacket{
				PktSeq:   seq,
				PktCount: uint16(n),
				PktIndex: uint16(i),
			}
			_, _ = pkt.Payload.Write(data[offset:end])

			result = append(result, pkt)
		}

		return result
	}
}

func (c *UdpPacketClient) GetSeq() uint32 {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.seq++
	return c.seq
}
