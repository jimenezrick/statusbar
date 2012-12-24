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
	lastTx, lastRcv int
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

func loadAvg() float64 {
	file := readFile("/proc/loadavg")
	return extractFloatCol(file, 1)
}

func usedMem() int {
	file := readFile("/proc/meminfo")
	memTotal := extractIntCol(extractLine(file, "MemTotal"), 2)
	memFree := extractIntCol(extractLine(file, "MemFree"), 2)
	return (memTotal - memFree) * 1024
}

func txRcvNet() (int, int) {
	file := readFile("/proc/net/dev")
	tx := extractIntCol(extractLine(file, iface), 10)
	rcv := extractIntCol(extractLine(file, iface), 2)
	return tx, rcv
}

func upDownNet() (int, int) {
	tx, rcv := txRcvNet()
	now := time.Now()
	up := int(float64(tx - nStats.lastTx) / now.Sub(nStats.lastTime).Seconds())
	down := int(float64(rcv - nStats.lastRcv) / now.Sub(nStats.lastTime).Seconds())
	nStats = netStats{tx, rcv, now}
	return up, down
}
