package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

func sendNotification(host, msg string) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	writeLine(conn, "notify "+hostname+":")
	writeLine(conn, msg)
}

func listener(addr string) {
	defer recoverErrorExit()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer recoverError()
	defer conn.Close()

	lhost, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		panic(err)
	}
	rhost, _, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		panic(err)
	}

	r := bufio.NewReader(conn)
	fs := strings.Fields(readLine(r))

	switch {
	case len(fs) == 2 && fs[0] == "notify" && strings.HasSuffix(fs[1], ":"):
		if rhost == lhost {
			notifications <- readLine(r)
		} else {
			notifications <- fs[1] + " " + readLine(r)
		}
	case len(fs) == 2 && fs[0] == "status" && strings.HasSuffix(fs[1], ":"):
		for {
			select {
			case remoteStats <- fs[1] + " " + readLine(r):
			default:
				// Don't enqueue stale updates
			}
		}
	}
}

func readLine(r *bufio.Reader) string {
	l, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return l[:len(l)-1]
}

func writeLine(w io.Writer, s string) {
	if _, err := w.Write([]byte(s + "\n")); err != nil {
		panic(err)
	}
}
