package upacket

import (
	"github.com/general252/gout/uerror"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/usafe"
	"log"
	"net"
)

type HandRecvPacket func(pkt *MulUdpPacket, payload []byte)
type HandTimeoutPacket func(pkt *MulUdpPacket)
type HandLossPacket func(pktIndex int)

type UdpServer struct {
	port int
	conn *net.UDPConn

	handRecvPacket    HandRecvPacket
	handTimeoutPacket HandTimeoutPacket
	handLossPacket    HandLossPacket
}

func NewUdpServer(port int, handRecvPacket HandRecvPacket, handTimeoutPacket HandTimeoutPacket, handLossPacket HandLossPacket) *UdpServer {
	return &UdpServer{port: port,
		handRecvPacket:    handRecvPacket,
		handTimeoutPacket: handTimeoutPacket,
		handLossPacket:    handLossPacket,
	}
}

func (c *UdpServer) Open() error {
	var err error
	c.conn, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: c.port})
	if err != nil {
		return uerror.WithErrorAndMessage(err, "listen udp fail")
	}

	// 设置缓冲区大小
	if err = c.conn.SetReadBuffer(16 * 1024 * 1024); err != nil {
		log.Printf("SetReadBuffer fail %v", err)
	}
	if err = c.conn.SetWriteBuffer(16 * 1024 * 1024); err != nil {
		log.Printf("SetWriteBuffer fail %v", err)
	}

	usafe.Go(func() {
		// session 长时间没有接收到数据, 删除
		var sessionList = make(map[string]*UdpSession)
		var buffer = make([]byte, 65536)

		for {
			n, addr, err := c.conn.ReadFromUDP(buffer)
			if err != nil {
				ulog.WarnF("==udp conn read udp end==")
				break
			}

			clientAddress := addr.String()

			session, ok := sessionList[clientAddress]
			if !ok {
				session = NewUdpSession(c, clientAddress)
				sessionList[clientAddress] = session
			}

			session.OnRecvData(buffer[:n])
		}
	})

	return nil
}

func (c *UdpServer) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
		return
	}
}
