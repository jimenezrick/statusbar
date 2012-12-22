package main

// #cgo CFLAGS: -Wall
// #cgo LDFLAGS: -lxcb
// #include <stdlib.h>
// #include "wm_name.h"
import "C"
import "unsafe"

func main() {
	xconn := C.connect_x()
	if xconn == nil {
		panic("can't connect to X server")
	}
	defer C.disconnect_x(xconn)

	if C.set_screen(xconn) == -1 {
		panic("can't find screen")
	}

	name := C.CString("ABC")
	if C.set_wm_name(xconn, name) == -1 { // XXX: Error checking
		panic("can't set window manager name")
	}
	defer C.free(unsafe.Pointer(name));
}
