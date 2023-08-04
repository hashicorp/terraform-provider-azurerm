// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestSqlPoolName(t *testing.T) {
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
			input:    "aBc_123",
			expected: true,
		},
		{
			// UTF-8
			input:    "販売管理",
			expected: true,
		},
		{
			// can't contain hyphen
			input:    "ab-c",
			expected: false,
		},
		{
			// can't end with .
			input:    "abc.",
			expected: false,
		},
		{
			// can't end with space
			input:    "abc ",
			expected: false,
		},
		{
			// 60 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefgh",
			expected: true,
		},
		{
			// 61 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghi",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := SqlPoolName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
