// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBucketName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "ab",
			expected: false,
		},
		{
			input:    "abc",
			expected: true,
		},
		{
			input:    "my-bucket",
			expected: true,
		},
		{
			input:    "my.bucket.name",
			expected: true,
		},
		{
			input:    "1bucket",
			expected: true,
		},
		{
			input:    "Bucket",
			expected: false,
		},
		{
			input:    "-bucket",
			expected: false,
		},
		{
			input:    ".bucket",
			expected: false,
		},
		{
			input:    "bucket-",
			expected: false,
		},
		{
			input:    "bucket.",
			expected: false,
		},
		{
			input:    "my..bucket",
			expected: false,
		},
		{
			input:    "my.-bucket",
			expected: false,
		},
		{
			input:    "192.168.1.1",
			expected: false,
		},
		{
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijk",
			expected: true,
		},
		{
			input:    "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyzabcdefghijkl",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := BucketName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t for %q (errors: %v)", v.expected, actual, v.input, errors)
		}
	}
}
