// Package optional provides
package optional

import "errors"

var ErrNoSuchElement = errors.New("Nothing does not contain value")

type Function[T any, U any] func(arg T) U

type Predicate[T any] Function[T, bool]

func False[T any](value T) bool {
	return false
}

func True[T any](value T) bool {
	return true
}

// Optional is a wrapper that can contain value, nothing or error.
type Optional[T any] interface {
	IsPresent() bool
	IfPresent(func(T))
	Filter(p Predicate[T]) Optional[T]
	OrElse(other T) T
	Get() (T, error)
	Map(func(T) T) Optional[T]
}

// Optional subtype that contains a value
type Just[T any] struct {
	Optional[T]

	value T
}

func (j Just[T]) IsPresent() bool {
	return true
}

func (j Just[T]) IfPresent(f func(T)) {
	f(j.value)
}

func (j Just[T]) Filter(p Predicate[T]) Optional[T] {
	if p(j.value) {
		return j
	}
	return Nothing[T]{}
}

func (j Just[T]) OrElse(other T) T {
	return j.value
}

func (j Just[T]) Get() (T, error) {
	return j.value, nil
}

func (j Just[T]) Map(f func(T) T) Optional[T] {
	return Just[T]{
		value: f(j.value),
	}
}

// Optional subtype that contains nothing
type Nothing[T any] struct {
}

func (j Nothing[T]) IsPresent() bool {
	return false
}

func (j Nothing[T]) IfPresent(_ func(T)) {}

func (j Nothing[T]) Filter(p Predicate[T]) Optional[T] {
	return j
}

func (j Nothing[T]) OrElse(other T) T {
	return other
}

func (j Nothing[T]) Get() (T, error) {
	return *new(T), ErrNoSuchElement
}

func (j Nothing[T]) Map(func(T) T) Optional[T] {
	return j
}

// Optional subtype that contains an error
type Error[T any] struct {
	Nothing[T]

	err error
}

func (e Error[T]) Get() (T, error) {
	return *new(T), e.err
}

// Constructs an Optional from a value type
func Of[T any](value T) Optional[T] {
	return Just[T]{
		value: value,
	}
}

// Constructs an Optional from a value and error (e.g. to wrap a function result)
func OfError[T any](value T, err error) Optional[T] {
	if err != nil {
		return Error[T]{
			err: err,
		}
	}
	return Of(value)
}

// Constructs an Optional from a pointer type (Nothing if value is nil)
func OfNullable[T any](value *T) Optional[*T] {
	if value == nil {
		return Nothing[*T]{}
	}
	return Just[*T]{
		value: value,
	}
}

// Constructs an Optional from a pointer type (Nothing if value is nil)
func OfErrorNullable[T any](value *T, err error) Optional[*T] {
	if err != nil {
		return Error[*T]{
			err: err,
		}
	}
	return OfNullable(value)
}

// Applies a function to an optional if the function returns value and error
func MapErr[A any, B any](opt Optional[A], f func(A) (B, error)) Optional[B] {
	maybeValue, maybeErr := opt.Get()
	if !opt.IsPresent() {
		if maybeErr == ErrNoSuchElement {
			return Nothing[B]{}
		}
		return OfError(*new(B), maybeErr)
	}
	return OfError(f(maybeValue))
}

// Applies a function to an optional if the function returns value
func Map[A any, B any](opt Optional[A], f func(A) B) Optional[B] {
	return MapErr(opt, func(a A) (B, error) {
		return f(a), nil
	})
}
