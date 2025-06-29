// TryAcquire
package main

import (
	"fmt"
	"sync"
	"time"
)

// начало решения

// Semaphore представляет семафор синхронизации.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore создает новый семафор указанной вместимости.
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, n)}
}

// Acquire занимает место в семафоре, если есть свободное.
// В противном случае блокирует вызывающую горутину.
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// TryAcquire занимает место в семафоре, если есть свободное,
// и возвращает true. В противном случае просто возвращает false.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Release освобождает место в семафоре и разблокирует
// одну из заблокированных горутин (если такие были).
func (s *Semaphore) Release() {
	<-s.ch
}

// конец решения

func main() {
	const maxConc = 4
	sema := NewSemaphore(maxConc)

	const nCalls = 12
	var wg sync.WaitGroup
	wg.Add(nCalls)

	var nOK, nBusy int
	for i := 0; i < nCalls; i++ {
		if !sema.TryAcquire() {
			nBusy++
			wg.Done()
			continue
		}
		go func() {
			defer wg.Done()
			defer sema.Release()
			time.Sleep(10 * time.Millisecond)
			fmt.Print(".")
			nOK++
		}()
	}

	wg.Wait()

	fmt.Println()
	fmt.Printf("%d calls: %d OK, %d busy\n", nCalls, nOK, nBusy)
	/*
		....
		12 calls: 4 OK, 8 busy
	*/
}
