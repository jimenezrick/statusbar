package main

import "flag"

func main() {
	flag.StringVar(&disk, "d", "sda", "disk device")
	flag.StringVar(&iface, "i", "eth0", "net interface")
	flag.IntVar(&pause, "p", 1, "update interval (seconds)")
	flag.Parse()

	go updater()
	updateStats()
}
