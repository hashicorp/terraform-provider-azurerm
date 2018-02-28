// Copyright (c) 2017, A. Stoewer <adrian.stoewer@rz.ifi.lmu.de>
// All rights reserved.

package strcase

import (
	"strings"
)

// SnakeCase converts a string into snake case.
func SnakeCase(s string) string {
	return lowerDelimiterCase(s, '_')
}

// lowerDelimiterCase converts a string into snake_case or kebab-case depending on
// the delimiter passed in as second argument.
func lowerDelimiterCase(s string, delimiter rune) string {
	s = strings.TrimSpace(s)
	buffer := make([]rune, 0, len(s)+3)

	var prev rune
	var curr rune
	for _, next := range s {
		if isDelimiter(curr) {
			if !isDelimiter(prev) {
				buffer = append(buffer, delimiter)
			}
		} else if isUpper(curr) {
			if isLower(prev) || (isUpper(prev) && isLower(next)) {
				buffer = append(buffer, delimiter)
			}
			buffer = append(buffer, toLower(curr))
		} else if curr != 0 {
			buffer = append(buffer, curr)
		}
		prev = curr
		curr = next
	}

	if len(s) > 0 {
		if isUpper(curr) && isLower(prev) && prev != 0 {
			buffer = append(buffer, delimiter)
		}
		buffer = append(buffer, toLower(curr))
	}

	return string(buffer)
}
