// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateFirewallRuleName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "a8b-c",
			expected: true,
		},
		{
			input:    "AbcD",
			expected: true,
		},
		{
			input:    "_abc",
			expected: true,
		},
		{
			input:    "-abc",
			expected: true,
		},
		{
			input:    "ab_",
			expected: true,
		},
		{
			input:    "abc-",
			expected: true,
		},
		{
			input:    "ab",
			expected: true,
		},
		{
			input:    "a",
			expected: true,
		},
		{
			input:    "a@",
			expected: false,
		},
		{
			input:    strings.Repeat("s", 127),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 128),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 129),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := FirewallRuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
