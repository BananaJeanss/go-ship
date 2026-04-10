//go:build ignore
package main

import (
	"math/rand"
	"time"
	"fmt"
)

var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
  rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}

func typewriterMonkey(i int, c chan string) {
	// typewriter monkey
	for i := 0; i < 10; i++ {
		yeah := StringWithCharset(10, charset)
		formatted := fmt.Sprintf("Monkey %d: %s", i+1, yeah)
		c <- formatted
	}
}

func main() {
	c := make(chan string)

	for i, j := 0, 10; i < j; i++ { // runs 10 times
		go typewriterMonkey(i, c)
	}
	
	// finally, read the output
	for i := 0; i < 100; i++ {
		// this reads, and then clears from the channel
		fmt.Println(<-c)
	}
}
