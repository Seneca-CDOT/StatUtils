package perfstats

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// AppVersion - current app version
const AppVersion = "0.4.0"

// StatEntry specific performance data sweep
type StatEntry struct {
	Stats interface{} `json:"stats"`
	Date  string      `json:"date"`
}

// SysStat Performance statistics
type SysStat struct {
	System interface{} `json:"system"` // platform details
	CPU    StatEntry   `json:"cpu"`    // CPU utilization
	Disk   StatEntry   `json:"disk"`   // Space for disk(s)
	Memory StatEntry   `json:"memory"` // Space for memory
}

type statGetter func() (StatEntry, error)

// Web API date format (utc timestamp)
func getDateFormatted() string {
	return time.Now().UTC().Format("20060102150405")
}

// Lowers only the first character in a string
func lowerFirst(str string) string {
	return strings.ToLower(string(str[0])) + str[1:]
}

// GetPlatformInfo get platform details
func GetPlatformInfo() (interface{}, error) {

	hostname, err := os.Hostname()
	machineDetails := map[string]string{
		"platform":        runtime.GOOS,
		"machine":         hostname,
		"softwareVersion": AppVersion,
	}
	return machineDetails, err
}

// PlatformSysStats Query performance stats
// (calls either _linux.go or _windows perfstats implementations
// depending on the platform it's running on)
func PlatformSysStats() (interface{}, error) {

	// find out current stats:
	statEntries := make(map[string]StatEntry)
	perfStatGetters := map[string]statGetter{
		"memory": getMemoryStats,
		"cpu":    getCPUStats,
		"disk":   getDiskStats,
	}

	for statName, getStatsFunc := range perfStatGetters {
		stats, err := getStatsFunc()
		if err != nil {
			return nil, fmt.Errorf("Cannot get %s details: %s", statName, err)
		}
		stats.Date = getDateFormatted()
		statEntries[statName] = stats
	}

	// query platform details:
	platformDetails, err := GetPlatformInfo()
	if err != nil {
		return nil, fmt.Errorf("Cannot get system platform details %s", err)
	}

	return SysStat{
		System: platformDetails,
		Memory: statEntries["memory"],
		CPU:    statEntries["cpu"],
		Disk:   statEntries["disk"],
	}, nil
}
