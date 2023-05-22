package collections

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
}

type SliceIterator[T any] struct {
	Iterator[[]T]
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

func IterateSlice[T any](slice []T) Iterator[T] {
	return &SliceIterator[T]{
		pos:   -1,
		slice: slice,
	}
}

// TreeIterator walks over a tree structure in a BFS way.
type TreeIterator[T any] struct {
	Iterator[T]

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
	i.toVisit = append(remain, i.getChildren(next)...)
	return next, true
}

func (i *TreeIterator[T]) ForEachRemaining(f func(item T) bool) {
	done := false
	for i.HasNext() && !done {
		result, _ := i.Next()
		done = f(result)
	}
}

// IterateTree creates a TreeIterator for a root node "root" and a method to get node children.
func IterateTree[T any](root T, getChildren func(T) []T) *TreeIterator[T] {
	return &TreeIterator[T]{
		toVisit:     []T{root},
		getChildren: getChildren,
	}
}
