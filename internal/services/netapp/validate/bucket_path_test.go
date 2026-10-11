// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestBucketPath(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{input: "/", expected: true},
		{input: "/sub/path", expected: true},
		{input: "/sub/path/with-dashes", expected: true},
		{input: "", expected: false},
		{input: "sub/path", expected: false},
		{input: "//sub", expected: true},
		{input: "/sub\\path", expected: false},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.input)

		_, errors := BucketPath(v.input, "path")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t for %q (errors: %v)", v.expected, actual, v.input, errors)
		}
	}
}
