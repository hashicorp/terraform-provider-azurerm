// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestOrchestratedDomainNameLabel(t *testing.T) {
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
			input:    "a-b2-c",
			expected: true,
		},
		{
			// can't contain upper case
			input:    "A-B2-C",
			expected: false,
		},
		{
			// can't contain underscores
			input:    "a_b2_c",
			expected: false,
		},
		{
			// can't start with an underscore
			input:    "_hello",
			expected: false,
		},
		{
			// can't start with dot
			input:    ".hello",
			expected: false,
		},
		{
			// dot in middle
			input:    "hello.world",
			expected: false,
		},
		{
			// hyphen in middle
			input:    "hello-world",
			expected: true,
		},
		{
			// can't end with hyphen
			input:    "helloworld-",
			expected: false,
		},
		{
			// can't contain an exclamation mark
			input:    "hello!",
			expected: false,
		},
		{
			// can't end with dot
			input:    "hello.",
			expected: false,
		},
		{
			// can't end with underscore
			input:    "helloworld_",
			expected: false,
		},
		{
			// 26 characters
			input:    "abcdeabcdeabcdeabcdeabcdea",
			expected: true,
		},
		{
			// 27 characters
			input:    "abcdeabcdeabcdeabcdeabcdeab",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q...", v.input)

		_, errors := OrchestratedDomainNameLabel(v.input, "domain_name_label")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
