package utils

import (
	"strings"
)

func SplitRemoveEmptyEntries(input string, delimiter string, removeWhitespace bool) []string {
	result := make([]string, 0)

	s := strings.Split(input, delimiter)

	for _, v := range s {
		if removeWhitespace {
			v = strings.TrimSpace(v)
		}
		if len(v) > 0 {
			result = append(result, v)
		}
	}

	return result
}
