package upacket

import (
	"github.com/general252/gout/uerror"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type UdpClient struct {
	host string
	port int
	conn *net.UDPConn

	seq  uint16
	lock sync.Mutex
}

func NewUdpClient(host string, port int) *UdpClient {
	return &UdpClient{
		host: host,
		port: port,
		seq:  0,
	}
}

func (c *UdpClient) Open() error {
	var err error
	c.conn, err = net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(c.host), Port: c.port})
	if err != nil {
		return uerror.WithErrorAndMessage(err, "net.DialUDP fail")
	}

	if err = c.conn.SetReadBuffer(16 * 1024 * 1024); err != nil {
		log.Printf("SetReadBuffer fail %v", err)
	}
	if err = c.conn.SetWriteBuffer(16 * 1024 * 1024); err != nil {
		log.Printf("SetWriteBuffer fail %v", err)
	}

	return nil
}

func (c *UdpClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
}

// 发送数据
func (c *UdpClient) Write(data []byte) error {
	return c.serialization(data)
}

// 拆包并发送 -> packet() -> udpWrite()
func (c *UdpClient) serialization(data []byte) error {
	var seq = c.getSeq()
	var magic = uint8(rand.Uint32())

	var dataLength = len(data)
	if dataLength <= UdpPacketPayloadMaxSize {

		return c.packet(data, 1, 0, seq, magic)
	} else {

		var packetCount = dataLength / UdpPacketPayloadMaxSize
		if dataLength%UdpPacketPayloadMaxSize > 0 {
			packetCount += 1
		}

		for i := 0; i < packetCount; i++ {
			var offset = i * UdpPacketPayloadMaxSize
			var end = offset + UdpPacketPayloadMaxSize
			if end > dataLength {
				end = dataLength
			}

			// 分包发送
			if err := c.packet(data[offset:end], uint16(packetCount), uint16(i), seq, magic); err != nil {
				return err
			}
		}

		return nil
	}
}

// 打包发送
func (c *UdpClient) packet(data []byte, count uint16, index uint16, seq uint16, magic uint8) error {
	var pkt = UdpPacket{
		PktSeq:   seq,
		PktCount: count,
		PktIndex: index,
		Magic:    magic,
	}

	if _, err := pkt.Payload.Write(data); err != nil {
		return uerror.WithErrorAndMessage(err, "buf write fail")
	}

	// 打包
	pktData, err := pkt.Packet()
	if err != nil {
		return err
	}

	// 发送
	if err := c.udpWrite(pktData); err != nil {
		return err
	}

	return nil
}

// udp 发送
func (c *UdpClient) udpWrite(data []byte) error {
	if c.conn == nil {
		return uerror.WithMessage("udp conn is nil")
	}

	if _, err := c.conn.Write(data); err != nil {
		return uerror.WithErrorAndMessage(err, "udp conn write fail")
	}

	return nil
}

func (c *UdpClient) getSeq() uint16 {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.seq++
	return c.seq
}
