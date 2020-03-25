package location

import "strings"

// Normalize transforms the human readable Azure Region/Location names (e.g. `West US`)
// into the canonical value to allow comparisons between user-code and API Responses
func Normalize(input string) string {
	return strings.Replace(strings.ToLower(input), " ", "", -1)
}

// NormalizeNilable normalizes the Location field even if it's nil to ensure this field
// can always have a value
func NormalizeNilable(input *string) string {
	if input == nil {
		return ""
	}

	return Normalize(*input)
}
