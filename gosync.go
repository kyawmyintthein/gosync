package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"gosync/local"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

// An atomic counter
type counter struct {
	val int32
}

func (c *counter) increment() {
	atomic.AddInt32(&c.val, 1)
}

func (c *counter) value() int32 {
	return atomic.LoadInt32(&c.val)
}

func (c *counter) reset() {
	atomic.StoreInt32(&c.val, 0)
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// var testFile = "temp/boo/hello.csv"
	var testDir = "temp"

	if err := watcher.Watch(testDir); err != nil {
		log.Fatal("watcher.Watch(%q) failed: %s", testDir, err)
	}

	//fileList := []string{}
	err = filepath.Walk(testDir, func(path string, f os.FileInfo, err error) error {
		watcher.Watch(path)
		//  fileList = append(fileList, path)
		return nil
	})

	if err != nil {
		log.Fatal("watcher.Watch(%q) failed: %s", testDir, err)
	}

	eventstream := watcher.Event
	var createReceived, modifyReceived, deleteReceived, renameReceived counter
	done := make(chan bool)
	go func() {
		for event := range eventstream {
			// Only count relevant events
			// if event.Name == filepath.Clean(testDir) || event.Name == filepath.Clean(testFile) {
			fmt.Println("event received: %s", event)
			var action int
			if event.IsDelete() {
				action = 0
				deleteReceived.increment()
			}
			if event.IsModify() {
				action = 1
				modifyReceived.increment()
			}
			if event.IsCreate() {
				action = 2
				createReceived.increment()
			}
			if event.IsRename() {
				action = 3
				renameReceived.increment()
			}
			if !excludeFileExt(event.Name) {
				local.Sync(event.Name, action)
			}
			// } else {
			//     log.Fatal("unexpected event received: %s", event)
			// }
		}
	}()

	done <- true

	log.Println("calling Close()")
	watcher.Close()
	log.Println("waiting for the event channel to become closed...")
	select {
	case <-done:
		log.Println("event channel closed")
	case <-time.After(2 * time.Second):
		log.Println("event stream was not closed after 2 seconds")
	}
}

func excludeFileExt(path string) bool {
	var extension = filepath.Ext(path)
	switch extension {
	case ".swp":
		return true
	case ".tmp":
		return true
	default:
		return false
	}
	return false
}
