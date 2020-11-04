package unet

import (
	"fmt"
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

// GetMacAddress 获取MAC地址
func GetMacAddress() (string, error) {
	var ip, err = GetHostIP()
	if err != nil {
		return "", err
	}
	interFaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	var findInterFace = func(ip string) (*net.Interface, error) {
		for _, face := range interFaces {
			if addrList, err := face.Addrs(); err == nil {
				for _, addr := range addrList {
					if strings.HasPrefix(addr.String(), ip) {
						return &face, nil
					}
				}
			}
		}
		return nil, fmt.Errorf("not found")
	}

	if face, err := findInterFace(ip); err != nil {
		return "", err
	} else {
		return face.HardwareAddr.String(), nil
	}
}
