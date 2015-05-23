# spion
File watcher for Golang based on libuv.

The difference to [fsnotify](https://github.com/go-fsnotify/fsnotify) is that spion
uses libuv's file watching methods. All the platform specific code is therefore
handled by libuv and not spion itself. fsnotify does not require a binding to C.
In addition spion should not be considered stable yet, there are some known
issues at the moment.

# Installation

```bash
go get github.com/Acconut/spion
```

Running spion requires libuv 1.x to be installed, version 0.10 is not supported.
The easiest way to obtain the shared object is by compiling libuv on your own.
By default `libuv.so.1` will be installed to `/usr/local/lib`, so ensure that
your linker looks there (or tell it using `LD_LIBRARY_PATH`).

# Support

Thanks to libuv, spion uses inotify on Linux, FSEvents on Darwin, kqueue on
BSDs, ReadDirectoryChangesW on Windows, event ports on Solaris but is
unsupported on Cygwin. For more information, consult the
[uvbook](http://nikhilm.github.io/uvbook/filesystem.html#file-change-events).

# Documentation

[![GoDoc](https://godoc.org/github.com/Acconut/spion?status.svg)](https://godoc.org/github.com/Acconut/spion)

# License

The MIT License (MIT)

Copyright (c) 2015 Marius Kleidl

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
