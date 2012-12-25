package main

import (
	"net"
	"bufio"
	"strings"
)

var address string

func listener() {
	defer recoverError()

	lis, err := net.Listen("tcp", address)
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

	r := bufio.NewReader(conn)
	fs := strings.Fields(readLine(r))

	switch {
	case len(fs) == 1 && fs[0] == "notification:":
		notifications <- readLine(r)
	case len(fs) == 2 && fs[0] == "status" && strings.HasSuffix(fs[1], ":"):
		for {
			select {
			case remoteStats <- fs[1] + " " + readLine(r):
			default:
				// Don't enqueue stale updates
			}
		}
	}

	if err := conn.Close(); err != nil {
		panic(err)
	}
}

func readLine(r *bufio.Reader) string {
	l, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return l[:len(l) - 1]
}
