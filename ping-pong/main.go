package main

import (
	"fmt"
	"sync"
)

func ping(ch1 chan string, ch2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		ch1 <- "ping"
		msg := <-ch2
		fmt.Println("Received:", msg)
	}
}

func pong(ch1 chan string, ch2 chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 5; i++ {
		msg := <-ch1
		fmt.Println("Received:", msg)
		ch2 <- "pong"
	}
}

func main() {
	var wg sync.WaitGroup
	ch1 := make(chan string)
	ch2 := make(chan string)

	wg.Add(2)
	go ping(ch1, ch2, &wg)
	go pong(ch1, ch2, &wg)

	wg.Wait()
	fmt.Println("Ping-pong complete!")
}
