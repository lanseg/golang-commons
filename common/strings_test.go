package common

import (
	"testing"
)

func TestEllipsis(t *testing.T) {

	for _, tc := range []struct {
		name         string
		in           string
		maxLen       int
		breakOnSpace bool
		want         string
	}{
		{
			name:   "normal ellipsis",
			in:     "some string to truncate",
			maxLen: 10,
			want:   "some st...",
		},
		{
			name:         "normal ellipsis break on space",
			in:           "some string to truncate",
			maxLen:       10,
			breakOnSpace: true,
			want:         "some...",
		},
		{
			name:         "normal ellipsis break on space no spaces",
			in:           "somestringtotruncate",
			maxLen:       15,
			breakOnSpace: true,
			want:         "...",
		},
		{
			name:   "normal ellipsis no truncation needed",
			in:     "some text to truncate",
			maxLen: 1000,
			want:   "some text to truncate",
		},
		{
			name:   "normal ellipsis small max length",
			in:     "some text to truncate",
			maxLen: 2,
			want:   "",
		},
		{
			name:   "normal ellipsis empty string",
			in:     "",
			maxLen: 10,
			want:   "",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result := Ellipsis(tc.in, tc.maxLen, tc.breakOnSpace)
			if result != tc.want {
				t.Errorf("Expected Ellipsis(%s, %d, %v) = %s, but got %s",
					tc.in, tc.maxLen, tc.breakOnSpace, tc.want, result)
			}
		})
	}
}
