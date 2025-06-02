package main

import (
	"fmt"
	"sync"
)

func worker(nums []int, ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	sum := 0
	for _, num := range nums {
		sum += num
	}
	ch <- sum
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	n := 4

	chunkSize := len(nums) / n
	ch := make(chan int, n)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		start := i * chunkSize
		end := start + chunkSize

		wg.Add(1)
		go worker(nums[start:end], ch, &wg)
	}

	wg.Wait()
	close(ch)

	totalSum := 0
	for partial := range ch {
		totalSum += partial
	}

	fmt.Println("Total Sum:", totalSum)
}
