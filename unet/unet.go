package unet

import (
	"net"
	"strings"
)

// get pc local host ip address
func GetHostIP() (string, error) {
	conn, err := net.Dial("udp", "192.192.192.192:80")
	if err != nil {
		return "127.0.0.1", err
	}
	defer func() {
		_ = conn.Close()
	}()

	var ip = strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip, nil
}
