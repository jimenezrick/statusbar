package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

type prevStats struct {
	v []int64
	t time.Time
}

var prevNetStats, prevIoStats prevStats

func initStats(disk, iface string) {
	snd, rcv := netSndRcv(iface)
	busy := ioBusy(disk)
	now := time.Now()
	prevNetStats = prevStats{[]int64{snd, rcv}, now}
	prevIoStats = prevStats{[]int64{busy}, now}
}

func updateStats(interval int, disk, iface string) {
	initStats(disk, iface)
	for {
		io := ioRate(disk)
		up, down := netRate(iface)
		stats := formatStats(time.Now(), loadAvg(), usedMem(), io, up, down)

		select {
		case localStats <- stats:
		default:
			// Don't enqueue stale updates
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func sendStats(host string, interval int, disk, iface string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	writeLine(conn, "status "+hostname+":")

	initStats(disk, iface)
	for {
		io := ioRate(disk)
		up, down := netRate(iface)
		stats := formatStats(time.Now(), loadAvg(), usedMem(), io, up, down)

		writeLine(conn, stats)
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

func formatStats(t time.Time, load float64, mem int64, io float64, up, down int64) string {
	return fmt.Sprintf(
		"%s [%.1fL] [%s] [%5.1fIO] [%s/%s]",
		t.Format(time.ANSIC),
		load,
		withUnits(mem, 5),
		io,
		withUnits(up, 5),
		withUnits(down, 5),
	)
}

func withUnits(x int64, width int) string {
	switch {
	case x >= 1000*1000*1000:
		return fmt.Sprintf("%*.1fG", width, float64(x)/(1024*1024*1024))
	case x >= 1000*1000:
		return fmt.Sprintf("%*.1fM", width, float64(x)/(1024*1024))
	case x >= 1000:
		return fmt.Sprintf("%*.1fK", width, float64(x)/1024)
	}
	return fmt.Sprintf("%*.1fB", width, float64(x))
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

func ioBusy(disk string) int64 {
	file := readFile("/proc/diskstats")
	return extractIntCol(extractLine(file, disk), 13)
}

func ioRate(disk string) float64 {
	busy := ioBusy(disk)
	now := time.Now()
	rate := float64(busy-prevIoStats.v[0]) / 1000 / now.Sub(prevIoStats.t).Seconds()
	prevIoStats = prevStats{[]int64{busy}, now}
	return rate * 100
}

func netSndRcv(iface string) (int64, int64) {
	file := readFile("/proc/net/dev")
	snd := extractIntCol(extractLine(file, iface), 10)
	rcv := extractIntCol(extractLine(file, iface), 2)
	return snd, rcv
}

func netRate(iface string) (int64, int64) {
	snd, rcv := netSndRcv(iface)
	now := time.Now()
	up := int64(float64(snd-prevNetStats.v[0]) / now.Sub(prevNetStats.t).Seconds())
	down := int64(float64(rcv-prevNetStats.v[1]) / now.Sub(prevNetStats.t).Seconds())
	prevNetStats = prevStats{[]int64{snd, rcv}, now}
	return up, down
}
