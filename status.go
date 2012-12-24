package main

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -lxcb
// #include <stdlib.h>
// #include "wm_name.h"
import "C"

import (
	"unsafe"
	"strings"
	"time"
)

var statsUpdates chan string = make(chan string)

func updater() {
	xconn := C.connect_x()
	if xconn == nil {
		panic("can't connect to X server")
	}
	defer C.disconnect_x(xconn)

	if C.set_screen(xconn) == -1 {
		panic("can't find screen")
	}

	warn(xconn, ">", 20, 25, 2)
	for {
		s := <-statsUpdates
		set_wm_name(xconn, s)
	}
}

func set_wm_name(xconn *C.xconn_t, name string) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str));
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
