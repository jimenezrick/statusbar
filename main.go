package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/blockdevice"
)

const (
	KiB = 1 << 10
	MiB = 1 << 20
	GiB = 1 << 30

	SectorBytes = 512
)

var (
	interval = flag.Duration("u", time.Second, "update interval")
	disk     = flag.String("d", "sda", "disk device")
	iface    = flag.String("i", "eth0", "net interface")

	lastValues = make(map[string]float64)
)

type Info struct {
	now          time.Time
	cpu          float64
	mem          uint64
	rio, wio     uint64
	rxnet, txnet uint64
}

func printSystemInfo(info Info) {
	fmt.Printf(
		"system %s [L:%3.1f] [M:%s] [IO:%s/%s] [N:%s/%s]\n",
		info.now.Format(time.ANSIC),
		info.cpu,
		withUnits(uint(info.mem)),
		withUnits(uint(info.rio)),
		withUnits(uint(info.wio)),
		withUnits(uint(info.rxnet)),
		withUnits(uint(info.txnet)),
	)
}

func withUnits(bytes uint) string {
	var quantity float64
	var unit string

	switch {
	case bytes > GiB:
		quantity = float64(bytes) / float64(GiB)
		unit = "G"
	case bytes > MiB:
		quantity = float64(bytes) / float64(MiB)
		unit = "M"
	case bytes > KiB:
		quantity = float64(bytes) / float64(KiB)
		unit = "K"
	default:
		quantity = float64(bytes)
		unit = "B"
	}

	if unit == "B" {
		return fmt.Sprintf("%5dB", bytes)
	}
	return fmt.Sprintf("%5.1f%s", quantity, unit)
}

func getRate(key string, val float64) float64 {
	oldVal, ok := lastValues[key]
	lastValues[key] = val
	if !ok {
		return 0
	}

	return (val - oldVal) / float64((*interval).Seconds())
}

func main() {
	flag.Parse()

	proc, err := procfs.NewDefaultFS()
	if err != nil {
		log.Fatal(err)
	}

	block, err := blockdevice.NewDefaultFS()
	if err != nil {
		log.Fatal(err)
	}

	for {
		stat, err := proc.Stat()
		if err != nil {
			log.Fatal(err)
		}
		numCPUs := float64(len(stat.CPU))
		cpuIdle := getRate("cpu_idle", stat.CPUTotal.Idle)

		mem, err := proc.Meminfo()
		if err != nil {
			log.Fatal(err)
		}

		iostat, _, err := block.SysBlockDeviceStat(*disk)
		if err != nil {
			log.Fatal(err)
		}
		rbytes := uint64(getRate("r_bytes", float64(iostat.ReadSectors)))
		wbytes := uint64(getRate("w_bytes", float64(iostat.WriteSectors)))

		netstat, err := proc.NetDev()
		if err != nil {
			log.Fatal(err)
		}
		rxbytes := uint64(getRate("rx_bytes", float64(netstat[*iface].RxBytes)))
		txbytes := uint64(getRate("tx_bytes", float64(netstat[*iface].TxBytes)))

		info := Info{
			now:   time.Now(),
			cpu:   numCPUs - cpuIdle,
			mem:   (mem.MemTotal - mem.MemAvailable) * KiB,
			rio:   rbytes * SectorBytes,
			wio:   wbytes * SectorBytes,
			rxnet: rxbytes,
			txnet: txbytes,
		}
		printSystemInfo(info)
		time.Sleep(*interval)
	}
}
