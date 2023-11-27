package collections

type TraverseOrder int

const (
	DepthFirst   = TraverseOrder(1)
	BreadthFirst = TraverseOrder(2)
)

// Iterator provides an iterator for a sequential data that cannot be put into a slice
// like custom list implementations or sequences with no known length.
type Iterator[T any] interface {
	// HasNext return true if there is still more data to proceed
	HasNext() bool

	// Next returns current value, true if the value exists and advances by one step.
	// Returns _, false if there was no value.
	Next() (T, bool)

	// ForEachRemaining invokes f for each element until the end of the underlying
	// sequence. If f returns true than iteration is done and we do not need any more
	// data.
	ForEachRemaining(f func(item T) bool)

	// Filter creates an iterator which provides only items that satisfy given predicate
	Filter(predicate func(item T) bool) Iterator[T]

	// Peek returns an iterator that additionally performs given action on each element
	Peek(func(item T)) Iterator[T]

	// Collect fetches all elements from an iterator and adds them to a slice
	Collect() []T
}

type Peeker[T any] struct {
	Iterator[T]

	peek   func(item T)
	parent Iterator[T]
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

func (i Peeker[T]) Filter(predicate func(item T) bool) Iterator[T] {
	return newFilteredIterator[T](i, predicate)
}

func (i Peeker[T]) Peek(peek func(item T)) Iterator[T] {
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

type FilterIterator[T any] struct {
	Iterator[T]

	ready bool
	value T

	predicate func(item T) bool
	parent    Iterator[T]
}

func (i *FilterIterator[T]) takeUntilMatch() (T, bool) {
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

func (i *FilterIterator[T]) HasNext() bool {
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

func (i *FilterIterator[T]) Next() (T, bool) {
	if i.ready {
		i.ready = false
		return i.value, true
	}
	return i.takeUntilMatch()
}

func (i *FilterIterator[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for !done {
		value, ok := i.Next()
		done = !ok || f(value)
	}
}

func (i *FilterIterator[T]) Filter(predicate func(item T) bool) Iterator[T] {
	return newFilteredIterator[T](i, predicate)
}

func (i *FilterIterator[T]) Peek(peek func(item T)) Iterator[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *FilterIterator[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

func newFilteredIterator[T any](parent Iterator[T], predicate func(item T) bool) Iterator[T] {
	return &FilterIterator[T]{
		predicate: predicate,
		parent:    parent,
	}
}

type SliceIterator[T any] struct {
	Iterator[T]
	pos int

	slice []T
}

func (i *SliceIterator[T]) HasNext() bool {
	return i.pos < len(i.slice)-1
}

func (i *SliceIterator[T]) Next() (T, bool) {
	if !i.HasNext() {
		return *new(T), false
	}
	i.pos++
	return i.slice[i.pos], true
}

func (i *SliceIterator[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *SliceIterator[T]) Filter(predicate func(item T) bool) Iterator[T] {
	return newFilteredIterator[T](i, predicate)
}

func (i *SliceIterator[T]) Peek(peek func(item T)) Iterator[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *SliceIterator[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

func IterateSlice[T any](slice []T) Iterator[T] {
	return &SliceIterator[T]{
		pos:   -1,
		slice: slice,
	}
}

// TreeIterator walks over a tree structure in a BFS way.
type TreeIterator[T any] struct {
	Iterator[T]

	order       TraverseOrder
	toVisit     []T
	getChildren func(node T) []T
}

func (i *TreeIterator[T]) HasNext() bool {
	return len(i.toVisit) > 0
}

func (i *TreeIterator[T]) Next() (T, bool) {
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

func (i *TreeIterator[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

func (i *TreeIterator[T]) Filter(predicate func(item T) bool) Iterator[T] {
	return newFilteredIterator[T](i, predicate)
}

func (i *TreeIterator[T]) Peek(peek func(item T)) Iterator[T] {
	return &Peeker[T]{
		peek:   peek,
		parent: i,
	}
}

func (i *TreeIterator[T]) Collect() []T {
	result := []T{}
	i.ForEachRemaining(func(item T) bool {
		result = append(result, item)
		return false
	})
	return result
}

// IterateTree creates a TreeIterator for a root node "root" and a method to get node children.
// Uses BFS by default
func IterateTree[T any](root T, order TraverseOrder, getChildren func(T) []T) *TreeIterator[T] {
	return &TreeIterator[T]{
		order:       order,
		toVisit:     []T{root},
		getChildren: getChildren,
	}
}
