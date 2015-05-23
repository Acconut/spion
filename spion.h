uv_loop_t* bdg_new_loop();
uv_fs_event_t* bdg_new_watcher(uv_loop_t* loop, void* callback);
int bdg_add_file(void* handle, char* filename);
