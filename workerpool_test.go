package workerpool

import (
	"errors"
	"reflect"
	"testing"
)

func TestExample(t *testing.T) {
	inputs := []string{"Lion", "Tiger", "Cheetah", "Jaguar", "Leopard", "Cat", "Cougar"}
	outputs := make([]string, len(inputs))
	wp := New(2)
	for n, v := range inputs {
		n := n
		v := v
		wp.AddJobFunc(func() error {
			outputs[n] = v
			return nil
		})
	}
	for _, err := range wp.Wait() {
		t.Error(err)
	}

	if !reflect.DeepEqual(inputs, outputs) {
		t.Errorf("\ngot: %#v \nwant: %#v", inputs, outputs)
	}
}

func TestMax(t *testing.T) {
	const want = 1
	var (
		errA = errors.New("error: A")
		errB = errors.New("error: B")
	)
	wp := New(0)
	if got := wp.max; got != want {
		t.Errorf("got: '%v', want: '%v'", got, want)
	}
	for _, err := range []error{errA, errB} {
		err := err
		wp.AddJobFunc(func() error {
			return err
		})
	}
	errs := wp.Wait()
	if got := len(errs); got != 2 {
		t.Errorf("got: '%v', want: '%v'", got, 2)
	}
	if got := errs[0]; !errors.Is(got, errA) {
		t.Errorf("got: '%v', want: '%v'", got, errA)
	}
	if got := errs[1]; !errors.Is(got, errB) {
		t.Errorf("got: '%v', want: '%v'", got, errB)
	}
}
