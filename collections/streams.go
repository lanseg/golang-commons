package collections

type TraverseOrder int

type Stream[T any] interface {
	Iterator[T]

	// ForEachRemaining invokes f for each element until the end of the underlying
	// sequence. If f returns true than iteration is done and we do not need any more
	// data.
	ForEachRemaining(f func(item T) bool)

	// Filter creates an iterator which provides only items that satisfy given predicate
	Filter(predicate func(item T) bool) Stream[T]

	// Peek returns an iterator that additionally performs given action on each element
	Peek(func(item T)) Stream[T]

	// Collect fetches all elements from an iterator and adds them to a slice
	Collect() []T
}

func StreamIterator[T any](iter Iterator[T]) Stream[T] {
	return &IteratorStream[T]{iterator: iter}
}

func SliceStream[T any](slice []T) Stream[T] {
	return StreamIterator[T](IterateSlice(slice))
}

// Base stream with the most common definitions
type IteratorStream[T any] struct {
	Stream[T]

	iterator Iterator[T]
}

func (i *IteratorStream[T]) HasNext() bool {
	return i.iterator.HasNext()
}

func (i *IteratorStream[T]) Next() (T, bool) {
	return i.iterator.Next()
}

func (i *IteratorStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *IteratorStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *IteratorStream[T]) Peek(peek func(item T)) Stream[T] {
	return &peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *IteratorStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

type peeker[T any] struct {
	Stream[T]

	peek   func(item T)
	parent Stream[T]
}

func (i peeker[T]) HasNext() bool {
	return i.parent.HasNext()
}

func (i peeker[T]) Next() (T, bool) {
	next, ok := i.parent.Next()
	if ok {
		i.peek(next)
	}
	return next, ok
}

func (i peeker[T]) ForEachRemaining(f func(item T) bool) {
	i.parent.ForEachRemaining(func(item T) bool {
		i.peek(item)
		return f(item)
	})
}

func (i peeker[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream(i, predicate)
}

func (i peeker[T]) Peek(peek func(item T)) Stream[T] {
	return &peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i peeker[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

type filterStream[T any] struct {
	Stream[T]

	ready bool
	value T

	predicate func(item T) bool
	parent    Stream[T]
}

func (i *filterStream[T]) takeUntilMatch() (T, bool) {
	done := false
	for !done {
		value, ok := i.parent.Next()
		if ok && i.predicate(value) {
			return value, true
		}
		done = !ok
	}
	return *new(T), false
}

func (i *filterStream[T]) HasNext() bool {
	if i.ready {
		return true
	}
	if !i.parent.HasNext() {
		return false
	}
	result, ok := i.takeUntilMatch()
	if !ok {
		return false
	}
	i.ready = true
	i.value = result
	return true
}

func (i *filterStream[T]) Next() (T, bool) {
	if i.ready {
		i.ready = false
		return i.value, true
	}
	return i.takeUntilMatch()
}

func (i *filterStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for !done {
		value, ok := i.Next()
		done = !ok || f(value)
	}
}

func (i *filterStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *filterStream[T]) Peek(peek func(item T)) Stream[T] {
	return &peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *filterStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

func newFilteredStream[T any](parent Stream[T], predicate func(item T) bool) Stream[T] {
	return &filterStream[T]{
		predicate: predicate,
		parent:    parent,
	}
}
