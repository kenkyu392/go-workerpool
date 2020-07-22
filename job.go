package workerpool

// Job is an interface for executing a single function.
// Do is called by the worker pool.
type Job interface {
	Do() error
}

// The JobFunc type is an adapter to allow the use of
// ordinary functions as jobs. If f is a function
// with the appropriate signature, JobFunc(f) is a
// Job that calls f.
type JobFunc func() error

// Do calls f().
func (f JobFunc) Do() error {
	return f()
}
