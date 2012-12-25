package main

import (
	"io/ioutil"
	"strconv"
	"strings"
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
	panic("line with substring " + substr + " not found")
}

func extractCol(s string, n int) string {
	fields := strings.Fields(s)
	return fields[n-1]
}

func extractIntCol(s string, n int) int64 {
	i, err := strconv.ParseInt(extractCol(s, n), 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func extractFloatCol(s string, n int) float64 {
	f, err := strconv.ParseFloat(extractCol(s, n), 64)
	if err != nil {
		panic(err)
	}
	return f
}
