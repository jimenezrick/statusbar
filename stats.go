package main

import (
	"fmt"
	"time"
)

func updateStats() {
	for {
		s := formatStats(time.Now(), 0.025, usedMem(), 0.2, 20, 300000000) // XXX
		select {
		case statsUpdates <- s:
		default:
			println("--- SKIP ---") // TODO: Remove
		}
		time.Sleep(time.Second)
	}
}

func formatStats(t time.Time, load float32, mem int, io float32, up, down int) string {
	return fmt.Sprintf("%s [%.2fL] [%s] [%.2fIO] [%s/%s]",
		t.Format(time.ANSIC), load, addUnits(mem), io, addUnits(up), addUnits(down))
}

func addUnits(x int) string {
	switch {
	case x >= 1024*1024*1024:
		return fmt.Sprintf("%.2fG", float32(x)/(1024*1024*1024))
	case x >= 1024*1024:
		return fmt.Sprintf("%.2fM", float32(x)/(1024*1024))
	case x >= 1024:
		return fmt.Sprintf("%.2fK", float32(x)/1024)
	}
	return fmt.Sprintf("%.2f", float32(x))
}

func usedMem() int {
	file := readFile("/proc/meminfo")
	memTotal := extractIntColumn(extractLine(file, "MemTotal"), 2)
	memFree := extractIntColumn(extractLine(file, "MemFree"), 2)
	return (memTotal - memFree) * 1024
}
