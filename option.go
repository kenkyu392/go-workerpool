package workerpool

import "errors"

// Errors
var (
	ErrInvalidWorkersSize = errors.New("workers must be set to 1 or more")
	ErrInvalidQueuesSize  = errors.New("queues must be set to 1 or more")
)

// Option ...
type Option func(w *WorkerPool) error

// MaxWorkers ...
func MaxWorkers(n int) Option {
	return func(w *WorkerPool) error {
		if n < 1 {
			return ErrInvalidWorkersSize
		}
		w.maxWorkers = n
		return nil
	}
}

// MaxQueues ...
func MaxQueues(n int) Option {
	return func(w *WorkerPool) error {
		if n < 1 {
			return ErrInvalidQueuesSize
		}
		w.maxQueues = n
		return nil
	}
}
