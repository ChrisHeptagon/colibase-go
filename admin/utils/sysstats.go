package utils

import (
	"fmt"
	"os"

	"github.com/mackerelio/go-osstat/memory"
)

func GetSysStats() ([]map[string]interface{}, error) {
	memory, err := memory.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return nil, err
	}
	var stats []map[string]interface{}
	stats = append(stats, map[string]interface{}{
		"total": memory.Total,
	})
	stats = append(stats, map[string]interface{}{
		"used": memory.Used,
	})
	stats = append(stats, map[string]interface{}{
		"cached": memory.Cached,
	})
	stats = append(stats, map[string]interface{}{
		"free": memory.Free,
	})
	stats = append(stats, map[string]interface{}{
		"active": memory.Active,
	})
	stats = append(stats, map[string]interface{}{
		"inactive": memory.Inactive,
	})
	return stats, nil
}
