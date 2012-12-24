package main

// XXX /proc/meminfo
// XXX /proc/net/dev
// XXX /proc/diskstats

func main() {
	go updater()
	updateStats()
}
