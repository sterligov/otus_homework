package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in nWorkers goroutines and stops its work when receiving nErrors errors from tasks.
func Run(tasks []Task, nWorkers, nErrors int) error {
	if nWorkers <= 0 {
		nWorkers = len(tasks)
	}
	if nErrors <= 0 {
		nErrors = len(tasks) + 1
	}

	errors := make(chan error)
	taskCh := make(chan Task)
	done := make(chan struct{})

	wg := &sync.WaitGroup{}
	wg.Add(nWorkers)

	for i := 0; i < nWorkers; i++ {
		go worker(done, taskCh, errors, wg)
	}

	go runTasks(done, taskCh, tasks)

	go func() {
		wg.Wait()
		close(errors)
	}()

	var errCounter int

	for err := range errors {
		if err != nil {
			errCounter++
			if errCounter == nErrors {
				close(done)
			}
		}
	}

	if errCounter >= nErrors {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func runTasks(done <-chan struct{}, taskCh chan<- Task, tasks []Task) {
	defer close(taskCh)

	for _, task := range tasks {
		select {
		case <-done:
			return
		case taskCh <- task:
		}
	}
}

func worker(done <-chan struct{}, taskCh <-chan Task, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range taskCh {
		select {
		case <-done:
			return
		default:
			errors <- task()
		}
	}
}
