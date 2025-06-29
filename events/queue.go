// Блокирующая очередь.
package main

import (
	"fmt"
	"sync"
)

// начало решения

// Queue - блокирующая FIFO-очередь.
type Queue struct {
	// TODO: переделать на срез и sync.Cond.
	cond *sync.Cond
	buf  []int
	mu   sync.Mutex
}

// NewQueue создает новую очередь.
func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.mu)
	return q
}

// Put добавляет элемент в очередь.
// Поскольку очередь безразмерная, никогда не блокируется.
func (q *Queue) Put(item int) {
	q.mu.Lock()
	q.buf = append(q.buf, item)
	q.cond.Signal()
	q.mu.Unlock()
}

// Get извлекает элемент из очереди.
// Если очередь пуста, блокируется до момента,
// пока в очереди не появится элемент.
func (q *Queue) Get() int {
	q.mu.Lock()
	for len(q.buf) == 0 {
		q.cond.Wait()
	}
	v := q.buf[0]
	q.buf = q.buf[1:]
	q.mu.Unlock()
	return v
}

// Len возвращает количество элементов в очереди.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.buf)
}

// конец решения

func main() {
	var wg sync.WaitGroup
	q := NewQueue()

	wg.Add(1)
	go func() {
		for i := range 100 {
			q.Put(i)
		}
		wg.Done()
	}()
	wg.Wait()

	total := 0

	wg.Add(1)
	go func() {
		for range 100 {
			total += q.Get()
		}
		wg.Done()
	}()
	wg.Wait()

	fmt.Println("Put x100, Get x100, Total:", total)
}
