package collections

type TraverseOrder int

const (
	DepthFirst   = TraverseOrder(1)
	BreadthFirst = TraverseOrder(2)
)

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

type sliceStream[T any] struct {
	Stream[T]
	pos int

	slice []T
}

func (i *sliceStream[T]) HasNext() bool {
	return i.pos < len(i.slice)-1
}

func (i *sliceStream[T]) Next() (T, bool) {
	if !i.HasNext() {
		return *new(T), false
	}
	i.pos++
	return i.slice[i.pos], true
}

func (i *sliceStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *sliceStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *sliceStream[T]) Peek(peek func(item T)) Stream[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *sliceStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

func IterateSlice[T any](slice []T) Stream[T] {
	return &sliceStream[T]{
		pos:   -1,
		slice: slice,
	}
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
	Iterator[T]

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

// treeStream walks over a tree structure in a BFS way.
type treeStream[T any] struct {
	Stream[T]

	order       TraverseOrder
	toVisit     []T
	getChildren func(node T) []T
}

func (i *treeStream[T]) HasNext() bool {
	return len(i.toVisit) > 0
}

func (i *treeStream[T]) Next() (T, bool) {
	if !i.HasNext() {
		return *new(T), false
	}
	next, remain := i.toVisit[0], i.toVisit[1:]
	switch i.order {
	case BreadthFirst:
		i.toVisit = append(remain, i.getChildren(next)...)
	case DepthFirst:
		i.toVisit = append(i.getChildren(next), remain...)
	}
	return next, true
}

func (i *treeStream[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *treeStream[T]) Filter(predicate func(item T) bool) Stream[T] {
	return newFilteredStream[T](i, predicate)
}

func (i *treeStream[T]) Peek(peek func(item T)) Stream[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *treeStream[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

// IterateTree creates a treeStream for a root node "root" and a method to get node children.
// Uses BFS by default
func IterateTree[T any](root T, order TraverseOrder, getChildren func(T) []T) *treeStream[T] {
	return &treeStream[T]{
		order:       order,
		toVisit:     []T{root},
		getChildren: getChildren,
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
