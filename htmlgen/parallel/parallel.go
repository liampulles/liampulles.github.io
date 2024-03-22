package parallel

import (
	"errors"
)

type Job func() error

// Run jobs concurrently. Result is joined errors, in arbitrary order.
func Concurrent(jobs []Job, concurrency int) error {
	active := make(chan struct{}, concurrency)
	done := make(chan error, len(jobs))

	// Fork them off
	for _, job := range jobs {
		go func() {
			// -> This will block if the active channel is full
			active <- struct{}{}
			err := job()
			<-active
			done <- err
		}()
	}

	// Wait for jobs, accumulate errors
	var err error
	pending := len(jobs)
	for pending > 0 {
		jErr := <-done
		err = errors.Join(err, jErr)
		pending--
	}

	return err
}
