package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in maxGoroutine goroutines and stops its work when receiving maxError errors from tasks.
func Run(tasks []Task, maxGoroutine, maxError int) error {
	if maxGoroutine <= 0 {
		maxGoroutine = len(tasks)
	}
	if maxError <= 0 {
		maxError = len(tasks) + 1
	}

	var (
		errCounter int32
		wg         sync.WaitGroup
	)
	sem := make(chan struct{}, maxGoroutine)

	for i := 0; i < len(tasks); i++ {
		sem <- struct{}{}
		if int(atomic.LoadInt32(&errCounter)) >= maxError {
			break
		}

		wg.Add(1)
		go func(task Task) {
			defer func() {
				wg.Done()
				<-sem
			}()

			if err := task(); err != nil {
				atomic.AddInt32(&errCounter, 1)
			}
		}(tasks[i])
	}
	wg.Wait()

	if int(errCounter) >= maxError {
		return ErrErrorsLimitExceeded
	}

	return nil
}
