package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
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

// TODO: Прошу помочь с написанием теста.
// В соответствии с комментарием от Олега Венгера:
// "привет, всё верно. необходимо ограничить сверху время ожидания завершения тасок.
// Но тут скорее нужно подумать, как можно запустить n тасок чтобы они все были заблокированы и чего-то ждали.
// А потом сразу всем сказать закончили.
// Тогда в целом при запуске eventually можно прикинуть, что для завершения всех тасок потребуется сколько-то времени"

// Я, с помощью sync.Cond, смог запустить все таски одновременно и
// заставить их ждать сигнала. Но, как посчитать время их выполнения для
// аргумента waitFor у require.Eventually, я не могу понять.
// Подскажите, пожалуйста, в верном направлении я двигаюсь? Мне жутко интересно, как это можно сделать

// Тест-кейс представлен ниже

// t.Run("tasks without errors with eventually", func(t *testing.T) {
//	tasksCount := 50
//	tasks := make([]Task, 0, tasksCount)
//
//	var runTasksCount int32
//
//	wg := sync.WaitGroup{}
//	once := sync.Once{}
//	cond := sync.NewCond(&sync.Mutex{})
//	workersCount := tasksCount
//	maxErrorsCount := 1
//
//	var start time.Time
//	var elapsed time.Duration
//
//	for i := 0; i < tasksCount; i++ {
//		tasks = append(tasks, func() error {
//			once.Do(func() {
//				start = time.Now()
//			})
//
//			atomic.AddInt32(&runTasksCount, 1)
//			cond.L.Lock()
//			cond.Wait()
//			cond.L.Unlock()
//			return nil
//		})
//	}
//
//	go func(wg *sync.WaitGroup) {
//		wg.Add(1)
//		defer wg.Done()
//		tasksCountInt32 := int32(tasksCount)
//		for {
//			if atomic.LoadInt32(&runTasksCount) == tasksCountInt32 {
//				cond.Broadcast()
//				elapsed = time.Since(start)
//				break
//			}
//		}
//	}(&wg)
//
//	condition := func() bool {
//		err := Run(tasks, workersCount, maxErrorsCount)
//		return err == nil
//	}
//
//	//err := Run(tasks, workersCount, maxErrorsCount)
//	//require.NoError(t, err)
//
//	require.Eventually(t, condition, time.Second*10, 1*time.Millisecond, "tasks were run sequentially?")
//	//require.Eventually(t, func() bool {
//	//	return true
//	//}, elapsed, 1*time.Millisecond, "tasks were run sequentially?")
//	require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
//	fmt.Println(elapsed)
//	wg.Wait()
// })
