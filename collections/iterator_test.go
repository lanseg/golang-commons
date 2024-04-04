package collections

import (
	"reflect"
	"testing"
)

func TestSliceIterator(t *testing.T) {

	for _, tc := range []struct {
		name  string
		items []any
		want  []any
	}{
		{"Simple slice", []any{1, 2, 3, 4}, []any{1, 2, 3, 4}},
		{"Empty slice", []any{}, []any{}},
		{"Nil slice", nil, []any{}},
		{"One element slice", []any{1}, []any{1}},
		{"Slice of interfaces", []any{"1", 2}, []any{"1", 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			i := IterateSlice(tc.items)
			result := []any{}
			for i.HasNext() {
				value, _ := i.Next()
				result = append(result, value)
			}

			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("Result after iterating (%v) expected to be (%v), but got (%v)",
					tc.items, tc.want, result)
			}

			if _, ok := i.Next(); ok {
				t.Errorf("Next after iterating (%v) expected to be (_, false, but got (_, true)",
					tc.items)
			}
		})
	}

	t.Run("Test multiple next after end", func(t *testing.T) {
		i := IterateSlice([]int{1, 2, 3})
		i.Next()
		i.Next()
		i.Next()
		i.Next()
		i.Next()
		if _, ok := i.Next(); ok {
			t.Errorf("Next after end should return false, but got true")
		}
		if i.HasNext() {
			t.Errorf("HasNext after end should return false, but got true")
		}
	})
}
