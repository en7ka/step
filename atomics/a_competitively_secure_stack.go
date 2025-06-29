// Конкурентно-безопасный стек.
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// начало решения

// Stack представляет конкурентно-безопасный стек без блокировок.
type Stack struct {
	top atomic.Pointer[Node]
}

// Node представляет элемент стека.
type Node struct {
	val  int
	next *Node
}

// Push добавляет значение на вершину стека.
func (s *Stack) Push(val int) {
	for {
		old := s.top.Load()
		node := &Node{val: val, next: old}
		if s.top.CompareAndSwap(old, node) {
			return
		}
	}
}

// Pop удаляет и возвращает вершину стека.
// Если стек пуст, возвращает false.
func (s *Stack) Pop() (int, bool) {
	for {
		old := s.top.Load()
		if old == nil {
			return 0, false
		}
		next := old.next
		if s.top.CompareAndSwap(old, next) {
			return old.val, true
		}
	}
}

// конец решения

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)

	stack := &Stack{}
	for i := range 1000 {
		go func() {
			time.Sleep(time.Millisecond)
			stack.Push(i)
			wg.Done()
		}()
	}

	wg.Wait()
	count := 0
	for _, ok := stack.Pop(); ok; _, ok = stack.Pop() {
		count++
	}
	fmt.Println(count)
}
