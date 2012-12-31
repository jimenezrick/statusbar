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

	// XXX XXX XXX: iface, disk, interval
	flag.IntVar(&updateInterval, "u", 1, "update interval (seconds)")
	flag.StringVar(&disk, "d", "sda", "disk device")
	flag.StringVar(&iface, "i", "eth0", "net interface")
	addr := flag.String("l", ":9000", "listen address ([host]:port)")
	// XXX XXX XXX

	host := flag.String("h", "localhost:9000", "connect to host ([host]:port)")
	msg := flag.String("n", "", "notification message")
	flag.Parse()

	if *msg != "" {
		sendNotification(*host, *msg)
	} else {
		go updater()
		go listener(*addr)
		updateStats()
	}
}
