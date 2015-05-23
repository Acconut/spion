package main

import (
	"fmt"
	"time"

	"github.com/Acconut/spion"
)

func main() {
	watcher, err := spion.New("./")
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				fmt.Printf("Change detected in %s: %s\n", event.Path, event.Filename)
			case <-time.After(5 * time.Second):
				fmt.Println("Stopping...")
				if err := watcher.Stop(); err != nil {
					fmt.Println(err)
				}
				return
			}
		}
	}()

	if err := watcher.Start(); err != nil {
		fmt.Println(err)
		return
	}
}
