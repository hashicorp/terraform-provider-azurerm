// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestNetAppVolumeQuotaRuleName(t *testing.T) {
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
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't end with a dash
			input:    "hello-",
			expected: true,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// dash in the middle
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			// can't end with a period
			input:    "hello.",
			expected: false,
		},
		{
			// 63 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			expected: true,
		},
		{
			// 64 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkj",
			expected: true,
		},
		{
			// 65 chars
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkja",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := VolumeQuotaRuleName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
