package utils

// NormalizeNilableFloat normalizes a nilable Float into a Float
// that is, if it's nil returns an empty Float else the value
func NormaliseNilableFloat64(input *float64) float64 {
	if input == nil {
		return 0
	}

	return *input
}
