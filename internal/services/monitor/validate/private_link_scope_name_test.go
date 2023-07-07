// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestPrivateLinkScopeName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "a",
			expected: true,
		},
		{
			input:    "A",
			expected: true,
		},
		{
			input:    "8",
			expected: true,
		},
		{
			input:    "a_-.b",
			expected: true,
		},
		{
			input:    "ab_",
			expected: true,
		},
		{
			input:    "ab-",
			expected: true,
		},
		{
			input:    "a+",
			expected: false,
		},
		{
			input:    "a.",
			expected: false,
		},
		{
			input:    "a.",
			expected: false,
		},
		{
			input:    strings.Repeat("s", 254),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 255),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 256),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := PrivateLinkScopeName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
