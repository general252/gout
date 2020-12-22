package example

import (
	"github.com/general252/gout/ulog"
	"github.com/general252/gout/unetio/udp_packet"
	"math"
	"net"
	"time"
)

func Client(port int, ms int) {
	cliConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	})
	if err != nil {
		ulog.ErrorF(" net.DialUDP fail %v", err)
		return
	}
	_ = cliConn.SetWriteBuffer(512 * 1024 * 1024)

	data := make([]byte, 362*12)
	for i := 0; i < len(data); i++ {
		data[i] = uint8(i % math.MaxUint8)
	}

	cli := udp_packet.NewUdpPacketClient()
	for {
		for i := 0; i < 20; i++ {
			pktArray := cli.Serialization(data)

			for _, udpPacket := range pktArray {
				n, err := cliConn.Write(udpPacket.Packet())
				if err != nil {
					ulog.ErrorF("send fail %v", err)
				} else {
					_ = n
				}
			}
		}

		time.Sleep(time.Millisecond * time.Duration(ms))
	}
}
