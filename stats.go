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

type netStats struct {
	lastTx, lastRcv int64
	lastTime time.Time
}

var nStats netStats

func init() {
	tx, rcv := txRcvNet()
	nStats = netStats{tx, rcv, time.Now()}
}

func updateStats() {
	for {
		io := 0.0 // XXX
		up, down := upDownNet()
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
	return fmt.Sprintf("%s [%.2fL] [%s] [%.2fIO] [%s/%s]",
		t.Format(time.ANSIC), load, withUnits(mem), io, withUnits(up), withUnits(down))
}

func withUnits(x int64) string {
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

func txRcvNet() (int64, int64) {
	file := readFile("/proc/net/dev")
	tx := extractIntCol(extractLine(file, iface), 10)
	rcv := extractIntCol(extractLine(file, iface), 2)
	return tx, rcv
}

func upDownNet() (int64, int64) {
	tx, rcv := txRcvNet()
	now := time.Now()
	up := int64(float64(tx - nStats.lastTx) / now.Sub(nStats.lastTime).Seconds())
	down := int64(float64(rcv - nStats.lastRcv) / now.Sub(nStats.lastTime).Seconds())
	nStats = netStats{tx, rcv, now}
	return up, down
}
