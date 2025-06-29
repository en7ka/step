// Статистика вызовов.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// начало решения

// External представляет внешний сервис.
type External struct {
	lastCall atomic.Value
	numCalls atomic.Int32
}

// NewExternal создает новый экземпляр External.
func NewExternal() *External {
	return &External{}
}

// Call вызывает внешний сервис.
func (e *External) Call() {
	// вызываем внешний сервис...
	e.lastCall.Store(time.Now())
	e.numCalls.Add(1)
}

// LastCall возвращает время последнего вызова.
func (e *External) LastCall() time.Time {
	return e.lastCall.Load().(time.Time)
}

// NumCalls возвращает количество вызовов.
func (e *External) NumCalls() int {
	return int(e.numCalls.Load())
}

// конец решения

func main() {
	const nConc = 4
	var wg sync.WaitGroup
	wg.Add(nConc)

	// Вызываем внешнюю систему из нескольких горутин.
	ext := NewExternal()
	for range nConc {
		go func() {
			defer wg.Done()
			for range 10 {
				ext.Call()
			}
		}()
	}

	wg.Wait()

	// Количество вызовов и время последнего вызова.
	fmt.Println("Calls:", ext.NumCalls())
	fmt.Println("Last call:", ext.LastCall().Format("15:04:05"))
	// Calls: 40
	// Last call: 15:04:05
}
