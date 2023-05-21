package collections

type Iterator[T any] interface {
	HasNext() bool
	Next() (T, bool)
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
