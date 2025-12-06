// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestExascaleDatabaseStorageVaultName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "v",
			expected: true,
		},
		{
			input:    "goodName",
			expected: true,
		},
		{
			input:    "_",
			expected: true,
		},
		{
			input:    "_good-Name-",
			expected: true,
		},
		{
			input:    "_G0od_name_",
			expected: true,
		},
		{
			input:    "good3Name",
			expected: true,
		},
		{
			input:    "good-name",
			expected: true,
		},
		{
			input:    "_G0od_name_",
			expected: true,
		},
		{
			input:    "Bad-Name2--",
			expected: false,
		},
		{
			input:    "1",
			expected: false,
		},
		{
			input:    "-Bad-Name2",
			expected: false,
		},
		{
			input:    "-Bad-Name2--",
			expected: false,
		},
		{
			input:    "--bad-Name",
			expected: false,
		},
		{
			input:    "2Bad-Name",
			expected: false,
		},
		{
			input:    "another--bad-name",
			expected: false,
		},
		{
			input:    "b@d2name",
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := ExascaleDatabaseResourceName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
