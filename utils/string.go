package utils

import "github.com/Azure/go-autorest/autorest/date"

// NormalizeNilableString normalizes a nilable string into a string
// that is, if it's nil returns an empty string else the value
func NormalizeNilableString(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func NormalizeTime(input *date.Time) string {
	if input != nil {
		return input.String()
	}
	return ""
}
