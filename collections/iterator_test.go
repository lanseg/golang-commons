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

	gtFilter := func(i int) bool { return i > 4 }
	for _, tc := range []struct {
		name   string
		filter func(x int) bool
		src    []int
		want   []int
	}{
		{
			name:   "normal slice filtered",
			filter: gtFilter,
			src:    []int{2, 5, 7, 9, 3, 1, 4, 0, 6},
			want:   []int{5, 7, 9, 6},
		},
		{
			name:   "normal slice, none match",
			filter: gtFilter,
			src:    []int{1, 2, 3, 4, 0},
			want:   []int{},
		},
		{
			name:   "normal slice, all match",
			filter: gtFilter,
			src:    []int{5, 6, 7, 8, 9},
			want:   []int{5, 6, 7, 8, 9},
		},
		{
			name:   "empty slice",
			filter: gtFilter,
			src:    []int{},
			want:   []int{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := []int{}
			IterateSlice(tc.src).Filter(tc.filter).ForEachRemaining(func(i int) bool {
				result = append(result, i)
				return false
			})
			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("Result after filtering (%v) expected to be (%v), but got (%v)",
					tc.src, tc.want, result)
			}
		})
	}

	t.Run("SliceIterator empty slice Filter", func(t *testing.T) {
	})
}

type someTree struct {
	value    string
	children []*someTree
}

func TestTreeIterator(t *testing.T) {

	aTree := &someTree{
		"root",
		[]*someTree{
			{"1", []*someTree{}},
			{"2", []*someTree{
				{"4", []*someTree{
					{"6", []*someTree{}},
				}},
			}},
			{"3", []*someTree{{"5", []*someTree{}}}},
		},
	}

	want := []string{"root", "1", "2", "3", "4", "5", "6"}
	getChildren := func(t *someTree) []*someTree {
		return t.children
	}

	t.Run("Iterate normal tree, forEachRemaining", func(t *testing.T) {
		iterator := IterateTree(aTree, getChildren)

		result := []string{}
		iterator.ForEachRemaining(func(t *someTree) bool {
			result = append(result, t.value)
			return false
		})

		if !reflect.DeepEqual(result, want) {
			t.Errorf("ForEachRemaining expected to iterate and get (%v), but got (%v)", want, result)
		}
	})

	t.Run("Iterate normal tree, for loop", func(t *testing.T) {
		iterator := IterateTree(aTree, getChildren)

		result := []string{}
		for iterator.HasNext() {
			r, _ := iterator.Next()
			result = append(result, r.value)
		}

		if !reflect.DeepEqual(result, want) {
			t.Errorf("For loop expected to iterate and get (%v), but got (%v)", want, result)
		}
	})

	getLeaves := func(i *someTree) bool {
		return len(i.children) == 0
	}
	for _, tc := range []struct {
		name   string
		filter func(t *someTree) bool
		root   *someTree
		want   []*someTree
	}{
		{
			name:   "filter normal tree",
			filter: getLeaves,
			root:   aTree,
			want: []*someTree{
				{"1", []*someTree{}},
				{"5", []*someTree{}},
				{"6", []*someTree{}},
			},
		},
		{
			name: "filter none matches",
			filter: func(i *someTree) bool {
				return false
			},
			root: aTree,
			want: []*someTree{},
		},
		{
			name:   "filter empty tree",
			filter: getLeaves,
			root:   &someTree{"root", []*someTree{}},
			want:   []*someTree{{"root", []*someTree{}}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := []*someTree{}
			IterateTree(tc.root, getChildren).Filter(tc.filter).ForEachRemaining(
				func(n *someTree) bool {
					result = append(result, n)
					return false
				})
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Filtering of %v expected to return %v, but got %v",
					tc.root, tc.want, result)
			}
		})
	}

}

func TestFilterIterator(t *testing.T) {

	src := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	t.Run("Nested filtered iterators", func(t *testing.T) {
		result := []int{}
		want := []int{0, 6, 12}

		IterateSlice(src).Filter(func(i int) bool {
			return i%2 == 0
		}).Filter(func(i int) bool {
			return i%3 == 0
		}).ForEachRemaining(func(i int) bool {
			result = append(result, i)
			return false
		})
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Two filters of %v should return %v, but got %v", src, want, result)
		}

	})
}
