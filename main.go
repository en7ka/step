package main

import (
	"fmt"
	"strings"
	"unicode"
)

type nextFunc func() string
type counter map[string]int

type pair struct {
	word  string
	count int
}

func countDigitsInWords(next nextFunc) counter {
	pending := make(chan string)
	counted := make(chan pair)

	go func() {
		for {
			w := next()
			pending <- w
			if w == "" {
				break
			}
		}
	}()

	go func() {
		defer close(counted)
		for w := range pending {
			counted <- pair{w, countDigits(w)}
			if w == "" {
				break
			}
		}
	}()

	stats := make(counter)
	for p := range counted {
		if p.word == "" {
			break
		}
		stats[p.word] = p.count
	}
	return stats
}

func countDigits(str string) int {
	count := 0
	for _, r := range str {
		if unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

func printStats(stats counter) {
	for word, count := range stats {
		fmt.Printf("%s: %d\n", word, count)
	}
}

func wordGenerator(phrase string) nextFunc {
	words := strings.Fields(phrase)
	i := 0
	return func() string {
		if i >= len(words) {
			return ""
		}
		w := words[i]
		i++
		return w
	}
}

func main() {
	phrase := "0ne 1wo thr33 4068"
	next := wordGenerator(phrase)
	stats := countDigitsInWords(next)
	printStats(stats)
}
