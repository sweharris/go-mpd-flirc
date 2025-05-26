package main

// We use this to get the current active window name
// We need this so we know whether to ignore the IR command we've received
//
// Code based on https://github.com/DeedleFake/ptt-fix/blob/v0.7.5/internal/xdo/xdo.go
//
// tidied up, stripped to minimum.  It works for me on Debian 12.  Maybe it'll
// work for you!

/*
#include <stdlib.h>
#include <xdo.h>
#cgo LDFLAGS: -lxdo
*/
import "C"

import (
    "unsafe"
    "runtime"
)

type Xdo struct {
    x *C.xdo_t
}

type Window struct {
    w C.Window
    x *Xdo
}

func xdo_finalizer(x *Xdo) {
    C.xdo_free (x.x);
}

func NewXdo() (*Xdo, bool) {
    x := C.xdo_new (nil);
    if x == nil {
      return nil, false
    }
    r := &Xdo{x};
    runtime.SetFinalizer(r, xdo_finalizer);
    return r, true;
}

func (x *Xdo) free() {
	if x.x != nil {
		C.xdo_free(x.x)
		x.x = nil
	}
}

func (x *Xdo) GetFocusedWindow() Window {
    window := C.Window(0);
    C.xdo_get_focused_window_sane(x.x, &window)
    return Window{window, x};
}

func (w *Window) GetName() string {
    var name_ret        *C.uchar;
    var name_len_ret    C.int;
    var whatever        C.int;
    C.xdo_get_window_name(w.x.x, w.w, &name_ret, &name_len_ret, &whatever);
    str := C.GoBytes(unsafe.Pointer(name_ret), name_len_ret);
    C.free(unsafe.Pointer(name_ret));
    return string(str);
}

func current_window() string {
	x,e := NewXdo()
	if !e {
		return ""
	}
		
	w := x.GetFocusedWindow()
	n := w.GetName()
	x.free()

	return n
}
