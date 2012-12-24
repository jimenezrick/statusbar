package main

// XXX /proc/diskstats

func main() {
	go updater()
	updateStats()
}
