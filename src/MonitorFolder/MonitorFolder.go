package main

import (
	"flag"
	"fmt"

	"github.com/fsnotify/fsnotify"
)

var folder = flag.String("folder", "./", "Folder to monitor")

// main
func main() {

	flag.Parse()

	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR", err)
	}
	defer watcher.Close()

	//
	done := make(chan bool)

	//
	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\t%s\t%s\n", event, event.Op.String(), event.String())

				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := watcher.Add(*folder); err != nil {
		fmt.Println("ERROR", err)
	}

	<-done
}
