package unet

import (
	"github.com/shirou/gopsutil/net"
)

// PidListenPorts 进程监听的端口号
func PidListenPorts(pid int) ([]int, error) {
	var ports []int
	connectionStats, err := net.ConnectionsPid("all", int32(pid))
	if err != nil {
		return ports, err
	}

	for _, stat := range connectionStats {
		if stat.Status == "LISTEN" {
			ports = append(ports, int(stat.Laddr.Port))
		}
	}

	return ports, nil
}
