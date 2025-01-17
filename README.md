# sysmetricslib

**sysmetricslib** is a Go library designed to retrieve and monitor system metrics, including CPU, memory, disk, and network statistics. This library provides a straightforward API for developers to integrate system monitoring capabilities into their applications.

## Features

- **CPU Metrics**: Monitor CPU usage and performance.
- **Memory Metrics**: Track memory consumption and availability.
- **Disk Metrics**: Observe disk usage and I/O statistics.
- **Network Metrics**: Analyze network throughput and activity.

## Installation

To include **sysmetricslib** in your Go project, run:

```bash
go get github.com/majdif47/sysmetricslib
```
**Usage Example: 
```bash
package main

import (
    "fmt"
    "github.com/majdif47/sysmetricslib/cpu" 
    "github.com/majdif47/sysmetricslib/memory" 
)

func main() {
    // Retrieve CPU metrics
    cpuMetrics, err := cpu.GetMetrics()
    if err != nil {
        fmt.Println("Error retrieving CPU metrics:", err)
        return
    }
    fmt.Printf("CPU Usage: %.2f%%\n", cpuMetrics.Usage)

    // Retrieve Memory metrics
    memMetrics, err := memory.GetMetrics()
    if err != nil {
        fmt.Println("Error retrieving memory metrics:", err)
        return
    }
    fmt.Printf("Memory Usage: %.2f%%\n", memMetrics.Usage)
}
```
This example demonstrates how to fetch and display CPU and memory usage metrics using the sysmetricslib library.
