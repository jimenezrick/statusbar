package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

const defaultPort = ":9000"

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
	addr := flag.String("l", defaultPort, "listen on address ([host]:port)")
	host := flag.String("h", "", "connect to host (host[:port])")
	msg := flag.String("n", "", "message notification")
	flag.Parse()

	if !strings.Contains(*host, ":") {
		*host += defaultPort
	}

	switch {
	case *msg != "" && *host == "":
		notifications <- *msg
	case *msg != "" && *host != "":
		sendNotification(*host, *msg)
	case *host != "":
		sendStats(*host, *interval, *disk, *iface)
	default:
		go updater()
		go listener(*addr)
		updateStats(*interval, *disk, *iface)
	}
}
