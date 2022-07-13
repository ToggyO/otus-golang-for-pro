package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, maxWorkersCount, maxErrorsCount int) error {
	if maxErrorsCount < 0 {
		maxErrorsCount = 0
	}

	taskCount := len(tasks)
	errorsCounter := int32(maxErrorsCount)
	workersCount := calcWorkersCount(taskCount, maxWorkersCount)
	workerQueue := make(chan Task, taskCount)
	wg := sync.WaitGroup{}

	for i := 1; i <= workersCount; i++ {
		wg.Add(1)
		go worker(workerQueue, &errorsCounter, &wg)
	}

	for _, t := range tasks {
		workerQueue <- t
	}
	close(workerQueue)

	wg.Wait()

	if atomic.LoadInt32(&errorsCounter) < 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(workerQueue <-chan Task, maxErrs *int32, wg *sync.WaitGroup) {
	defer wg.Done()

	for atomic.LoadInt32(maxErrs) >= 0 {
		t, ok := <-workerQueue
		if !ok {
			break
		}
		if t() != nil {
			atomic.AddInt32(maxErrs, -1)
		}
	}
}

func calcWorkersCount(tasksCount, maxWorkersCount int) int {
	workersCount := maxWorkersCount
	if tasksCount < maxWorkersCount {
		workersCount = tasksCount
	}

	return workersCount
}
