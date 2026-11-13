// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestForwardingRuleName(t *testing.T) {
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
			input:    "rule",
			expected: true,
		},
		{
			// starts with a number
			input:    "1rule",
			expected: true,
		},
		{
			// contains a period
			input:    "rule.name",
			expected: true,
		},
		{
			// contains a hyphen and underscore
			input:    "rule-name_1",
			expected: true,
		},
		{
			// can't start with an underscore
			input:    "_rule",
			expected: false,
		},
		{
			// can't start with a hyphen
			input:    "-rule",
			expected: false,
		},
		{
			// can't start with a period
			input:    ".rule",
			expected: false,
		},
		{
			// can't end with a hyphen
			input:    "rule-",
			expected: false,
		},
		{
			// can't end with a period
			input:    "rule.",
			expected: false,
		},
		{
			// can't contain a space
			input:    "rule name",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "rule!",
			expected: false,
		},
		{
			// one character
			input:    "r",
			expected: true,
		},
		{
			// 80 characters
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzab",
			expected: true,
		},
		{
			// 81 characters
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabc",
			expected: false,
		},
	}

	for _, v := range testData {
		_, errors := ForwardingRuleName()(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
