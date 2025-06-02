package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Node[T any] struct {
	value   T
	next    *Node[T]
	version int32
}

type Queue[T any] struct {
	head  *Node[T]
	tail  *Node[T]
	mu    sync.Mutex
	retry int
	size  int32
}

func (q *Queue[T]) enqueue(element T) bool {

	for i := 0; i < q.retry; i++ {
		if q.enqueueOptimistic(element) {
			atomic.AddInt32(&q.size, 1)
			return true
		}
	}

	if q.enqueuePessimistic(element) {
		atomic.AddInt32(&q.size, 1)
		return true
	}
	return false
}

func (q *Queue[T]) enqueueOptimistic(element T) bool {

	if q.tail == nil {
		newNode := &Node[T]{value: element, next: nil, version: 1}
		q.tail = newNode
		q.head = newNode
		return true
	}

	current := atomic.LoadInt32(&q.tail.version)

	if atomic.CompareAndSwapInt32(&q.tail.version, current, current+1) {
		newNode := &Node[T]{value: element, next: nil, version: current + 1}
		q.tail.next = newNode
		q.tail = newNode
		return true
	}

	return false
}

func (q *Queue[T]) enqueuePessimistic(element T) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	newNode := &Node[T]{value: element, next: nil, version: 1}

	if q.tail == nil {
		q.tail = newNode
		q.head = newNode
	} else {
		q.tail.next = newNode
		q.tail = newNode
	}

	return true
}

func (q *Queue[T]) dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.head == nil {
		var zeroValue T
		return zeroValue, false
	}

	value := q.head.value
	q.head = q.head.next

	if q.head == nil {
		q.tail = nil
	}

	atomic.AddInt32(&q.size, -1)
	return value, true
}

func (q *Queue[T]) Size() int32 {
	return atomic.LoadInt32(&q.size)
}

func main() {

	q := &Queue[int]{retry: 3}

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			q.enqueue(i)
		}(i)
	}

	wg.Wait()
	fmt.Printf("Queue size after enqueuing 1000 elements: %d\n", q.Size())
}
