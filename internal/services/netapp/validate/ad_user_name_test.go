// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateNetAppActiveDirectoryUserName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "validuser",
			expected: true,
		},
		{
			input:    "user_with-characters!",
			expected: true,
		},
		{
			input:    "user@domain",
			expected: false,
		},
		{
			input:    "user\\domain",
			expected: false,
		},
		{
			input:    "",
			expected: false,
		},
		{
			// exactly 255 characters
			input:    strings.Repeat("a", 255),
			expected: true,
		},
		{
			// more than 255 characters
			input:    strings.Repeat("a", 256),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := ActiveDirectoryUserName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t for %q (errors: %v)", v.expected, actual, v.input, errors)
		}
	}
}
