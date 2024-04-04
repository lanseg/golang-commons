package collections

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
}

type sliceIterator[T any] struct {
	Iterator[T]

	cursor int
	items  []T
}

func (i *sliceIterator[T]) HasNext() bool {
	return i.cursor < len(i.items)-1
}

func (i *sliceIterator[T]) Next() (T, bool) {
	if !i.HasNext() {
		return *new(T), false
	}
	i.cursor += 1
	return i.items[i.cursor], true
}

func IterateSlice[T any](items []T) Iterator[T] {
	return &sliceIterator[T]{
		cursor: -1,
		items:  items,
	}
}

type treeIterator[T any] struct {
	Iterator[T]

	order       TraverseOrder
	toVisit     []T
	getChildren func(node T) []T
}

func (i *treeIterator[T]) HasNext() bool {
	return len(i.toVisit) > 0
}

func (i *treeIterator[T]) Next() (T, bool) {
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

// IterateTree creates a treeIterator for a root node "root" and a method to get node children.
// Uses BFS by default
func IterateTree[T any](root T, order TraverseOrder, getChildren func(T) []T) Iterator[T] {
	return &treeIterator[T]{
		order:       order,
		toVisit:     []T{root},
		getChildren: getChildren,
	}
}
