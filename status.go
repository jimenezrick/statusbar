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

const notificationPause = 5 * time.Second

var (
	notifications = make(chan string, 5)
	remoteStats   = make(chan string, 5)
	localStats    = make(chan string)
)

func updater(timeoutRemote int) {
	defer recoverErrorExit()

	xconn := C.connect_x()
	if xconn == nil {
		panic("can't connect to X server")
	}
	defer C.disconnect_x(xconn)

	if C.set_screen(xconn) == -1 {
		panic("can't find screen")
	}

	for {
		select {
		case s := <-notifications:
			notify(xconn, s)
		case s := <-remoteStats:
			set_wm_name(xconn, s)
			priorizeRemoteStats(xconn, timeoutRemote)
		case s := <-localStats:
			set_wm_name(xconn, s)
		}
	}
}

func priorizeRemoteStats(xconn *C.xconn_t, timeoutRemote int) {
	for {
		t := time.NewTimer(time.Duration(timeoutRemote) * time.Second)

		select {
		case s := <-notifications:
			notify(xconn, s)
		case s := <-remoteStats:
			set_wm_name(xconn, s)
		case <-t.C:
			return
		}
		t.Stop()
	}
}

func set_wm_name(xconn *C.xconn_t, name string) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	if C.set_wm_name(xconn, str) == -1 {
		panic("can't set window manager name")
	}
}

func warn(xconn *C.xconn_t, msg, pattern string, len int, pause int, times int) {
	for t := 0; t < times; t++ {
		for l := 0; l < len; l++ {
			animation := strings.Repeat(" ", l) + pattern + strings.Repeat(" ", len-l)
			set_wm_name(xconn, animation+msg)
			time.Sleep(time.Millisecond * time.Duration(pause))
		}
	}
}

func notify(xconn *C.xconn_t, msg string) {
	p := time.After(notificationPause)
	warn(xconn, msg, "> > >", 20, 30, 3)
	set_wm_name(xconn, msg)
	<-p
}
