package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func say(done chan<- struct{}, id int, phrase string) {
	for _, word := range strings.Fields(phrase) {
		fmt.Printf("Worker #%d says: %s...\n", id, word)
		dur := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(dur)
	}
	done <- struct{}{}

}
func main() {
	phrases := []string{
		"go is awesome",
		"cats are cute",
		"rain is wet",
		"channels are hard",
		"floor is lava",
	}

	done := make(chan struct{})

	for idx, phrase := range phrases {
		go say(done, idx, phrase)
	}

	for i := 0; i < len(phrases); i++ {
		<-done
	}
}
