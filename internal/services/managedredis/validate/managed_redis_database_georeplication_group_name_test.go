// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestManagedRedisDatabaseGeoreplicationGroupName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid: single character",
			input:    "a",
			expected: true,
		},
		{
			name:     "Valid: alphanumeric",
			input:    "group1",
			expected: true,
		},
		{
			name:     "Valid: with hyphens",
			input:    "my-group-1",
			expected: true,
		},
		{
			name:     "Valid: starts with number",
			input:    "1group",
			expected: true,
		},
		{
			name:     "Valid: ends with number",
			input:    "group1",
			expected: true,
		},
		{
			name:     "Valid: 63 characters",
			input:    strings.Repeat("a", 63),
			expected: true,
		},
		{
			name:     "Valid: mixed case",
			input:    "MyGroup1",
			expected: true,
		},
		{
			name:     "Invalid: empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Invalid: too long (64 characters)",
			input:    strings.Repeat("a", 64),
			expected: false,
		},
		{
			name:     "Invalid: starts with hyphen",
			input:    "-group",
			expected: false,
		},
		{
			name:     "Invalid: ends with hyphen",
			input:    "group-",
			expected: false,
		},
		{
			name:     "Invalid: consecutive hyphens",
			input:    "my--group",
			expected: false,
		},
		{
			name:     "Invalid: contains underscore",
			input:    "my_group",
			expected: false,
		},
		{
			name:     "Invalid: contains space",
			input:    "my group",
			expected: false,
		},
		{
			name:     "Invalid: contains special characters",
			input:    "my@group",
			expected: false,
		},
		{
			name:     "Invalid: multiple consecutive hyphens",
			input:    "my---group",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, errors := ManagedRedisDatabaseGeoreplicationGroupName(testCase.input, "test")
			valid := len(errors) == 0

			if valid != testCase.expected {
				t.Errorf("Expected %t for input %q, got %t. Errors: %v", testCase.expected, testCase.input, valid, errors)
			}
		})
	}
}

func TestManagedRedisDatabaseGeoreplicationGroupName_WrongType(t *testing.T) {
	_, errors := ManagedRedisDatabaseGeoreplicationGroupName(123, "test")
	if len(errors) == 0 {
		t.Error("Expected error for non-string input")
	}
	if !strings.Contains(errors[0].Error(), "expected type of \"test\" to be string") {
		t.Errorf("Expected type error, got: %v", errors[0])
	}
}
