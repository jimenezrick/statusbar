package main

// XXX /proc/net/dev
// XXX /proc/diskstats

func main() {
	go updater()
	updateStats()
}
