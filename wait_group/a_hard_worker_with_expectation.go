// Работяга
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// начало решения

// Worker выполняет заданную функцию в цикле, пока не будет остановлен.
type Worker struct {
	fn func() error
	// TODO: добавить поля
	started bool
	stopped bool
	wg      sync.WaitGroup
}

// NewWorker создает новый экземпляр Worker с заданной функцией.
func NewWorker(fn func() error) *Worker {
	return &Worker{fn: fn}
}

// Start запускает отдельную горутину, в которой циклически
// выполняет заданную функцию, пока не будет вызван метод Stop,
// либо пока функция не вернет ошибку.
// Повторные вызовы Start игнорируются.
// Гарантируется, что Start не вызывается из разных горутин.
func (w *Worker) Start() {
	// TODO: реализовать требования
	if w.started == true {
		return
	}
	w.started = true
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for !w.stopped {
			err := w.fn()
			if err != nil {
				w.stopped = true
				return
			}
		}
	}()
}

// Stop останавливает выполнение цикла.
// Вызов Stop до Start игнорируется.
// Повторные вызовы Stop игнорируются.
// Гарантируется, что Stop не вызывается из разных горутин.
func (w *Worker) Stop() {
	// TODO: реализовать требования
	if w.stopped == true {
		return
	}
	w.stopped = true
}

// Wait блокирует вызвавшую его горутину до тех пор,
// пока Worker не будет остановлен (из-за ошибки или вызова Stop).
// Wait может вызываться несколько раз, в том числе из разных горутин.
// Wait может вызываться до Start. Это не приводит к блокировке.
// Wait может вызываться после Stop. Это не приводит к блокировке.
func (w *Worker) Wait() {
	// TODO: реализовать требования
	if !w.started {
		return
	}
	w.wg.Wait()
}

// конец решения

func main() {
	{
		// Завершение по ошибке
		count := 3
		fn := func() error {
			fmt.Print(count, " ")
			count--
			if count == 0 {
				return errors.New("count is zero")
			}
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		worker := NewWorker(fn)
		worker.Start()
		time.Sleep(25 * time.Millisecond)

		fmt.Println()
		// 3 2 1
	}
	{
		// Завершение по Stop
		count := 3
		fn := func() error {
			fmt.Print(count, " ")
			count--
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		worker := NewWorker(fn)
		worker.Start()
		time.Sleep(25 * time.Millisecond)
		worker.Stop()

		fmt.Println()
		// 3 2 1
	}
	{
		// Ожидание завершения через Wait
		count := 3
		fn := func() error {
			fmt.Print(count, " ")
			count--
			time.Sleep(10 * time.Millisecond)
			return nil
		}

		worker := NewWorker(fn)
		worker.Start()

		// эта горутина остановит работягу через 25 мс
		go func() {
			time.Sleep(25 * time.Millisecond)
			worker.Stop()
		}()

		// подождем, пока кто-нибудь остановит работягу
		worker.Wait()
		fmt.Println("done")

		// 3 2 1 done
	}
}
