package uhost

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/general252/cpu_percent/cpu_percent"
	"github.com/general252/gout/ushell"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
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

type SystemInfo struct {
	Name        string // 系统名称
	Version     string // 系统版本
	Id          string // 系统ID
	InstallDate string // 系统安装时间
	StartDate   string // 系统启动时间
	BIOS        string // BIOS 版本
}

// GetSystemInfo get system info
func GetSystemInfo() (*SystemInfo, error) {
	if runtime.GOOS == "windows" {
		return getSystemInfoWindows()
	} else {
		if rs, err := getSystemInfoLinux(); err == nil {
			return rs, nil
		}
		if rs, err := getSystemInfoLinuxProc(); err == nil {
			return rs, nil
		}
		return nil, fmt.Errorf("get fail")
	}
}

func getSystemInfoWindows() (*SystemInfo, error) {
	stdOut, _, err := ushell.ShellCommand("systeminfo")
	if err != nil {
		return nil, err
	}

	var lines []string
	var r = bufio.NewReader(strings.NewReader(stdOut))
	for i := 0; i < 18; i++ {
		if data, _, err := r.ReadLine(); err != nil {
			break
		} else {
			if i == 0 || i == 8 || i == 15 || i == 16 {
				continue
			}

			v := strings.SplitN(string(data), ":", 2)
			if len(v) == 2 {
				lines = append(lines, strings.TrimSpace(v[1]))
			}
		}
	}

	if lines == nil || len(lines) < 14 {
		return nil, fmt.Errorf("get fail")
	}

	var info = &SystemInfo{
		Name:        lines[1],
		Version:     lines[2],
		Id:          lines[7],
		InstallDate: lines[8],
		StartDate:   lines[9],
		BIOS:        lines[13],
	}
	return info, nil
}

func getSystemInfoLinux() (*SystemInfo, error) {
	stdOut, _, err := ushell.ShellCommand("systeminfo")
	if err != nil {
		return nil, err
	}

	var lines []string
	var r = bufio.NewReader(strings.NewReader(stdOut))
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}

		v := strings.SplitN(string(line), ":", 2)
		if len(v) == 2 {
			lines = append(lines, strings.TrimSpace(v[1]))
		}

		lines = append(lines)
	}

	if lines == nil || len(lines) < 5 {
		return nil, fmt.Errorf("get fail")
	}

	var info = &SystemInfo{
		Name:        lines[1],
		Version:     lines[3],
		Id:          "",
		InstallDate: "",
		StartDate:   "",
		BIOS:        "",
	}

	return info, nil
}

func getSystemInfoLinuxProc() (*SystemInfo, error) {
	data, err := ioutil.ReadFile("/proc/version")
	if err != nil {
		return nil, err
	}

	var procInfo = strings.ToLower(string(data))
	if strings.Contains(procInfo, "red hat") {
		if dataVer, err := ioutil.ReadFile("/etc/centos-release"); err == nil {
			return &SystemInfo{
				Name:    "",
				Version: string(dataVer),
			}, nil
		}
	} else if strings.Contains(procInfo, "ubuntu") {
		a := strings.Index(procInfo, "ubuntu")
		if a >= 0 {
			b := strings.Index(procInfo[a:], ")")
			if b >= 0 {
				return &SystemInfo{
					Name:    "",
					Version: procInfo[a : a+b],
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("get fail")
}
