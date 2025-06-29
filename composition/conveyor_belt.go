package main

import (
	//"compress/flate"
	"fmt"
	//"go/constant"
	"math/rand"
)

// начало решения

// генерит случайные слова из 5 букв
// с помощью randomWord(5)
func generate(cancel <-chan struct{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				return
			case out <- randomWord(5):
			}
		}
	}()
	return out
}

// выбирает слова, в которых не повторяются буквы,
// abcde - подходит
// abcda - не подходит
func takeUnique(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				seen := make(map[rune]struct{})
				unique := true
				for _, r := range val {
					if _, exists := seen[r]; exists {
						unique = false
						break
					}
					seen[r] = struct{}{}
				}
				if !unique {
					continue
				}
				select {
				case out <- val:
				case <-cancel:
					return
				}
			}
		}
	}()
	return out
}

// переворачивает слова
// abcde -> edcba
func reverse(cancel <-chan struct{}, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-cancel:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				rune := []rune(val)
				for i, j := 0, len(rune)-1; i < j; i, j = i+1, j-1 {
					rune[i], rune[j] = rune[j], rune[i]
				}
				select {
				case out <- string(rune):
				case <-cancel:
					return
				}
			}
		}
	}()
	return out
}

// объединяет c1 и c2 в общий канал
func merge(cancel <-chan struct{}, c1, c2 <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for c1 != nil || c2 != nil {
			select {
			case <-cancel:
				return
			case val1, ok := <-c1:
				if !ok {
					c1 = nil
				} else {
					select {
					case out <- val1:
					case <-cancel:
						return
					}
				}
			case val2, ok := <-c2:
				if !ok {
					c2 = nil
				} else {
					select {
					case out <- val2:
					case <-cancel:
						return
					}

				}
			}
		}
	}()
	return out
}

// печатает первые n результатов
func print(cancel <-chan struct{}, in <-chan string, n int) {
	for i := 0; i < n; i++ {
		select {
		case <-cancel:
			return
		case val, ok := <-in:
			if !ok {
				return
			}
			runes := []rune(val)
			for a, b := 0, len(runes)-1; a < b; a, b = a+1, b-1 {
				runes[a], runes[b] = runes[b], runes[a]
			}
			fmt.Printf("%s -> %s\n", val, string(runes))
		}
	}
}

// конец решения

// генерит случайное слово из n букв
func randomWord(n int) string {
	const letters = "aeiourtnsl"
	chars := make([]byte, n)
	for i := range chars {
		chars[i] = letters[rand.Intn(len(letters))]
	}
	return string(chars)
}

func main() {
	cancel := make(chan struct{})
	defer close(cancel)

	c1 := generate(cancel)
	c2 := takeUnique(cancel, c1)
	c3_1 := reverse(cancel, c2)
	c3_2 := reverse(cancel, c2)
	c4 := merge(cancel, c3_1, c3_2)
	print(cancel, c4, 10)
}
