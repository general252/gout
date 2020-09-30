package uhost

import (
	"github.com/general252/cpu_percent/cpu_percent"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"time"
)

type MemoryStat struct {
	UsedPercent float64
	Total       uint64
	Used        uint64
	Available   uint64
}

type CpuInfo struct {
	Name       string
	PhysicalID string
	VendorId   string
	Cores      int32
	MHz        float64
}
type DiskInfo struct {
	Device     string
	MountPoint string
	FSType     string

	TotalSize   uint64
	FreeSize    uint64
	UsedSize    uint64
	UsedPercent float64
}

// disk info
func DiskDeviceInfo() []DiskInfo {
	var diskInfoArray []DiskInfo

	if parts, err := disk.Partitions(true); err == nil {
		for _, part := range parts {
			diskInfo, _ := disk.Usage(part.Mountpoint)
			diskInfoArray = append(diskInfoArray, DiskInfo{
				Device:      part.Device,
				MountPoint:  part.Mountpoint,
				FSType:      part.Fstype,
				TotalSize:   diskInfo.Total,
				FreeSize:    diskInfo.Free,
				UsedSize:    diskInfo.Used,
				UsedPercent: diskInfo.UsedPercent,
			})
		}
	}
	return diskInfoArray
}

func MemInfo() MemoryStat {
	var memInfo MemoryStat

	if v, err := mem.VirtualMemory(); err == nil {
		memInfo.Total = v.Total
		memInfo.Available = v.Available
		memInfo.Used = v.Used
		memInfo.UsedPercent = v.UsedPercent
	}

	return memInfo
}

func CPUInfo() CpuInfo {
	var cpuInfo CpuInfo
	if v, err := cpu.Info(); err == nil {
		if len(v) > 0 {
			var info = v[0]
			cpuInfo.Name = info.ModelName
			cpuInfo.VendorId = info.VendorID
			cpuInfo.PhysicalID = info.PhysicalID
			cpuInfo.Cores = info.Cores
			cpuInfo.MHz = info.Mhz
		}
	}

	return cpuInfo
}

// example: CPUPercent(time.Second)
func CPUPercent(interval time.Duration) float64 {
	var cpuPercent float64
	//if percent, err := cpu.Percent(time.Second, false); err == nil {
	//	if len(percent) > 0 {
	//		cpuPercent = percent[0]
	//	}
	//}
	if percent, err := cpu_percent.GetCpuPercent(interval); err == nil {
		cpuPercent = percent
	}

	return cpuPercent
}
