package main

import (
	"fmt"
	"time"
)

var (
	Processed, Unprocessable, Unprocessed int
)

func SetTitle(title string) {
	fmt.Printf("\033]0;%s\007", title)
}

func ConsoleWorker() {
	t := time.NewTicker(time.Millisecond)

	select {
	case <-t.C:
		SetTitle(fmt.Sprintf("Unprocessed: %d, Processed: %d, Unprocessable: %d", Unprocessed, Processed, Unprocessable))
	}
}
