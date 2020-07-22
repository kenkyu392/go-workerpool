package workerpool

import (
	"sync"
)

// New ...
func New(max int) *WorkerPool {
	if max < 1 {
		max = 1
	}
	wp := &WorkerPool{
		max: max,
		job: make(chan Job, max),
	}
	wp.wg.Add(wp.max)
	for i := 0; i < wp.max; i++ {
		go wp.dispatch()
	}
	return wp
}

// WorkerPool ...
type WorkerPool struct {
	wg   sync.WaitGroup
	job  chan Job
	max  int
	errs []error
	once sync.Once
}

// Errors ...
func (w *WorkerPool) Errors() []error {
	return w.errs
}

// AddJob ...
func (w *WorkerPool) AddJob(fn Job) {
	if fn != nil {
		w.job <- fn
	}
}

// AddJobFunc ...
func (w *WorkerPool) AddJobFunc(fn JobFunc) {
	w.AddJob(fn)
}

// Wait ...
func (w *WorkerPool) Wait() []error {
	w.once.Do(func() {
		close(w.job)
		w.wg.Wait()
	})
	return w.Errors()
}

func (w *WorkerPool) dispatch() {
	defer w.wg.Done()
	for {
		select {
		default:
		case worker, ok := <-w.job:
			if !ok {
				return
			}
			if err := worker.Do(); err != nil {
				w.errs = append(w.errs, err)
			}
		}
	}
}
