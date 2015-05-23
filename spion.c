#include <stdlib.h>
#include <stdio.h>

#include <uv.h>

#include "spion.h"
#include "_cgo_export.h"

void on_change(uv_fs_event_t* handle, const char* filename, int events, int status) {
	char path[1024];
	size_t size = 1023;
	// Does not handle error if path is longer than 1023.
	uv_fs_event_getpath(handle, path, &size);
	path[++size] = '\0';

	goCallback(handle->data, path, filename, events);
}

uv_loop_t* bdg_new_loop() {
	uv_loop_t* loop = malloc(sizeof(uv_loop_t));
	uv_loop_init(loop);
	return loop;
}

uv_fs_event_t* bdg_new_watcher(uv_loop_t* loop, void* callback) {
	uv_fs_event_t* handle = malloc(sizeof(uv_fs_event_t));
	uv_fs_event_init(loop, handle);

	handle->data = callback;

	return handle;
}

int bdg_add_file(void* handle, char* filename) {
	return uv_fs_event_start((uv_fs_event_t*) handle, on_change, filename, 0);
}
