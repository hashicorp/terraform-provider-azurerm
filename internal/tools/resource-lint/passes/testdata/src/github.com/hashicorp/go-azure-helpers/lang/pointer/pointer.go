// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package pointer

// Mock pointer helper functions for testing

func To[T any](v T) *T {
	return &v
}

func ToEnum[T any](v T) *T {
	return &v
}

func From[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}
