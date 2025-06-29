// Атомарные операции (sync/atomic)
// Увеличение общего int в нескольких горутинах
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var wg sync.WaitGroup

	var total atomic.Int32
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10000 {
				total.Add(1)
			}
		}()
	}

	wg.Wait()
	fmt.Println("total", total.Load())
}
