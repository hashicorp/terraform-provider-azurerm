// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestLinuxAdminPassword(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			// bad: empty
			input:    "",
			expected: false,
		},
		{
			// bad: all lower characters
			input:    "juanjo",
			expected: false,
		},
		{
			// bad: only some caps
			input:    "JuanJo",
			expected: false,
		},
		{
			// bad: can't use reserved words
			input:    "P@$$w0rd",
			expected: false,
		},
		{
			// bad: can't be longer than 72 characters
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghiJklmno9q.stuvwxyzabcdefghijkjlmnopqrst",
			expected: false,
		},
		{
			// 72 characters its fine
			input:    "abcdefghijklmnopqrstuvwxyzabcD.9ghijklmnopqrstuvwxyzabcdefghijkjlmnopqrs",
			expected: true,
		},
		{
			// bad: can't be less than 6 characters
			input:    "A9.de",
			expected: false,
		},
		{
			// bad: "_" doesnt count as special character
			input:    "A9BC_7",
			expected: false,
		},
		{
			// "/" count as special character
			input:    "A9BC/7",
			expected: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := LinuxAdminPassword(v.input, "admin_password")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
