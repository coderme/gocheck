package main

import (
//	"fmt"
)

var (
	urlQueue    = make(chan string, 100)
	resultQueue = make(chan *Result, 100)
	limit       = make(chan struct{}, *concurrency)
)

func main() {
	urlQueue <- website

	for {
		select {
		case u := <-urlQueue:
			limit <- struct{}{}

			go func() {
				resultQueue <- fetch(u)
			}()
		case r := <-resultQueue:
			if r == nil {
				continue
			}
			if *outJSON {
				r.PrintJSON()
			} else {
				r.PrintText()
			}
		}
	}
}
