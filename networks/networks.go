package networkinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type NetInfo struct {
	Name        string
	State       string
	Speed       int64 // interface cap
	RxBytes     uint64
	TxBytes     uint64
	RxErrors    uint64
	TxErrors    uint64
}

func readFileAsString(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %w", path, err)
	}
	return strings.TrimSpace(string(content)), nil
}

func readFileAsUint64(path string) (uint64, error) {
	content, err := readFileAsString(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(content, 10, 64)
}

func readFileAsInt64(path string) (int64, error) {
	content, err := readFileAsString(path)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(content, 10, 64)
}

func getNetworkInterfaces() ([]NetInfo, error) {
	var interfaces []NetInfo

	entries, err := os.ReadDir("/sys/class/net")
	if err != nil {
		return nil, fmt.Errorf("failed to read /sys/class/net: %w", err)
	}

	for _, entry := range entries {


		ifName := entry.Name()
		basePath := filepath.Join("/sys/class/net", ifName)

		state, err := readFileAsString(filepath.Join(basePath, "operstate"))
		if err != nil {
			fmt.Printf("Warning: Could not read state for %s: %v\n", ifName, err)
			state = "unknown"
		}

		speed, err := readFileAsInt64(filepath.Join(basePath, "speed"))
		if err != nil {
			fmt.Printf("Info: Speed not available for %s: %v\n", ifName, err)
			speed = -1
		}

		rxBytes, err := readFileAsUint64(filepath.Join(basePath, "statistics/rx_bytes"))
		if err != nil {
			fmt.Printf("Warning: Could not read RX bytes for %s: %v\n", ifName, err)
		}
		txBytes, err := readFileAsUint64(filepath.Join(basePath, "statistics/tx_bytes"))
		if err != nil {
			fmt.Printf("Warning: Could not read TX bytes for %s: %v\n", ifName, err)
		}
		rxErrors, err := readFileAsUint64(filepath.Join(basePath, "statistics/rx_errors"))
		if err != nil {
			fmt.Printf("Warning: Could not read RX errors for %s: %v\n", ifName, err)
		}
		txErrors, err := readFileAsUint64(filepath.Join(basePath, "statistics/tx_errors"))
		if err != nil {
			fmt.Printf("Warning: Could not read TX errors for %s: %v\n", ifName, err)
		}

		interfaces = append(interfaces, NetInfo{
			Name:     ifName,
			State:    state,
			Speed:    speed,
			RxBytes:  rxBytes,
			TxBytes:  txBytes,
			RxErrors: rxErrors,
			TxErrors: txErrors,
		})
	}

	return interfaces, nil
}
