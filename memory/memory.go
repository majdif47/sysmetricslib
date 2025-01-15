package memory

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type MemInfo struct {
  TotalMemory float64
  UsedMemory  float64
  FreeMemory  float64
  TotalSwap   float64
  UsedSwap    float64
  FreeSwap    float64
}

func GetMemoryStats() (*MemInfo, error) {
  file, err := os.Open("/proc/meminfo")
  if err != nil {
    return nil, err
  }
  var memTotal, memUsed, memFree, swapTotal, swapUsed, swapFree float64
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    
    val, err := strconv.ParseFloat(fields[1],64)
    if err != nil {
      return nil,err
    }
    switch fields[0] {
    case "MemTotal:":
      memTotal = val
		case "MemAvailable:":
			memFree = val
		case "SwapTotal:":
			swapTotal = val
		case "SwapFree:":
			swapFree = val
    }
  }
  swapUsed = swapTotal- swapFree
  memUsed = memTotal - memFree
  return &MemInfo{
    TotalMemory: ((memTotal/1024)/1024),
    UsedMemory: ((memUsed/1024)/1024),
    FreeMemory: ((memFree/1024)/1024),
    TotalSwap: ((swapTotal/1024)/1024),
    UsedSwap: ((swapUsed/1024)/1024),
    FreeSwap: ((swapFree/1024)/1024),
  },nil

}



