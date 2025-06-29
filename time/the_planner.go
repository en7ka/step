package main

import (
	"fmt"
	"time"
)

// начало решения

func schedule(dur time.Duration, fn func()) func() {
	canceled := make(chan struct{})
	ticker := time.NewTicker(dur)
	sem := make(chan struct{}, 1)
	sem <- struct{}{}
	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case <-sem:
					go func() {
						defer func() {
							sem <- struct{}{}
						}()
						fn()
					}()
				default:

				}
			case <-canceled:
				return
			}
		}
	}()

	cancel := func() {
		defer func() { recover() }()
		ticker.Stop()
		close(canceled)
	}
	return cancel
}

// конец решения

func main() {
	work := func() {
		at := time.Now()
		fmt.Printf("%s: work done\n", at.Format("15:04:05.000"))
	}

	cancel := schedule(50*time.Millisecond, work)
	defer cancel()

	// хватит на 5 тиков
	time.Sleep(260 * time.Millisecond)
}
