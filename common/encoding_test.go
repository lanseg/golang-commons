package common

import (
	"reflect"
	"testing"
)

type someStruct struct {
	value any
	child *someStruct
}

func TestSet(t *testing.T) {

	someA := &someStruct{
		value: 123,
		child: &someStruct{
			value: 456,
		},
	}

	for _, tc := range []struct {
		name   string
		src    any
		target any
	}{
		{
			name:   "primitive type set",
			src:    1,
			target: 2,
		},
		{
			name:   "primitive and nil set",
			src:    1,
			target: nil,
		},
		{
			name:   "nil and primitive set",
			src:    nil,
			target: "123",
		},
		{
			name:   "nil and nil set",
			src:    nil,
			target: nil,
		},
		{
			name:   "struct fields",
			src:    someA.value,
			target: someA.child.value,
		},
		{
			name:   "struct itself",
			src:    someA,
			target: nil,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			Set(&tc.target)(tc.src)
			if !reflect.DeepEqual(tc.target, tc.src) {
				t.Errorf("Expected target == src, but got %v and %v", tc.src, tc.target)
			}
		})
	}
}
