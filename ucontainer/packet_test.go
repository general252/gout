package ucontainer

import (
	"fmt"
	"log"
	"time"
)

func ExampleNewPacketBuffer() {
	pktBuffer := NewPacketBuffer()

	go func() {
		for i := 0; i < 10; i++ {
			_, _ = pktBuffer.Write([]byte(fmt.Sprintf("hello %v", i)))
		}

		_ = pktBuffer.Close()
	}()

	time.Sleep(time.Second * 3)
	for {
		pkt, err := pktBuffer.Read()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(string(pkt))
	}

	// output:
}
