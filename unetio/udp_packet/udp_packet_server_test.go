package udp_packet

import (
	"bytes"
	"context"
	"github.com/general252/gout/unetio"
	"log"
	"math"
	"net"
	"testing"
	"time"
)

// 测试序列化反序列化
func TestNewUdpPacketServer(t *testing.T) {
	var srcData = make([]byte, 2020)
	for i := 0; i < len(srcData); i++ {
		srcData[i] = uint8(i % math.MaxUint8)
	}

	var cli = NewUdpPacketClient()
	var ser = NewUdpPacketServer(context.TODO(), func(pktHeadInfo *MulUdpPacket, payload []byte) {
		log.Println(pktHeadInfo)
	}, func(pktHeadInfo *MulUdpPacket) {
		// 超时 10 sec
	}, func(lossSeq uint32) {
		// 丢包 5 sec
	})

	srcPktList := cli.Serialization(srcData)
	var dstPktList []*UdpPacket
	for i := 0; i < len(srcPktList); i++ {
		tmpPkt, err := ser.Deserialization(srcPktList[i].Packet())
		if err != nil {
			t.Error(err)
		}

		dstPktList = append(dstPktList, tmpPkt)
	}

	var dstData bytes.Buffer
	for i := 0; i < len(dstPktList); i++ {
		dstData.Write(dstPktList[i].PayloadData())
	}

	if bytes.Compare(srcData, dstData.Bytes()) != 0 {
		t.Errorf("fail")
	}

	t.Logf("success")
}

var (
	port = 8800
)

func ExampleNewUdpPacketServer_Server() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	udpConnServer, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		log.Printf("ListenUDP fail %v", err)
		return
	}
	defer udpConnServer.Close()
	_ = udpConnServer.SetReadBuffer(512 * 1024 * 1024)
	_ = udpConnServer.SetWriteBuffer(512 * 1024 * 1024)

	var server = func() {
		var connMap = make(map[string]*UdpPacketServer)
		buffer := make([]byte, unetio.UdpPacketMaxSize*2)
		for {
			n, addr, err := udpConnServer.ReadFromUDP(buffer)
			if err != nil {
				log.Printf("ReadFromUDP fail %v", err)
				return
			}
			if n <= 0 {
				continue
			}

			var pktConnServer = connMap[addr.String()]
			if pktConnServer == nil {
				log.Printf("new connection %v", addr)

				pktConnServer = NewUdpPacketServer(context.TODO(), func(pktHeadInfo *MulUdpPacket, payload []byte) {
					//log.Printf("recv msg len: %v", len(payload))
				}, func(pktHeadInfo *MulUdpPacket) {
					// 超时 10 sec
					log.Printf("超时 %v", pktHeadInfo.GetSeq())
				}, func(lossSeq uint32) {
					// 丢包 5 sec
					log.Printf("丢包 %v", lossSeq)
				})
				connMap[addr.String()] = pktConnServer
			}

			pktConnServer.PushPacketData(*addr, buffer[:n])
		}
	}
	server()

	// output:
	//
}

func ExampleNewUdpPacketServer_Client() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var client = func(port int) {
		udpConnCli, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: port})
		if err != nil {
			log.Printf("DialUDP fail %v", err)
			return
		}
		defer udpConnCli.Close()
		_ = udpConnCli.SetReadBuffer(512 * 1024 * 1024)
		_ = udpConnCli.SetWriteBuffer(512 * 1024 * 1024)

		cli := NewUdpPacketClient()

		var sendData bytes.Buffer
		sendData.WriteString("hello world")
		sendData.Write(make([]byte, 20000))
		for {
			for i := 0; i < 10; i++ {
				var pktArray = cli.Serialization(sendData.Bytes())
				for _, packet := range pktArray {
					if _, err := udpConnCli.Write(packet.Packet()); err != nil {
						log.Printf("udp write fail %v", err)
					}
				}
			}

			time.Sleep(time.Millisecond * 50)
		}
	}

	go client(port)
	go client(port)
	go client(port)
	go client(port)
	go client(port)
	go client(port)
	go client(port)
	client(port)

	// output:
	//
}
