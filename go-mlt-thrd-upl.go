package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"time"
)

func watcherHandler(w *watcher.Watcher) {
	for {
		select {
		case event := <-w.Event:
			if event.IsDir() && event.Op == watcher.Write {
				continue
			}
			fmt.Println(event)
		case err := <-w.Error:
			log.Fatalln(err)
		case <-w.Closed:
			return
		}
	}
}

func main() {
	w := watcher.New()

	go watcherHandler(w)

	if err := w.AddRecursive("./monitoring_folder"); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
