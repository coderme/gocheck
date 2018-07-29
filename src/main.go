package main

import (
	//	"fmt"
	"time"
)

var (
	urlQueue    = make(chan string, 100)
	resultQueue = make(chan *Result, 100)
	limit       = make(chan struct{}, *concurrency)
	touched     time.Time
)

func main() {
	urlQueue <- website
	errs := 0
	tick := time.Tick(defaultWait)

	for {
		select {
		case u := <-urlQueue:
			touched = time.Now()
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
			if r.ErrorServer || r.ErrorClient {
				errs++
			}

			if errs >= *maxErrsCount {
				exit(1)
			}
		case <-tick:
			if time.Since(touched) > *timeWait {
				showDebug("Timewait:", "Done")
				exit(0)
			}
		}
	}
}
