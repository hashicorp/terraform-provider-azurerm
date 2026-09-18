// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"reflect"
	"testing"
)

func TestRemoveFromStringArray(t *testing.T) {
	testCases := []struct {
		name     string
		elements []string
		remove   string
		expected []string
	}{
		{
			name:     "element not found",
			elements: []string{"a", "b", "c"},
			remove:   "z",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "element found once",
			elements: []string{"a", "b", "c"},
			remove:   "b",
			expected: []string{"a", "c"},
		},
		{
			name:     "element found multiple times",
			elements: []string{"a", "b", "b", "c"},
			remove:   "b",
			expected: []string{"a", "c"},
		},
		{
			name:     "element at beginning",
			elements: []string{"a", "b", "c"},
			remove:   "a",
			expected: []string{"b", "c"},
		},
		{
			name:     "element at end",
			elements: []string{"a", "b", "c"},
			remove:   "c",
			expected: []string{"a", "b"},
		},
		{
			name:     "empty slice",
			elements: []string{},
			remove:   "a",
			expected: []string{},
		},
		{
			name:     "nil slice",
			elements: nil,
			remove:   "a",
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := RemoveFromStringArray(tc.elements, tc.remove)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected: %v, got: %v", tc.expected, actual)
			}
		})
	}
}
