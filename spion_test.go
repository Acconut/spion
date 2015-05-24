package spion

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestNotExist(t *testing.T) {
	watcher, err := New("./does-not-exist")
	if watcher != nil {
		t.Error("expected no watcher")
	}
	if err != os.ErrNotExist {
		t.Error("expected ErrNotExist to be returned")
	}
}

func TestWatcher(t *testing.T) {
	// Create temporary directory
	dirName, err := ioutil.TempDir("", "spion-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirName)

	watcher, err := New(dirName)
	if err != nil {
		t.Fatal(err)
	}

	// This goroutine will receive the events emitted by the watcher and test
	// the values.
	go func() {
		for i := 0; i < 3; i++ {
			event := <-watcher.Event
			if event.Path != dirName {
				t.Errorf("%d: expected path to be %s but got %s", i, dirName, event.Path)
			}
			if event.Filename != "file.txt" {
				t.Errorf("%d: expected filename to be file.txt but got %s", i, event.Filename)
			}
		}

		fmt.Println("stopping")
		if err := watcher.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	// This goroutine causes the watcher to emit events by creating, writing to
	// and removing files.
	go func() {
		// Create a new file in the directory
		file, err := os.Create(dirName + "/file.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		// Write to a file
		if _, err := file.Write([]byte("hello world")); err != nil {
			t.Fatal(err)
		}

		// Delete the file
		if err := os.Remove(dirName + "/file.txt"); err != nil {
			t.Fatal(err)
		}
	}()

	if err := watcher.Start(); err != nil {
		t.Fatal(err)
	}
}
