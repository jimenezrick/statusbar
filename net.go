package main

import (
	"net"
	"bufio"
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
	switch readLine(r) {
	case "notification:":
		notifications <- readLine(r)
	case "status:":
		//
		// XXX: also help usage
		//
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
