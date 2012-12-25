package main

import (
	"os"
	"fmt"
	"flag"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], r)
		}
	}()

	flag.StringVar(&disk, "d", "sda", "disk device")
	flag.StringVar(&iface, "i", "eth0", "net interface")
	flag.IntVar(&pause, "p", 1, "update interval (seconds)")
	flag.Parse()

	go updater()
	updateStats()
}
