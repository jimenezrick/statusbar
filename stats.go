package main

import (
	"fmt"
	"time"
)

func updateStats() {
	for {
		s := formatStats(time.Now(), loadavg(), usedMem(), 0.2, 20, 300000000) // XXX
		select {
		case statsUpdates <- s:
		default:
			println("--- SKIP ---") // TODO: Remove
		}
		time.Sleep(time.Second)
	}
}

func formatStats(t time.Time, load float64, mem int, io float64, up, down int) string {
	return fmt.Sprintf("%s [%.2fL] [%s] [%.2fIO] [%s/%s]",
		t.Format(time.ANSIC), load, addUnits(mem), io, addUnits(up), addUnits(down))
}

func addUnits(x int) string {
	switch {
	case x >= 1024*1024*1024:
		return fmt.Sprintf("%.2fG", float64(x)/(1024*1024*1024))
	case x >= 1024*1024:
		return fmt.Sprintf("%.2fM", float64(x)/(1024*1024))
	case x >= 1024:
		return fmt.Sprintf("%.2fK", float64(x)/1024)
	}
	return fmt.Sprintf("%.2f", float64(x))
}

func loadavg() float64 {
	file := readFile("/proc/loadavg")
	return extractFloatCol(file, 1)
}

func usedMem() int {
	file := readFile("/proc/meminfo")
	memTotal := extractIntCol(extractLine(file, "MemTotal"), 2)
	memFree := extractIntCol(extractLine(file, "MemFree"), 2)
	return (memTotal - memFree) * 1024
}
