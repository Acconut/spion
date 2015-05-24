package spion

// #cgo pkg-config: libuv
// #include <stdlib.h> // free
// #include <uv.h> // uv_strerror, uv_fs_event_t, uv_loop_t
// #include "spion.h"
import "C"

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

// This error is returned if you try to reuse a watcher, Even after stopping
// a watcher it can not be restarted.
var ErrInUse = errors.New("spion: watcher already in use")

// Event represents a single event for a single file
type Event struct {
	// The name of the file which changed
	Filename string
	// The path where the changed file is or was located
	Path string
}

// A Watcher represents an active or stopped watcher which reports events using
// the Event channel.
type Watcher struct {
	Event   chan Event
	Stopped bool
	Started bool

	handle   *C.uv_fs_event_t
	loop     *C.uv_loop_t
	callback callbackFunc
}

type callbackFunc *func(path *C.char, filename *C.char, events C.int)

// New creates a watcher. If the paramater points to a file the single file will
// be watched. If it points to a directory will be watched but not recursivly
// (sub-directories are ignored).
func New(filename string) (*Watcher, error) {
	eventChan := make(chan Event)

	callback := func(path *C.char, filename *C.char, events C.int) {
		event := Event{
			Filename: C.GoString(filename),
			Path:     C.GoString(path),
		}
		eventChan <- event
	}

	// Create new loop for each watcher. After the watcher is stopped, the loop
	// will be closed and freed.
	loop := C.bdg_new_loop()
	// Inititate a new handle for each watcher. Handles can not be reused and have
	// to be freed after stopping. Reusing a handle causes and EINVAL (argument invalid)
	// error to be returned.
	handle := C.bdg_new_watcher(loop, unsafe.Pointer(&callback))

	str := C.CString(filename)
	defer C.free(unsafe.Pointer(str))

	if err := uvErr(C.bdg_add_file(unsafe.Pointer(handle), str)); err != nil {
		return nil, err
	}

	watcher := &Watcher{
		Event:    eventChan,
		handle:   handle,
		loop:     loop,
		callback: &callback,
	}

	return watcher, nil
}

// Start the watcher. This call will block as long as the watcher is active and
// will only return after Stop() is called. Even after stopping the watcher,
// the structure can not be reused. A new watcher must be created,
func (w *Watcher) Start() error {
	if w.Started {
		return ErrInUse
	}

	w.Started = true
	return uvErr(C.uv_run(w.loop, C.UV_RUN_DEFAULT))
}

// Stop the watcher and clean up internal structures. If the watcher has not
// been started yet or has already been stopped, this function call will do
// nothing.
func (w *Watcher) Stop() error {
	if !w.Started || w.Stopped {
		return nil
	}

	w.Stopped = true

	// Free loop and FS event handler
	defer C.free(unsafe.Pointer(w.loop))
	defer C.free(unsafe.Pointer(w.handle))

	// Remove the handler from the loop to cause uv_run() to return
	if err := uvErr(C.uv_fs_event_stop(w.handle)); err != nil {
		return err
	}
	C.uv_stop(w.loop)
	return nil
}

// UvError represents an error returned by libuv. This struct is returned if
// spion is not able to translate an UvError into an equivalent error in Go's
// standart library.
type UvError struct {
	Code int
}

func (err UvError) Error() string {
	message := C.GoString(C.uv_strerror(C.int(err.Code)))
	return fmt.Sprintf("uv_error: %s (%d)", message, err.Code)
}

//export goCallback
func goCallback(ptr unsafe.Pointer, path *C.char, filename *C.char, events C.int) {
	(*(callbackFunc)(ptr))(path, filename, events)
}

func uvErr(code C.int) error {
	// libuv's ENOENT error is the equivalent of Go's os.ErrNotExist
	if code == -2 {
		return os.ErrNotExist
	}

	if code < 0 {
		return UvError{
			Code: int(code),
		}
	}
	return nil
}
