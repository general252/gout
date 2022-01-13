package ushutdown

import (
	"log"
	"testing"
)

func TestNewHook(t *testing.T) {
	NewHook().WithStopHandler(
		func() {
			log.Println("http server stop")
		},
		func() {
			log.Println("sip server stop")
		},
		func() {
			log.Println("rtsp server stop")
		},
	).GracefulStop()
}
