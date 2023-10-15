package common

import (
	"strings"
	"unicode"
)

// Ellipsis truncates string to be at most maxLength runes long
// maxLength = resultStringLength + len("...")
func Ellipsis(s string, maxLength int, breakOnSpace bool) string {
	if maxLength < 3 {
		return ""
	}
	maxTextLength := maxLength - 3
	if len(s) < maxLength {
		return s
	}
	sb := strings.Builder{}
	lastSpace := 0
	prevRune := ' '
	for i, r := range s {
		if unicode.IsSpace(r) && !unicode.IsSpace(prevRune) {
			lastSpace = i
		}
		if i >= maxTextLength {
			break
		}
		sb.WriteRune(r)
		prevRune = r
	}
	result := sb.String()
	if breakOnSpace {
		result = result[:lastSpace]
	}
	return result + "..."
}
