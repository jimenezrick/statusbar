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

	interval := flag.Int("u", 1, "update interval (seconds)")
	disk := flag.String("d", "sda", "disk device")
	iface := flag.String("i", "eth0", "net interface")
	addr := flag.String("l", ":9000", "listen on address ([host]:port)")
	host := flag.String("h", "localhost:9000", "connect to host ([host]:port)")
	msg := flag.String("n", "", "message notification")
	flag.Parse()

	if *msg != "" {
		sendNotification(*host, *msg)
	} else {
		go updater()
		go listener(*addr)
		updateStats(*interval, *disk, *iface)
	}
}
