package main

import "net"

var address string

func listener() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			printError(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	//
	// XXX
	//
}
