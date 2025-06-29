// Безопасная карта
package main

import (
	"fmt"
	"sync"
)

// начало решения

// Counter представляет безопасную карту частот слов.
// Ключ - строка, значение - целое число.
type Counter struct {
	lock sync.Mutex
	m    map[string]int
}

// Increment увеличивает значение по ключу на 1.
func (c *Counter) Increment(str string) {
	c.lock.Lock()
	c.m[str]++
	c.lock.Unlock()
}

// Value возвращает значение по ключу.
func (c *Counter) Value(str string) int {
	c.lock.Lock()
	result := c.m[str]
	c.lock.Unlock()
	return result
}

// Range проходит по всем записям карты,
// и для каждой вызывает функцию fn, передавая в нее ключ и значение.
func (c *Counter) Range(fn func(key string, val int)) {
	for word := range c.m {
		c.lock.Lock()
		fn(word, c.m[word])
		c.lock.Unlock()
	}
}

// NewCounter создает новую карту частот.
func NewCounter() *Counter {
	c := &Counter{m: make(map[string]int)}
	return c
}

// конец решения

func main() {
	counter := NewCounter()

	var wg sync.WaitGroup
	wg.Add(3)

	increment := func(key string, val int) {
		defer wg.Done()
		for ; val > 0; val-- {
			counter.Increment(key)
		}
	}

	go increment("one", 100)
	go increment("two", 200)
	go increment("three", 300)

	wg.Wait()

	fmt.Println("two:", counter.Value("two"))

	fmt.Print("{ ")
	counter.Range(func(key string, val int) {
		fmt.Printf("%s:%d ", key, val)
	})
	fmt.Println("}")
}
