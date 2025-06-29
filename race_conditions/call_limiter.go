// Ограничитель вызовов
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrBusy = errors.New("busy")
var ErrCanceled = errors.New("canceled")

// начало решения

// throttle следит, чтобы функция fn выполнялась не более limit раз в секунду.
// Возвращает функции handle (выполняет fn с учетом лимита) и cancel (останавливает ограничитель).
func throttle(limit int, fn func()) (handle func() error, cancel func()) {
	t := struct {
		mu      sync.Mutex
		start   time.Time
		count   int
		limit   int
		stopped bool
	}{
		start: time.Now(),
		limit: limit,
	}
	handle = func() error {
		t.mu.Lock()
		if t.stopped {
			t.mu.Unlock()
			return ErrCanceled
		}
		now := time.Now()
		if now.Sub(t.start) >= time.Second {
			t.start = now
			t.count = 0
		}
		if t.count >= t.limit {
			t.mu.Unlock()
			return ErrBusy
		}
		t.count++
		t.mu.Unlock()

		fn()
		return nil
	}

	cancel = func() {
		t.mu.Lock()
		if !t.stopped {
			t.stopped = true
		}
		t.mu.Unlock()
	}
	return handle, cancel
}

// конец решения

func main() {
	work := func() {
		fmt.Print(".")
	}

	handle, cancel := throttle(5, work)
	defer cancel()

	const n = 8
	var nOK, nErr int
	for i := 0; i < n; i++ {
		err := handle()
		if err == nil {
			nOK += 1
		} else {
			nErr += 1
		}
	}
	fmt.Println()
	fmt.Printf("%d calls: %d OK, %d busy\n", n, nOK, nErr)
}
