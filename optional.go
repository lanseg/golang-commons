package optional

import "errors"

var ErrNoSuchElement = errors.New("Nothing contains value")

type Function[T any, U any] func(arg T) U

type Predicate[T any] Function[T, bool]

func False[T any](value T) bool {
	return false
}

func True[T any](value T) bool {
	return true
}

type Optional[T any] interface {
	IsPresent() bool
	Filter(p Predicate[T]) Optional[T]
	OrElse(other T) T
	Get() (T, error)
}

type Just[T any] struct {
	Optional[T]

	value T
}

func (j Just[T]) IsPresent() bool {
	return true
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

type Nothing[T any] struct {
}

func (j Nothing[T]) IsPresent() bool {
	return false
}

func (j Nothing[T]) Filter(p Predicate[T]) Optional[T] {
	return j
}

func (j Nothing[T]) OrElse(other T) T {
	return other
}

func (j Nothing[T]) Get() (T, error) {
	return *new(T), ErrNoSuchElement
}

type Error[T any] struct {
	Nothing[T]

	err error
}

func (e Error[T]) Get() (T, error) {
	return *new(T), e.err
}

func Of[T any](value T) Optional[T] {
	return Just[T]{
		value: value,
	}
}

func OfError[T any](value T, err error) Optional[T] {
	if err != nil {
		return Error[T]{
			err: err,
		}
	}
	return Of(value)
}

func OfNullable[T any](value *T) Optional[*T] {
	if value == nil {
		return Nothing[*T]{}
	}
	return Just[*T]{
		value: value,
	}
}

func OfErrorNullable[T any](value *T, err error) Optional[*T] {
	if err != nil {
		return Error[*T]{
			err: err,
		}
	}
	return OfNullable(value)
}
