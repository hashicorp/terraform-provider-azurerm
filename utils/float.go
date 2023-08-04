// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// NormalizeNilableFloat normalizes a nilable float64 into a float64 value
func NormalizeNilableFloat(input *float64) float64 {
	if input == nil {
		return 0
	}

	return *input
}

// NormalizeNilableFloat32 normalizes a nilable float32 into a float32 value
func NormalizeNilableFloat32(input *float32) float32 {
	if input == nil {
		return 0
	}

	return *input
}
