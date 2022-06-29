package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("max errors count equals 0 or less", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := -1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("no tasks", func(t *testing.T) {
		tasksCount := 0
		tasks := make([]Task, 0, tasksCount)

		workersCount := 5
		maxErrorsCount := 1

		err := Run(tasks, workersCount, maxErrorsCount)
		require.NoError(t, err)
	})
}

// TODO: Скажу честно - выполнил задание со звездочкой после занятия по разбору ДЗ.
// Применил sync.Cond, чтобы заставить таски ожидать команды к завершению
func TestConcurrencyEventually(t *testing.T) {
	tasksCount := 5
	workersCount := 5
	maxErrorsCount := 1

	tasks := make([]Task, 0, tasksCount)
	cond := sync.NewCond(&sync.Mutex{})
	wg := sync.WaitGroup{}

	var runTasksCount int32
	for i := 0; i < tasksCount; i++ {
		tasks = append(tasks, func() error {
			atomic.AddInt32(&runTasksCount, 1)
			cond.L.Lock()
			cond.Wait()
			cond.L.Unlock()
			return nil
		})
	}

	go func(wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		tasksCountInt32 := int32(tasksCount)
		for {
			if atomic.LoadInt32(&runTasksCount) == tasksCountInt32 {
				cond.Broadcast()
				break
			}
		}
	}(&wg)

	errChan := make(chan error)
	go func() {
		errChan <- Run(tasks, workersCount, maxErrorsCount)
	}()

	require.Eventually(t, func() bool {
		return atomic.LoadInt32(&runTasksCount) == int32(workersCount)
	}, time.Second, time.Millisecond)
	wg.Wait()

	require.NoError(t, <-errChan)
}
