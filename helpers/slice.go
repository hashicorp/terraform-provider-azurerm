// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"slices"
)

// RemoveFromStringArray removes all matching values from a string array
func RemoveFromStringArray(elements []string, remove string) []string {
	if i := slices.Index(elements, remove); i != -1 {
		return RemoveFromStringArray(slices.Delete(elements, i, i+1), remove)
	}

	return elements
}

// SliceContainsValue
// Deprecated
// use slices.Contains instead
func SliceContainsValue(input []string, value string) bool {
	return slices.Contains(input, value)
}
