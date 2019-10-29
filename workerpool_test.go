package workerpool

import (
	"errors"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	t.Parallel()

	const (
		maxQueues  = 4
		maxWorkers = 2
	)

	var (
		errFinished = errors.New("finished")
		started     = make(chan struct{}, maxQueues)
		finished    = make(chan struct{}, maxQueues)
		sync        = make(chan struct{})
	)

	wp, err := NewWorkerPool(
		MaxQueues(maxQueues),
		MaxWorkers(maxWorkers),
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < maxQueues; i++ {
		n := i
		wp.AddWorkerFunc(func() error {
			started <- struct{}{}
			<-sync
			finished <- struct{}{}
			if n%2 == 0 {
				return errFinished
			}
			return nil
		})
	}

	if n := wp.Queues(); n != maxQueues {
		t.Fatalf("got %d, want %d", n, maxQueues)
	}

	if n := wp.MaxWorkers(); n != maxWorkers {
		t.Fatalf("got %d, want %d", n, maxWorkers)
	}

	if n := wp.MaxQueues(); n != maxQueues {
		t.Fatalf("got %d, want %d", n, maxQueues)
	}

	close(sync)

	time.Sleep(time.Second)
	if n := len(finished); n != maxQueues {
		t.Fatalf("got %d, want %d", n, maxQueues)
	}

	if errs := wp.Wait(); len(errs) != 2 {
		t.Fatalf("got %v, want %v: %v", len(errs), 0, errs)
	}

	if errs := wp.Wait(); len(errs) != 2 {
		t.Fatalf("got %v, want %v: %v", len(errs), 0, errs)
	}
}

func TestOptions(t *testing.T) {
	t.Parallel()

	t.Run("default options", func(t *testing.T) {
		wp, err := NewWorkerPool()
		if err != nil {
			t.Fatal(err)
		}
		if wp.MaxQueues() != 1 {
			t.Errorf("got %d, want %d", wp.MaxQueues(), 1)
		}
		if wp.MaxWorkers() != 1 {
			t.Errorf("got %d, want %d", wp.MaxWorkers(), 1)
		}
	})

	t.Run("valid options", func(t *testing.T) {
		const (
			maxQueues  = 6
			maxWorkers = 2
		)
		wp, err := NewWorkerPool(
			MaxQueues(maxQueues),
			MaxWorkers(maxWorkers),
		)
		if err != nil {
			t.Fatal(err)
		}
		if wp.MaxQueues() != maxQueues {
			t.Errorf("got %d, want %d", wp.MaxQueues(), maxQueues)
		}
		if wp.MaxWorkers() != maxWorkers {
			t.Errorf("got %d, want %d", wp.MaxWorkers(), maxWorkers)
		}
	})

	t.Run("invalid options", func(t *testing.T) {
		if _, err := NewWorkerPool(
			MaxQueues(0),
		); err != ErrInvalidQueuesSize {
			t.Errorf("got %d, want %d", err, ErrInvalidQueuesSize)
		}
		if _, err := NewWorkerPool(
			MaxWorkers(0),
		); err != ErrInvalidWorkersSize {
			t.Errorf("got %d, want %d", err, ErrInvalidWorkersSize)
		}
	})

}
