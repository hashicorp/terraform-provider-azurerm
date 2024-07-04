// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestComputeClusterName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// empty
			input:    "",
			expected: false,
		},
		{
			// basic example
			input:    "hello",
			expected: true,
		},
		{
			// cannot start with a hyphen
			input:    "-hello",
			expected: false,
		},
		{
			// cannot start with a digit
			input:    "1hello",
			expected: false,
		},
		{
			// can end with a digit
			input:    "hello2",
			expected: true,
		},
		{
			// cannot end with a hyphen
			input:    "hello-",
			expected: false,
		},
		{
			// cannot contain other special symbols other than hyphens
			input:    "hello.world",
			expected: false,
		},
		{
			// cannot contain underscore
			input:    "hello_world",
			expected: false,
		},
		{
			// hyphen in the middle
			input:    "hello-world",
			expected: true,
		},
		{
			// 2 chars
			input:    "ab",
			expected: false,
		},
		{
			// 3 chars
			input:    "abc",
			expected: true,
		},
		{
			// 32 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdef",
			expected: true,
		},
		{
			// 33 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefg",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := ComputeClusterName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
