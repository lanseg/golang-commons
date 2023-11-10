package concurrent

import (
	"time"

	"github.com/lanseg/golang-commons/optional"
)

const (
	defaultWaiterInterval = 1 * time.Second
	defaultWaiterRetries  = 30
)

type ValueProvider[T any] func() optional.Optional[T]
type Condition[T any] func(optional.Optional[T]) bool

type waiter[T any] struct {
	retryCount int
	interval   time.Duration
	provider   ValueProvider[T]
	condition  Condition[T]
}

func (w *waiter[T]) Wait() optional.Optional[T] {
	ticker := time.NewTicker(w.interval)
	for i := 0; i < w.retryCount; i++ {
		select {
		case <-ticker.C:
			value := w.provider()
			if w.condition(value) {
                ticker.Stop()
				return value
			}
		}
	}
	return optional.Nothing[T]{}

}

func newWaiter[T any]() *waiter[T] {
	return &waiter[T]{
		retryCount: defaultWaiterRetries,
		interval:   defaultWaiterInterval,
	}
}

func WaitForSomething[T any](src ValueProvider[T]) optional.Optional[T] {
	waiter := newWaiter[T]()
	waiter.provider = src
	waiter.condition = func(o optional.Optional[T]) bool {
		return o.IsPresent()
	}
	return waiter.Wait()
}

func WaitForValue[T comparable](value T, src ValueProvider[T]) {
	waiter := newWaiter[T]()
	waiter.provider = src
	waiter.condition = func(o optional.Optional[T]) bool {
		return o.Filter(func(v T) bool {
			return v == value
		}).IsPresent()
	}
	waiter.Wait()
}
