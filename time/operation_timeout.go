package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// выполняет какую-то операцию,
// обычно быстро, но иногда медленно
func work() int {
	if rand.Intn(10) < 8 {
		time.Sleep(10 * time.Millisecond)
	} else {
		time.Sleep(200 * time.Millisecond)
	}
	return 42
}

// выполняет функцию fn() c таймаутом timeout и возвращает результат
// если в течение timeout функция не вернула ответ - возвращает ошибку
func withTimeout(fn func() int, timeout time.Duration) (int, error) {
	var result int

	done := make(chan struct{})
	go func() {
		result = fn()
		close(done)
	}()

	select {
	case <-done:
		return result, nil
	case <-after(timeout):
		return 0, errors.New("timeout")
	}
}

// начало решения

// возвращает канал, в котором появится значение
// через промежуток времени dur
func after(dur time.Duration) <-chan time.Time {
	result := make(chan time.Time, 1)
	go func() {
		time.Sleep(dur)
		result <- time.Now()
	}()
	return result
}

// конец решения

func main() {
	for i := 0; i < 10; i++ {
		start := time.Now()
		timeout := 50 * time.Millisecond
		if answer, err := withTimeout(work, timeout); err != nil {
			fmt.Printf("Took %v. Error: %v\n", time.Since(start), err)
		} else {
			fmt.Printf("Took %v. Result: %v\n", time.Since(start), answer)
		}
	}
}
