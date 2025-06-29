package main

import (
	"errors"
	"fmt"
)

var ErrFull = errors.New("Queue is full")
var ErrEmpty = errors.New("Queue is empty")

// начало решения

// Queue - FIFO-очередь на n элементов
type Queue struct {
	ch chan int
}

// Get возвращает очередной элемент.
// Если элементов нет и block = false -
// возвращает ошибку.
func (q Queue) Get(block bool) (int, error) {
	if block {
		return <-q.ch, nil
	}
	select {
	case v := <-q.ch:
		return v, nil
	default:
		return 0, ErrEmpty
	}
}

// Put помещает элемент в очередь.
// Если очередь заполнения и block = false -
// возвращает ошибку.
func (q Queue) Put(val int, block bool) error {
	if block {
		q.ch <- val
		return nil
	}
	select {
	case q.ch <- val:
		return nil
	default:
		return ErrFull
	}
}

// MakeQueue создает новую очередь
func MakeQueue(n int) Queue {
	ch := make(chan int, n)
	return Queue{ch: ch}
}

// конец решения
func main() {
	q := MakeQueue(2)

	err := q.Put(1, false)
	fmt.Println("put 1:", err)

	err = q.Put(2, false)
	fmt.Println("put 2:", err)

	err = q.Put(3, false)
	fmt.Println("put 3:", err)

	res, err := q.Get(false)
	fmt.Println("get:", res, err)

	res, err = q.Get(false)
	fmt.Println("get:", res, err)

	res, err = q.Get(false)
	fmt.Println("get:", res, err)
}
