package main

import (
	"flag"
	"fmt"
	"os"
)

func printError(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
}

func recoverError() {
	if err := recover(); err != nil {
		printError(err)
	}
}

func recoverErrorExit() {
	if err := recover(); err != nil {
		printError(err)
		os.Exit(1)
	}
}

func main() {
	defer recoverErrorExit()

	flag.IntVar(&updateInterval, "u", 1, "update interval (seconds)")
	flag.StringVar(&disk, "d", "sda", "disk device")
	flag.StringVar(&iface, "i", "eth0", "net interface")
	flag.StringVar(&address, "l", ":9000", "listen address ([host]:port)")
	flag.Parse()

	go updater()
	go listener()
	updateStats()
}
