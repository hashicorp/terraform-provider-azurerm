// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestWindowsAdminUsername(t *testing.T) {
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
			// can't be longer than 20 characters
			input:    "abcdefghijklmnopqrstuvwxyz",
			expected: false,
		},
		{
			// 20 characters its fine
			input:    "abcdefghijklmnopqrst",
			expected: true,
		},
		{
			// cannot end in a dot
			input:    "juanjo.",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := WindowsAdminUsername(v.input, "admin_username")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
