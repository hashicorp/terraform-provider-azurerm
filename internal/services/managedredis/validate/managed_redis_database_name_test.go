// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestManagedRedisDatabaseName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Invalid Empty database name",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid database name",
			input:    "My Database",
			expected: false,
		},
		{
			name:     "Valid database name",
			input:    "default",
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := ManagedRedisDatabaseName(tc.input, "name")
			actual := len(errors) == 0
			if tc.expected != actual {
				t.Fatalf("Expected %t but got %t", tc.expected, actual)
			}
		})
	}
}
