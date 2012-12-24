package main

import (
	"strings"
	"io/ioutil"
)

func readFile(name string) string {
	buf, err := ioutil.ReadFile(name)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func extractLine(s, substr string) string {
	lines := strings.Split(s, "\n")
	for _, l := range lines {
		if strings.Contains(l, substr) {
			return l
		}
	}
	panic("line not found")
}

func extractColumn(s string, n int) string {
	fields := strings.Fields(s)
	return fields[n - 1]
}
