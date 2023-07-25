// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDataStoreName(t *testing.T) {
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
			// can end with a hyphen
			input:    "hello-",
			expected: true,
		},
		{
			// cannot contain other special symbols other than hyphens
			input:    "hello.world",
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
			expected: true,
		},
		{
			// 255 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopa",
			expected: true,
		},
		{
			// 256 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqrstuvwxyzabcdefghabcdefghijklmnopqa",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := DataStoreName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
