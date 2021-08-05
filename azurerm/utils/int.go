package utils

// NormaliseNilableInt takes a pointer to an int and returns a zero value or
// the real value if present
func NormaliseNilableInt(input *int) int {
	if input == nil {
		return 0
	}

	return *input
}

// NormaliseNilableInt32 takes a pointer to an int32 and returns a zero value or
// the real value if present
func NormaliseNilableInt32(input *int32) int32 {
	if input == nil {
		return 0
	}

	return *input
}
