package collections

import (
	"reflect"
	"sort"
	"testing"
)

func TestGroupBy(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		items    []string
		key      func(a string) string
		expected map[string]([]string)
	}{
		{
			desc:  "successful operation for a simple array",
			items: []string{"hello world", "hello there", "goodbye world"},
			key: func(a string) string {
				return string(a[0])
			},
			expected: map[string]([]string){
				"h": []string{"hello world", "hello there"},
				"g": []string{"goodbye world"},
			},
		},
		{
			desc:     "successful operation for an empty array",
			items:    []string{},
			key:      identity[string],
			expected: map[string]([]string){},
		},
		{
			desc:  "successful operation for an identity function key",
			items: []string{"hello world", "hello there", "goodbye world"},
			key:   identity[string],
			expected: map[string]([]string){
				"hello world":   []string{"hello world"},
				"hello there":   []string{"hello there"},
				"goodbye world": []string{"goodbye world"},
			},
		},
		{
			desc:  "duplicates should stay in the result",
			items: []string{"ab", "ab", "ab", "ba", "ba", "ba"},
			key: func(a string) string {
				return string(a[0])
			},
			expected: map[string]([]string){
				"a": []string{"ab", "ab", "ab"},
				"b": []string{"ba", "ba", "ba"},
			},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			result := GroupBy(tc.items, tc.key)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("GroupBy(%v) expected to be %v, but got %v", tc.items,
					tc.expected, result)
			}
		})
	}
}

func TestKeysValues(t *testing.T) {
	for _, tc := range []struct {
		desc       string
		data       map[string]string
		wantKeys   []string
		wantValues []string
	}{
		{
			desc:       "successful operation",
			data:       map[string]string{"a": "1", "b": "2", "c": "!@#", "hello": "world"},
			wantKeys:   []string{"a", "b", "c", "hello"},
			wantValues: []string{"1", "2", "!@#", "world"},
		},
		{
			desc:       "empty map returns empty key value",
			data:       map[string]string{},
			wantKeys:   []string{},
			wantValues: []string{},
		},
		{
			desc:       "duplicate values preserved",
			data:       map[string]string{"a": "1", "b": "1", "c": "1"},
			wantKeys:   []string{"a", "b", "c"},
			wantValues: []string{"1", "1", "1"},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			keys := Keys(tc.data)
			values := Values(tc.data)
			sort.Strings(keys)
			sort.Strings(values)
			sort.Strings(tc.wantKeys)
			sort.Strings(tc.wantValues)
			if !reflect.DeepEqual(tc.wantKeys, keys) {
				t.Errorf("Keys(%v) expected to be %v, but got %v", tc.data, tc.wantKeys, keys)
			}
			if !reflect.DeepEqual(tc.wantValues, values) {
				t.Errorf("Values(%v) expected to be %v, but got %v", tc.data, tc.wantValues, values)
			}
		})
	}
}

func TestNewMap(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		keys   []string
		values []string
		want   map[string]string
	}{
		{
			desc:   "successful operation",
			keys:   []string{"a", "b", "c"},
			values: []string{"0", "1", "2"},
			want:   map[string]string{"a": "0", "b": "1", "c": "2"},
		},
		{
			desc:   "more keys than values",
			keys:   []string{"a", "b", "c", "d"},
			values: []string{"0", "1"},
			want:   map[string]string{"a": "0", "b": "1", "c": "", "d": ""},
		},
		{
			desc:   "more values than keys",
			keys:   []string{"a", "b"},
			values: []string{"0", "1", "2", "3"},
			want:   map[string]string{"a": "0", "b": "1"},
		},
		{
			desc:   "empty keys and values",
			keys:   []string{},
			values: []string{},
			want:   map[string]string{},
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			result := NewMap(tc.keys, tc.values)
			if !reflect.DeepEqual(tc.want, result) {
				t.Errorf("NewMap(%v, %v) expected to be %v, but got %v", tc.keys, tc.values,
					tc.want, result)
			}
		})
	}
}

func TestSetContainsAll(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		tset  []string
		items []string
		want  bool
	}{
		{
			desc:  "set has exactly same items should return true",
			tset:  []string{"a", "b", "c", "d"},
			items: []string{"a", "b", "c", "d"},
			want:  true,
		},
		{
			desc:  "set has elements from items should return true",
			tset:  []string{"a", "b", "c", "d", "e"},
			items: []string{"b", "d", "c"},
			want:  true,
		},
		{
			desc:  "set contains some elements from items should return false",
			tset:  []string{"a", "b", "c", "d"},
			items: []string{"b", "c", "e"},
			want:  false,
		},
		{
			desc:  "set contains no elements from items should return false",
			tset:  []string{"a", "b", "c", "d"},
			items: []string{"g", "h", "i"},
			want:  false,
		},
		{
			desc:  "empty set should return false",
			tset:  []string{},
			items: []string{"b", "c", "e"},
			want:  false,
		},
        {
            desc: "empty items should return true",
            tset: []string{"1", "2", "3"},
            items: []string{},
            want: true,
        },
	} {
		t.Run(tc.desc, func(t *testing.T) {
			if tc.want != NewSet(tc.tset).ContainsAll(tc.items) {
				t.Errorf("Expected Set(%v).ContainsAll(%v) to be %v, but got %v",
					tc.tset, tc.items, tc.want, !tc.want)
			}
		})
	}
}
