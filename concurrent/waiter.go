package concurrent

import (
	"time"

	"github.com/lanseg/golang-commons/optional"
)

const (
	defaultWaiterInterval = 1 * time.Second
	defaultWaiterRetries  = 10
)

type ValueProvider[T any] func() optional.Optional[T]
type Condition[T any] func(optional.Optional[T]) bool

func NewIsPresent[T any]() Condition[T] {
	return func(o optional.Optional[T]) bool {
		return o.IsPresent()
	}
}

func NewEquals[T comparable](value T) Condition[T] {
	return func(o optional.Optional[T]) bool {
		return o.Filter(func(v T) bool {
			return v == value
		}).IsPresent()
	}
}

type Waiter[T any] struct {
	retryCount int
	interval   time.Duration
	provider   ValueProvider[T]
	condition  Condition[T]
}

func (w *Waiter[T]) Wait() optional.Optional[T] {
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
	ticker.Stop()
	return optional.Nothing[T]{}
}

// WaitForSomething blocks until src returns optional with value or with error.
// By default, it tries 30 times with 1 second interval between retries
func WaitForSomething[T any](src ValueProvider[T]) optional.Optional[T] {
	return (&Waiter[T]{
		condition:  NewIsPresent[T](),
		provider:   src,
		interval:   defaultWaiterInterval,
		retryCount: defaultWaiterRetries,
	}).Wait()
}

// WaitForValue blocks until src returns optional with given value or with error.
// By default, it tries 30 times with 1 second interval between retries
func WaitForValue[T comparable](value T, src ValueProvider[T]) {
	(&Waiter[T]{
		condition:  NewEquals[T](value),
		provider:   src,
		interval:   defaultWaiterInterval,
		retryCount: defaultWaiterRetries,
	}).Wait()
}

func WaitForValueRetries[T comparable](value T, src ValueProvider[T], retryCount int) {
	(&Waiter[T]{
		condition:  NewEquals[T](value),
		provider:   src,
		interval:   defaultWaiterInterval,
		retryCount: retryCount,
	}).Wait()
}
