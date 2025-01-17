package networks

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

func GetNetworkInterfaces() ([]NetInfo, error) {
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
			state = "unknown"
		}

		speed, err := readFileAsInt64(filepath.Join(basePath, "speed"))
		if err != nil {
			speed = -1
		}

		rxBytes, _ := readFileAsUint64(filepath.Join(basePath, "statistics/rx_bytes"))
		txBytes, _ := readFileAsUint64(filepath.Join(basePath, "statistics/tx_bytes"))
		rxErrors, _ := readFileAsUint64(filepath.Join(basePath, "statistics/rx_errors"))
		txErrors, _ := readFileAsUint64(filepath.Join(basePath, "statistics/tx_errors"))
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
