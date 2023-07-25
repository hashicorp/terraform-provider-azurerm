// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestLinuxAdminUsername(t *testing.T) {
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
			input:    "juanjo",
			expected: true,
		},
		{
			// basic example with caps
			input:    "JuanJo",
			expected: true,
		},
		{
			// can't use reserved words
			input:    "administrator",
			expected: false,
		},
		{
			// can't be longer than 64 characters
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkjl",
			expected: false,
		},
		{
			// 64 characters its fine
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkj",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := LinuxAdminUsername(v.input, "admin_username")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
