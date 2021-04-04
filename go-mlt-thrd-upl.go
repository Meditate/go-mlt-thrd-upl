package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"time"
)

func main() {
	w := watcher.New()

	go watcherHandler(w)

	if err := w.AddRecursive("./monitoring_folder"); err != nil {
		log.Fatalln(err)
	}

	fmt.Println()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func watcherHandler(w *watcher.Watcher) {
	for {
		select {
		case event := <-w.Event:
			if event.IsDir() && event.Op == watcher.Write {
				continue
			} else {
				fmt.Println(event.Path)

				processFile(event.Path)
			}
		case err := <-w.Error:
			log.Fatalln(err)
		case <-w.Closed:
			return
		}
	}
}
