package workerpool

import (
	"sync"
)

// NewWorkerPool ...
func NewWorkerPool(opts ...Option) (*WorkerPool, error) {
	wp := new(WorkerPool)
	for _, opt := range opts {
		if err := opt(wp); err != nil {
			return nil, err
		}
	}

	if wp.maxQueues < 1 {
		wp.maxQueues = 1
	}

	if wp.maxWorkers < 1 {
		wp.maxWorkers = 1
	}

	wp.qs = make(chan Worker, wp.maxQueues)
	wp.wg.Add(wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		go wp.dispatch()
	}

	return wp, nil
}

// WorkerPool ...
type WorkerPool struct {
	wg         sync.WaitGroup
	qs         chan Worker
	errs       []error
	closed     bool
	maxQueues  int
	maxWorkers int
}

// MaxQueues ...
func (w *WorkerPool) MaxQueues() int {
	return w.maxQueues
}

// MaxWorkers ...
func (w *WorkerPool) MaxWorkers() int {
	return w.maxWorkers
}

// Queues ...
func (w *WorkerPool) Queues() int {
	return len(w.qs)
}

// Errors ...
func (w *WorkerPool) Errors() []error {
	return w.errs
}

// AddWorker ...
func (w *WorkerPool) AddWorker(fn Worker) {
	if fn != nil {
		w.qs <- fn
	}
}

// AddWorkerFunc ...
func (w *WorkerPool) AddWorkerFunc(fn WorkerFunc) {
	w.AddWorker(fn)
}

// Wait ...
func (w *WorkerPool) Wait() []error {
	if w.closed {
		return w.Errors()
	}
	close(w.qs)
	w.wg.Wait()
	w.closed = true
	return w.Errors()
}

func (w *WorkerPool) dispatch() {
	defer w.wg.Done()
	for {
		select {
		case worker, ok := <-w.qs:
			if !ok {
				return
			}
			if err := worker.Do(); err != nil {
				w.errs = append(w.errs, err)
			}
		}
	}
}
