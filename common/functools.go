package common

// Set creates a function that updates target to given value.
func Set[T any](target *T) func(T) {
	return func(value T) {
		*target = value
	}
}
