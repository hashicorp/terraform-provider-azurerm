package helpers

import "strings"

func trimURLScheme(input string) string {
	schemes := []string{
		"https://",
		"HTTPS://",
		"http://",
		"HTTP://",
	}

	for _, v := range schemes {
		if strings.HasPrefix(input, v) {
			return strings.TrimPrefix(input, v)
		}
	}

	return input
}
