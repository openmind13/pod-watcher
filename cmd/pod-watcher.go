package main

import (
	"log"
	"pod-watcher/watcher"
)

func main() {
	watcher, err := watcher.New()
	if err != nil {
		log.Fatal(err)
	}
	watcher.Start()
}
