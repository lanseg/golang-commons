package common

import (
	"testing"
)

type Something struct {
    value int
}

func TestIfNull (t *testing.T) {
    someA  := &Something { value: 1}
    someB := &Something { value : 2} 
	for _, tc := range []struct {
		desc     string
        a *Something
        b *Something
        want *Something
	}{
		{
			desc:  "both nil return nil",
            a: nil,
            b: nil,
            want: nil,
		},
        {
            desc: "a non nil, b nil returns a",
            a: someA,
            b: nil,
            want: someA,
        },
        {
            desc: "a nil, b non nil returns b",
            a: nil,
            b: someB,
            want: someB,
        },
        {
            desc: "a non nil, b non nil, returns a",
            a: someA,
            b: someB,
            want: someA,
        },
	} {
		t.Run(tc.desc, func(t *testing.T) {
			result := IfNull(tc.a, tc.b)
			if result != tc.want {
				t.Errorf("IfNull(%v, %v) expected to be %v, but got %v", tc.a, tc.b,
					tc.want, result)
			}
		})
	}
}

