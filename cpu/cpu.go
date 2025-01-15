package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpuInfo struct {
  ModelName     string
  CoresCount    int
  ThreadsCount  int
  CurrentFreq   string
  CacheSize     string
  Temperature   string
  // PowerUsage    string
}


type cpuTime struct {
  User    uint64
  System  uint64
  Idle    uint64
  Total   uint64
}


type CPUINFO struct {
  InfoCPU       cpuInfo
  TimeCPU       cpuTime
  TimePerThread map[string]*cpuTime 
  UsageCPU      string
  UsageThreads  map[string]float64
}

func GetCpuInfo() (*CPUINFO, error) {
  infoCpu, err := getCpuInfo()
  if err != nil {
    return nil,err
  }
  timeCpu, err := getCpuTime()
  if err != nil {
    return nil,err
  }
  timePerThread,err := getCpuTimePerThread()
  if err != nil {
    return nil,err
  }
  usageCpu, err := getCpuUsage()
  if err != nil {
    return nil,err
  }
  usageThreads, err := getThreadsUsage()
  if err != nil {
    return nil,err
  }
  return &CPUINFO{
    InfoCPU: *infoCpu,
    TimeCPU: *timeCpu,
    TimePerThread: timePerThread,
    UsageCPU: usageCpu,
    UsageThreads: usageThreads,
  },nil
}



func getCpuInfo() (*cpuInfo, error) {
  file, err := os.Open("/proc/cpuinfo")
  if err != nil {
    return nil,err
  }
  defer file.Close()

  var modelName, currentFreq, cacheSize, temperature, currentThread string
  var coreCount, threadsCount int
  var weightedFrequencySum float64
  var totalWeight float64
  freqPerThread := make(map[string]float64)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()
    if strings.HasPrefix(line, "model name") {
      modelName = strings.TrimSpace(strings.Split(line, ":")[1])
    } else if strings.HasPrefix(line, "cache size") {
			cacheSize = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.HasPrefix(line, "cpu cores") {
			coreCount, _ = strconv.Atoi(strings.TrimSpace(strings.Split(line, ":")[1]))
		} else if strings.HasPrefix(line, "processor") {
      threadsCount++
      currentThread = "cpu" + strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.HasPrefix(line, "cpu MHz") {
      threadFreq, _ := strconv.ParseFloat(strings.TrimSpace(strings.Split(line, ":")[1]),32)
      freqPerThread[currentThread] = threadFreq
    }
  }
  threadsUsage, err := getThreadsUsage()
  if err != nil {
    return nil, err
  }
  for thread, freq := range freqPerThread {
    weightedFrequencySum += freq * threadsUsage[thread]
    totalWeight += threadsUsage[thread]
  }
  currentFreq = strconv.FormatFloat((weightedFrequencySum/totalWeight),'f',2,64)
  temperature, err = getCpuTemp()
  // powerUsage, err = getCpuPowerUsage()
  return &cpuInfo {
    ModelName:    modelName,
    CoresCount:   coreCount,
    ThreadsCount: threadsCount,
    CurrentFreq:  currentFreq,
    CacheSize:    cacheSize,
    Temperature:  temperature,   
  // PowerUsage    string
  },nil
}
func findThermalZone() (string, error) {
	thermalZonesPath := "/sys/class/thermal/"

	files, err := os.ReadDir(thermalZonesPath)
	if err != nil {
		return "", err
	}
  
	for _, file := range files {
		if strings.Contains(strings.TrimSpace(file.Name()), "thermal") {
			thermalZonePath := thermalZonesPath + file.Name() + "/"
			typeFilePath := thermalZonePath + "type"
			typeFile, err := os.Open(typeFilePath)
			if err != nil {
				continue 
			}
			defer typeFile.Close()

			var zoneType string
			_, err = fmt.Fscanf(typeFile, "%s", &zoneType)
			if err != nil {
				continue
			}

			if zoneType == "x86_pkg_temp" {
				return thermalZonePath, nil
			}
		}
	}

	return "", fmt.Errorf("x86_pkg_temp not found in any thermal zone")
}

func getCpuTemp() (string, error) {
	thermalZonePath, err := findThermalZone()
	if err != nil {
		return "", err
	}

	tempFilePath := thermalZonePath + "temp"
	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	var tempInMC int
	_, err = fmt.Fscanf(tempFile, "%d", &tempInMC)
	if err != nil {
		return "", err
	}

	return strconv.FormatFloat(float64(tempInMC)/1000, 'f', 2, 64) + " ℃", nil
}



func getCpuTimePerThread() (map[string]*cpuTime,error) {
  file, err := os.Open("/proc/stat")
  if err != nil {
    return nil, err
  }
  defer file.Close()
  
  cpuTimePerThread := make(map[string]*cpuTime)
  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    line := scanner.Text()
    fields := strings.Fields(line)
    if len(fields) > 0 && strings.HasPrefix(fields[0], "cpu") && fields[0] != "cpu" {
      thread := fields[0]
      var user,system,idle,total uint64
      for i, val := range fields[1:]{
        num, err := strconv.ParseUint(val,10,64)
        if err != nil {
          return nil,err
        }
        total += num 
        switch i {
        case 0:
          user = num
        case 2:
          system = num
        case 3:
          idle  = num
        }
      }
      cpuTimePerThread[thread] = &cpuTime{
        User:   user,
        System: system,
        Idle:   idle,
        Total:  total,
      }
    }
  }
  return cpuTimePerThread,nil
}


func getCpuTime() (*cpuTime, error) {
  file, err := os.Open("/proc/stat")
  if err != nil {
    return nil,err
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Scan()
  firstLine := scanner.Text()
  fields := strings.Fields(firstLine)
  var total,user,system,idle uint64
  if len(fields) > 0 && fields[0] == "cpu" {
    for i, val := range fields[1:]{
      num, err := strconv.ParseUint(val,10,64)
      if err != nil {
        return nil,err
      }
      total += num 
      switch i {
      case 0:
        user = num
      case 2:
        system = num
      case 3:
        idle  = num
      }
    }
  }
  return &cpuTime {
    User:   user,
    System: system,
    Idle:   idle,
    Total:  total,
  },nil
}


func getCpuUsage() (cpuUsage string, err error) {
  prevCpuTime,err := getCpuTime()
  if err != nil {
    return "", err
  }
  prevTot := prevCpuTime.Total
  prevIdle := prevCpuTime.Idle
  time.Sleep(time.Second)
  currCpuTime, err := getCpuTime()
  if err != nil {
    return  "", err
  }
  currTot := currCpuTime.Total
  currIdle := currCpuTime.Idle
  return strconv.FormatFloat((1 - (float64(currIdle - prevIdle)/float64(currTot - prevTot))) * 100,'f',2,64) + "%", nil
}


func getThreadsUsage() (map[string]float64, error) {
  prevThreadsTime, err := getCpuTimePerThread()
  if err != nil {
    return nil,err
  }
  time.Sleep(time.Second)
  currThreadsTime, err := getCpuTimePerThread()
  if err != nil {
    return nil,err
  }
  threadsUsage := make(map[string]float64)
  for thread, prevTime := range prevThreadsTime {
    currTime, exists := currThreadsTime[thread]
    if !exists {
      return nil, fmt.Errorf("No current stats(time) for thread %s",thread)
    }
    prevTot := prevTime.Total
    prevIdle := prevTime.Idle
    currTot := currTime.Total
    currIdle := currTime.Idle
    threadsUsage[thread] = (1 - (float64(currIdle - prevIdle)/float64(currTot - prevTot))) * 100
  }
  return threadsUsage, nil
}

