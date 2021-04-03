package main

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"io"
	"log"
	"math"
	"os"
	"time"
)

func processFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}

	fileStat, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	defer file.Close()

	streamSize := math.Round(float64((fileStat.Size() / 3.0)) + 0.5)

	fmt.Printf("file size: %d\n", fileStat.Size())

	p := make([]byte, int(streamSize))

	for i := 0; i < 3; i++ {
		n, err := file.Read(p)

		if err == io.EOF {
			log.Fatalln(err)
		}

		fmt.Printf("%d bytes read\n", n)
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
