package utils

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
)

func GetStats() ([]map[string]interface{}, error) {
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
	var stats []map[string]interface{}
	var memoryStats []map[string]interface{}
	var cpuStats []map[string]interface{}
	memoryStats = append(memoryStats, map[string]interface{}{
		"total": memory.Total,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"used": memory.Used,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"cached": memory.Cached,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"free": memory.Free,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"active": memory.Active,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"inactive": memory.Inactive,
	})
	memoryStats = append(memoryStats, map[string]interface{}{
		"swap_total": memory.SwapTotal,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"total": cpu.Total,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"user": cpu.User,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"system": cpu.System,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"idle": cpu.Idle,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"nice": cpu.Nice,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"irq": cpu.Irq,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"softirq": cpu.Softirq,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"steal": cpu.Steal,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"guest": cpu.Guest,
	})
	cpuStats = append(cpuStats, map[string]interface{}{
		"guestnice": cpu.GuestNice,
	})
	stats = append(stats, map[string]interface{}{
		"memory": memoryStats,
		"cpu":    cpuStats,
	})

	return stats, nil
}
