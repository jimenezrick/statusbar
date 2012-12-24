package main

import (
	"fmt"
	"time"
)

const (
	// TODO: Take from command line
	disk = "sda"
	iface = "eth0"
	pause = time.Second
)

type prevStats struct {
	v []int64
	t time.Time
}

var prevNetStats prevStats

func init() {
	snd, rcv := netSndRcv()
	prevNetStats = prevStats{[]int64{snd, rcv}, time.Now()}
}

func updateStats() {
	for {
		io := 0.0 // XXX
		up, down := netRate()
		stats := formatStats(time.Now(), loadAvg(), usedMem(), io, up, down)

		select {
		case statsUpdates <- stats:
		default:
			// Don't queue stale updates
			println("--- SKIP ---") // TODO: Remove
		}
		time.Sleep(pause)
	}
}

func formatStats(t time.Time, load float64, mem int64, io float64, up, down int64) string {
	return fmt.Sprintf(
		"%s [%.2fL] [%s] [%.2fIO] [%s/%s]",
		t.Format(time.ANSIC),
		load,
		withUnits(mem, 6),
		io,
		withUnits(up, 6),
		withUnits(down, 6),
	)
}

func withUnits(x int64, width int) string {
	switch {
	case x >= 1024*1024*1024:
		return fmt.Sprintf("%*.2fG", width, float64(x)/(1024*1024*1024))
	case x >= 1024*1024:
		return fmt.Sprintf("%*.2fM", width, float64(x)/(1024*1024))
	case x >= 1024:
		return fmt.Sprintf("%*.2fK", width, float64(x)/1024)
	}
	return fmt.Sprintf("%*.2fB", width, float64(x))
}

func loadAvg() float64 {
	file := readFile("/proc/loadavg")
	return extractFloatCol(file, 1)
}

func usedMem() int64 {
	file := readFile("/proc/meminfo")
	memTotal := extractIntCol(extractLine(file, "MemTotal"), 2)
	memFree := extractIntCol(extractLine(file, "MemFree"), 2)
	return (memTotal - memFree) * 1024
}













// io....
// ioRate













func netSndRcv() (int64, int64) {
	file := readFile("/proc/net/dev")
	snd := extractIntCol(extractLine(file, iface), 10)
	rcv := extractIntCol(extractLine(file, iface), 2)
	return snd, rcv
}

func netRate() (int64, int64) {
	snd, rcv := netSndRcv()
	now := time.Now()
	up := int64(float64(snd - prevNetStats.v[0]) / now.Sub(prevNetStats.t).Seconds())
	down := int64(float64(rcv - prevNetStats.v[1]) / now.Sub(prevNetStats.t).Seconds())
	prevNetStats = prevStats{[]int64{snd, rcv}, now}
	return up, down
}
