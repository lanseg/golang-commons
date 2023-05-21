package collections

import (
	"reflect"
	"testing"
)

type someType struct {
}

func TestTuples(t *testing.T) {

	t.Run("AsPair", func(t *testing.T) {
		result := AsPair("A", "B")
		want := Pair[string, string]{
			First:  "A",
			Second: "B",
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("AsPair() is expected to be %v, but got %v", want, result)
		}
	})

	t.Run("AsPair Nullable non null", func(t *testing.T) {
		result := AsPair(&someType{}, "B")
		want := Pair[*someType, string]{
			First:  &someType{},
			Second: "B",
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("AsPair() is expected to be %v, but got %v", want, result)
		}
	})

	t.Run("AsPair Nullable null", func(t *testing.T) {
		result := AsPair[*someType, string](nil, "B")
		want := Pair[*someType, string]{
			First:  nil,
			Second: "B",
		}
		if !reflect.DeepEqual(result, want) {
			t.Errorf("AsPair() is expected to be %v, but got %v", want, result)
		}
	})
}
