package main

import (
	"errors"
	"fmt"
	"time"
)

var ErrCanceled = errors.New("canceled")

func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	interval := time.Second / time.Duration(limit)
	tokens := make(chan struct{}, 1)
	stop := make(chan struct{})

	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				select {
				case tokens <- struct{}{}:
				default:
				}
			case <-stop:
				return
			}
		}
	}()

	handle = func() error {
		select {
		case <-stop:
			return ErrCanceled
		default:
		}
		select {
		case <-stop:
			return ErrCanceled
		case <-tokens:
			go fn()
			return nil
		}
	}

	var done bool
	cancel = func() {
		if done {
			return
		}
		done = true
		close(stop)
	}

	return
}

func main() {
	work := func() { fmt.Print(".") }

	handle, cancel := throttle(5, work)
	defer cancel()

	start := time.Now()
	for i := 0; i < 10; i++ {
		handle()
	}
	fmt.Println()
	fmt.Printf("10 queries took %v\n", time.Since(start))
}
