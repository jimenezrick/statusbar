package main

import (
	"flag"
	"time"
)

var (
	interval = flag.Int("u", 1, "update interval (seconds)")
	disk     = flag.String("d", "sda", "disk device")
	iface    = flag.String("i", "eth0", "net interface")
	ac       = flag.String("a", "AC", "ac device")
	bat      = flag.String("b", "BAT0", "battery device")
)

func main() {
	flag.Parse()

	for {
		printStats(*interval, *disk, *iface, *ac, *bat)
		time.Sleep(time.Second)
	}
}
