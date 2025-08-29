// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestManagedRedisName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Invalid empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid To short",
			input:    "in",
			expected: false,
		},
		{
			name:     "Invalid characters underscores",
			input:    "invalid_Exports_Name",
			expected: false,
		},
		{
			name:     "Invalid characters space",
			input:    "invalid Managed Redis name",
			expected: false,
		},
		{
			name:     "Invalid name starts with hyphen",
			input:    "-invalidManagedRedisName",
			expected: false,
		},
		{
			name:     "Invalid name ends with hyphen",
			input:    "invalidManagedRedisName-",
			expected: false,
		},
		{
			name:     "Invalid name with consecutive hyphens",
			input:    "validStorageInsightConfigName--2",
			expected: false,
		},
		{
			name:     "Invalid name over max length",
			input:    "thisIsToLoooooooooooooooooooooo000oooooooongForAManagedRedisName",
			expected: false,
		},
		{
			name:     "Valid name max length",
			input:    "thisIsTheLooooooooooooooooooooooooooooongestManagedRedisName",
			expected: true,
		},
		{
			name:     "Valid name",
			input:    "validManagedRedisName",
			expected: true,
		},
		{
			name:     "Valid name with hyphen",
			input:    "validStorageInsightConfigName-2",
			expected: true,
		},
		{
			name:     "Valid name min length",
			input:    "val",
			expected: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, errors := ManagedRedisName(tc.input, "name")
			result := len(errors) == 0
			if result != tc.expected {
				t.Fatalf("Expected the result to be %v but got %v (and %d errors)", tc.expected, result, len(errors))
			}
		})
	}
}
