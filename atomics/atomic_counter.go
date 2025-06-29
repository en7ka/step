// Атомарный счетчик
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// начало решения

// Total представляет атомарный счетчик.
type Total struct {
	total int32
	count atomic.Int32
}

// Increment увеличивает счетчик на 1.
func (t *Total) Increment() {
	t.count.Add(1)
}

// Value возвращает значение счетчика.
func (t *Total) Value() int {
	total := t.count.Load()
	return int(total)
}

// конец решения

func main() {
	var wg sync.WaitGroup

	var total Total
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10000 {
				total.Increment()
			}
		}()
	}

	wg.Wait()
	fmt.Println("total", total.Value())
}
