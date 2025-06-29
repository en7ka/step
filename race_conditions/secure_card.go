// Конкурентно-безопасная карта.
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ConcMap - безопасная в многозадачной среде карта.
type ConcMap[K comparable, V any] struct {
	items map[K]V
	lock  sync.Mutex
	wg    sync.WaitGroup
}

// NewConcMap создает новую карту.
func NewConcMap[K comparable, V any]() *ConcMap[K, V] {
	return &ConcMap[K, V]{items: map[K]V{}}
}

// Get возвращает значение по ключу.
func (cm *ConcMap[K, V]) Get(key K) V {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	return cm.items[key]
}

// Set устанавливает значение по ключу.
func (cm *ConcMap[K, V]) Set(key K, val V) {
	cm.lock.Lock()
	defer cm.lock.Unlock()
	cm.items[key] = val
}

// начало решения

// SetIfAbsent устанавливает новое значение по ключу
// и возвращает его, но только если такого ключа нет в карте.
// Если ключ уже есть - возвращает старое значение по ключу.
func (cm *ConcMap[K, V]) SetIfAbsent(key K, val V) V {
	// TODO: реализовать требования
	var v V
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if old, ok := cm.items[key]; ok {
		return old
	} else {
		cm.items[key] = val
	}
	v = cm.items[key]
	return v
}

// Compute устанавливает значение по ключу, применяя к нему функцию.
// Возвращает новое значение. Функция выполняется атомарно.
func (cm *ConcMap[K, V]) Compute(key K, f func(V) V) V {
	// TODO: реализовать требования
	var v V
	cm.lock.Lock()
	defer cm.lock.Unlock()
	if q, ok := cm.items[key]; ok {
		v = q
	}
	newval := f(v)
	cm.items[key] = newval
	return newval
}

// конец решения

func getSet() {
	var wg sync.WaitGroup
	wg.Add(2)

	m := NewConcMap[string, int]()

	go func() {
		defer wg.Done()
		m.Set("hello", rand.Intn(100))
	}()

	go func() {
		defer wg.Done()
		m.Set("hello", rand.Intn(100))
	}()

	wg.Wait()
	fmt.Println("hello =", m.Get("hello"))
	// hello = 71 (случайное)
}

func setIfAbsent() {
	var wg sync.WaitGroup
	wg.Add(2)

	m := NewConcMap[string, int]()

	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Millisecond)
		m.SetIfAbsent("hello", 42)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		m.SetIfAbsent("hello", 84)
	}()

	wg.Wait()
	fmt.Println("hello =", m.Get("hello"))
	// hello = 42 (от первой горутины)
}

func compute() {
	var wg sync.WaitGroup
	wg.Add(2)

	m := NewConcMap[string, int]()

	go func() {
		defer wg.Done()
		for range 100 {
			m.Compute("hello", func(v int) int {
				return v + 1
			})
		}
	}()

	go func() {
		defer wg.Done()
		for range 100 {
			m.Compute("hello", func(v int) int {
				return v + 1
			})
		}
	}()

	wg.Wait()
	fmt.Println("hello =", m.Get("hello"))
	// hello = 200 (каждая горутина увеличила hello на 100)
}

func main() {
	getSet()
	fmt.Println("---")
	setIfAbsent()
	fmt.Println("---")
	compute()
}
