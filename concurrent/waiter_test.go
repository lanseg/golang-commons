package concurrent

import (
	"testing"
	"time"

	"github.com/lanseg/golang-commons/optional"
)

const (
	testInterval = 100 * time.Millisecond
	testRetries  = 10
)

func newValueProvider[T any](value T) ValueProvider[T] {
	return func() optional.Optional[T] {
		return optional.Of(value)
	}
}

func newValueAfterRetries[T any](value T, setAfter int) ValueProvider[T] {
	return func() optional.Optional[T] {
		setAfter--
		if setAfter <= 0 {
			return optional.Of[T](value)
		}
		return optional.Nothing[T]{}
	}
}

func newNothingProvider[T any]() ValueProvider[T] {
	return func() optional.Optional[T] {
		return optional.Nothing[T]{}
	}
}

func newWaiter[T any](condition Condition[T], provider ValueProvider[T]) *Waiter[T] {
	return &Waiter[T]{
		interval:   testInterval,
		retryCount: testRetries,
		provider:   provider,
		condition:  condition,
	}
}

func TestWaitForSomething(t *testing.T) {

	t.Run("successful run got no value", func(t *testing.T) {
		result := newWaiter[string](
			NewIsPresent[string](),
			newNothingProvider[string]()).Wait()
		if result.IsPresent() {
			t.Errorf("Expected result not to be present, but got %v", result)
		}
	})

	t.Run("successful run got a value immediately", func(t *testing.T) {
		result := newWaiter[string](
			NewIsPresent[string](),
			newValueProvider[string]("SomeValue")).Wait()
		value, _ := result.Get()
		if value != "SomeValue" {
			t.Errorf("Expected to get value, but got %v", result)
		}
	})

	t.Run("successful run got a value after retries", func(t *testing.T) {
		result := newWaiter[string](
			NewIsPresent[string](),
			newValueAfterRetries[string]("SomeValue", 4)).Wait()
		value, _ := result.Get()
		if value != "SomeValue" {
			t.Errorf("Expected to get value, but got %v", result)
		}
	})
}
