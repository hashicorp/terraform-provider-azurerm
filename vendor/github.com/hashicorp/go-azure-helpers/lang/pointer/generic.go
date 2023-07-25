// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pointer

// From is a generic function that returns the value of a pointer
// If the pointer is nil, a zero value for the underlying type of the pointer is returned.
func From[T any](input *T) (output T) {
	var v T
	if input != nil {
		return *input
	}
	return v
}

// To is a generic function that returns a pointer to the value provided.
func To[T any](input T) *T {
	return &input
}
