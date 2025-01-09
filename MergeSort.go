package main

import (
	"fmt"
	"sync"
)

func MergeSortedArrays(arr []int, low int, high int, mid int) {
	temp := make([]int, 0)

	l := low
	r := mid + 1

	for l <= mid && r <= high {
		if arr[l] < arr[r] {
			temp = append(temp, arr[l])
			l++
		} else {
			temp = append(temp, arr[r])
			r++
		}
	}

	for l <= mid {
		temp = append(temp, arr[l])
		l++
	}

	for r <= high {
		temp = append(temp, arr[r])
		r++
	}

	for i := 0; i < len(temp); i++ {
		arr[low+i] = temp[i]
	}
}

func Merge(arr []int, low int, high int, wg *sync.WaitGroup) {
	defer wg.Done()

	if low < high {
		mid := (low + high) / 2

		subWg := &sync.WaitGroup{}
		subWg.Add(2)

		go Merge(arr, low, mid, subWg)
		go Merge(arr, mid+1, high, subWg)

		subWg.Wait()

		MergeSortedArrays(arr, low, high, mid)
	}
}

func main() {
	arr := []int{38, 27, 43, 3, 9, 82, 10}
	fmt.Println("Before sorting:", arr)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go Merge(arr, 0, len(arr)-1, wg)

	wg.Wait()

	fmt.Println("After sorting:", arr)
}
