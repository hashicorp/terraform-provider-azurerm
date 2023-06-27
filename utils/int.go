// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

// NormaliseNilableInt64 takes a pointer to an int64 and returns a zero value or
// the real value if present
func NormaliseNilableInt64(input *int64) int64 {
	if input == nil {
		return 0
	}

	return *input
}
