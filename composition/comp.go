package main

import (
	"fmt"
)

// начало решения

// count отправляет в канал числа от start до бесконечности
func count(cancel <-chan struct{}, start int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := start; ; i++ {
			select {
			case out <- i:
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// take выбирает первые n чисел из in и отправляет в выходной канал
func take(cancel <-chan struct{}, in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case out <- <-in:
			case <-cancel:
				return
			}
		}
	}()
	return out
}

// конец решения

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	stream := take(cancel, count(cancel, 10), 5)
	first := <-stream
	second := <-stream
	third := <-stream

	fmt.Println(first, second, third)
}
