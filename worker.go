package workerpool

// Worker is an interface for executing a single function.
// Do is called by the worker pool.
type Worker interface {
	Do() error
}

// The WorkerFunc type is an adapter to allow the use of
// ordinary functions as workers. If f is a function
// with the appropriate signature, WorkerFunc(f) is a
// Worker that calls f.
type WorkerFunc func() error

// Do calls f().
func (f WorkerFunc) Do() error {
	return f()
}
