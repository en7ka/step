package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var mu sync.Mutex
var delta atomic.Int32
var counter atomic.Int32

func increment() {
	delta.Add(1)
	sleep(10)
	counter.Add(delta.Load())
}

func sleep(maxMs int) {
	dur := time.Duration(rand.Intn(maxMs) * int(time.Millisecond))
	time.Sleep(dur)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(100)

	for range 100 {
		go func() {
			increment()
			wg.Done()
		}()
	}

	wg.Wait()
	d, c := delta.Load(), counter.Load()
	fmt.Println(d, c)
}
