package main

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -lxcb
// #include <stdlib.h>
// #include "wm_name.h"
import "C"

import (
	"strings"
	"time"
	"unsafe"
)

var (
	notifications = make(chan string, 5)
	statsUpdates  = make(chan string)
)

func updater() {
	defer recoverError()

	xconn := C.connect_x()
	if xconn == nil {
		panic("can't connect to X server")
	}
	defer C.disconnect_x(xconn)

	if C.set_screen(xconn) == -1 {
		panic("can't find screen")
	}

	for {
		var s string

		select {
		case s = <-notifications:
			warn(xconn, ">", 20, 15, 2) // XXX: Todo sleep 5 after a notification
		case s = <-statsUpdates:
		}
		set_wm_name(xconn, s)
	}
}

func set_wm_name(xconn *C.xconn_t, name string) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	if C.set_wm_name(xconn, str) == -1 {
		panic("can't set window manager name")
	}
}

func warn(xconn *C.xconn_t, pattern string, len int, pause int, times int) {
	for t := 0; t < times; t++ {
		for l := len; l > 0; l-- {
			set_wm_name(xconn, strings.Repeat(pattern, l))
			time.Sleep(time.Duration(pause) * time.Millisecond)
		}
	}
}
