// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// Bool
// Deprecated: please use the `To` function in the `pointer` package
func Bool(input bool) *bool {
	return &input
}

// Int32
// Deprecated: please use the `To` function in the `pointer` package
func Int32(input int32) *int32 {
	return &input
}

// Int64
// Deprecated: please use the `To` function in the `pointer` package
func Int64(input int64) *int64 {
	return &input
}

// Float
// Deprecated: please use the `To` function in the `pointer` package
func Float(input float64) *float64 {
	return &input
}

// String
// Deprecated: please use the`To` function in the `pointer` package
func String(input string) *string {
	return &input
}

// StringSlice
// Deprecated: please use the `To` function in the `pointer` package
func StringSlice(input []string) *[]string {
	if input == nil {
		return nil
	}
	return &input
}
