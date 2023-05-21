package collections

import (
	"reflect"
	"testing"
)

func TestSliceIterator(t *testing.T) {
	t.Run("SliceIterator.HasNext for an empty slice", func(t *testing.T) {
		slice := IterateSlice([]string{})
		if slice.HasNext() {
			t.Errorf("HasNext for an empty slice should return false")
		}
	})

	t.Run("SliceIterator.Next for an empty slice", func(t *testing.T) {
		slice := IterateSlice([]string{})
		_, ok := slice.Next()
		if ok {
			t.Errorf("Next for an empty slice should return nil, false")
		}
	})

	t.Run("SliceIterator ForEachRemainig an empty slice", func(t *testing.T) {
		slice := IterateSlice([]string{})
		callCount := 0
		slice.ForEachRemaining(func(s string) bool {
			callCount++
			return true
		})
		if callCount > 0 {
			t.Errorf("ForEachRemainig for an empty slice should not be called")
		}
	})

	t.Run("SliceIterator.HasNext for a slice", func(t *testing.T) {
		slice := IterateSlice([]string{"1", "2", "3", "4"})
		for _, expect := range []bool{true, true, true, true, false} {
			if slice.HasNext() != expect {
				t.Errorf("HasNext expected to return %v, but got %v",
					slice.HasNext, expect)
				break
			}
			slice.Next()
		}
	})

	t.Run("SliceIterator.Next for a slice", func(t *testing.T) {
		slice := IterateSlice([]string{"1", "2"})
		aValue, aOk := slice.Next()
		bValue, bOk := slice.Next()
		_, cOk := slice.Next()
		if !(aOk == true && bOk == true && aValue == "1" && bValue == "2") {
			t.Errorf("Next expected to return value, true for existing value")
		}
		if cOk {
			t.Errorf("Next expected to return _, false for existing value")
		}

	})

	t.Run("SliceIterator ForEachRemainig a slice", func(t *testing.T) {
		src := []string{"1", "2", "3", "4"}
		slice := IterateSlice(src)
		resultA := []string{}
		resultB := []string{}

		slice.ForEachRemaining(func(s string) bool {
			resultA = append(resultA, s)
			return true
		})
		slice.ForEachRemaining(func(s string) bool {
			resultB = append(resultB, s)
			return false
		})
		if !reflect.DeepEqual(resultA, []string{"1"}) || !reflect.DeepEqual(resultB, []string{"2", "3", "4"}) {
			t.Errorf("ForEachRemaining expected to iterate over all values, but got %v and %v", resultA, resultB)
		}
	})
}
