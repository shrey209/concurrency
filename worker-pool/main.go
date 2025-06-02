package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id int
}

func (t *Task) process() {
	fmt.Printf("Processing task ID: %d\n", t.Id)
	time.Sleep(time.Second)
}

type workerPool struct {
	Tasks       []Task
	Concurrency int
	taskChan    chan Task
	wg          sync.WaitGroup
}

func (wp *workerPool) worker() {
	for task := range wp.taskChan {
		task.process()
		wp.wg.Done()
	}
}

func (wp *workerPool) run() {
	wp.taskChan = make(chan Task, len(wp.Tasks))

	for i := 0; i < wp.Concurrency; i++ {
		go wp.worker()
	}

	wp.wg.Add(len(wp.Tasks))

	for _, task := range wp.Tasks {
		wp.taskChan <- task
	}

	close(wp.taskChan)
	wp.wg.Wait()
}

func main() {
	// Create 20 tasks
	tasks := make([]Task, 20)
	for i := 0; i < 20; i++ {
		tasks[i] = Task{Id: i + 1}
	}

	wp := workerPool{
		Tasks:       tasks,
		Concurrency: 5,
	}

	wp.run()
	fmt.Println("Finished all tasks.")
}
