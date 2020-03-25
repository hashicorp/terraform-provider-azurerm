package location

import "strings"

// NormalizeLocation transforms the human readable Azure Region/Location names (e.g. `West US`)
// into the canonical value to allow comparisons between user-code and API Responses
func NormalizeLocation(input string) string {
	return strings.Replace(strings.ToLower(input), " ", "", -1)
}

// NormalizeNilableLocation normalizes the Location field even if it's nil to ensure this field
// can always have a value
func NormalizeNilableLocation(input *string) string {
	if input == nil {
		return ""
	}

	return NormalizeLocation(*input)
}
