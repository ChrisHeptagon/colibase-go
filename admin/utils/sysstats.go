package utils

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func GetStats() (map[string]map[string]interface{}, error) {
	var stats = make(map[string]map[string]interface{})
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	cpu, err := cpu.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	var memoryStats = make(map[string]interface{})
	var cpuStats = make(map[string]interface{})
	memoryStats["total"] = memory.Total
	memoryStats["used"] = memory.Used
	memoryStats["cached"] = memory.Cached
	memoryStats["free"] = memory.Free
	memoryStats["active"] = memory.Active
	memoryStats["inactive"] = memory.Inactive
	memoryStats["swap_total"] = memory.SwapTotal
	cpuStats["total"] = cpu.Total
	cpuStats["user"] = cpu.User
	cpuStats["system"] = cpu.System
	cpuStats["idle"] = cpu.Idle
	cpuStats["nice"] = cpu.Nice
	cpuStats["irq"] = cpu.Irq
	cpuStats["softirq"] = cpu.Softirq
	cpuStats["steal"] = cpu.Steal
	cpuStats["guest"] = cpu.Guest
	cpuStats["guestnice"] = cpu.GuestNice
	stats["memory"] = memoryStats
	stats["cpu"] = cpuStats
	return stats, err
}
