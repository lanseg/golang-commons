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

// Stream definitions
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
	return &Peeker[T]{
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

type Peeker[T any] struct {
	Stream[T]

	peek   func(item T)
	parent Stream[T]
}

func (i Peeker[T]) HasNext() bool {
	return i.parent.HasNext()
}

func (i Peeker[T]) Next() (T, bool) {
	next, ok := i.parent.Next()
	if ok {
		i.peek(next)
	}
	return next, ok
}

func (i Peeker[T]) ForEachRemaining(f func(item T) bool) {
	i.parent.ForEachRemaining(func(item T) bool {
		i.peek(item)
		return f(item)
	})
}

func (i Peeker[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i Peeker[T]) Peek(peek func(item T)) Stream[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i Peeker[T]) Collect() []T {
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
	return &Peeker[T]{
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

type multiStream[T any] struct {
	Stream[T]

	current   int
	iterators []Stream[T]
}

func (i *multiStream[T]) HasNext() bool {
	for _, is := range i.iterators {
		if is.HasNext() {
			return true
		}
	}
	return false
}

func (i *multiStream[T]) Next() (T, bool) {
	for cur := 0; cur < len(i.iterators); cur++ {
		iterIndex := (cur + i.current) % len(i.iterators)
		iter := i.iterators[iterIndex]
		if iter.HasNext() {
			i.current = (iterIndex + 1) % len(i.iterators)
			return iter.Next()
		}
	}
	return *new(T), false
}

func (i *multiStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *multiStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *multiStream[T]) Peek(peek func(item T)) Stream[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *multiStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

// Union joins multiple iterator to work as stream-of-streams.
//
// Items are picked one by one from each iterator: first item of first iterator, first item of
// second iterator..., second item of first iterator, second item of second iterator...
func Union[T any](streams ...Stream[T]) Stream[T] {
	return &multiStream[T]{
		iterators: streams,
	}
}

type mergedStream[T any] struct {
	Stream[T]

	iterators []Stream[T]
}

func (i *mergedStream[T]) HasNext() bool {
	for _, is := range i.iterators {
		if is.HasNext() {
			return true
		}
	}
	return false
}

func (i *mergedStream[T]) Next() (T, bool) {
	for _, is := range i.iterators {
		if is.HasNext() {
			return is.Next()
		}
	}
	return *new(T), false
}

func (i *mergedStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *mergedStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *mergedStream[T]) Peek(peek func(item T)) Stream[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *mergedStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

// Concat joins multiple iterators and traverses them sequentially
// When reaching last element of an iterator, starting the next iterator
func Concat[T any](i ...Stream[T]) Stream[T] {
	return &mergedStream[T]{
		iterators: i,
	}
}
