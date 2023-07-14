// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

// NormalizeNilableString normalizes a nilable string into a string
// that is, if it's nil returns an empty string else the value
func NormalizeNilableString(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}
