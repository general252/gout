package example

import (
	"context"
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/unetio/udp_packet"
	"net"
)

func Server(port int) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		ulog.Error(err)
		return
	}
	defer udpConn.Close()

	// 设置写缓冲区大小
	_ = udpConn.SetWriteBuffer(512 * 1024 * 1024)

	var connMap = make(map[string]*udp_packet.UdpPacketServer)
	var buffer = make([]byte, 65536)

	for {
		n, addr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			ulog.Error(err)
			continue
		}

		var pktFactory = connMap[addr.String()]
		if pktFactory == nil {
			pktFactory = udp_packet.NewUdpPacketServer(context.TODO(), func(payload []byte) {
				ulog.Info("收到数据: ", len(payload))
			})
			connMap[addr.String()] = pktFactory
		}

		pktFactory.PushPacketData(*addr, buffer[:n])
	}
}
