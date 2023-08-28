// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestBackupPolicyName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "hello",
			expected: true,
		},
		{
			input:    "hello_",
			expected: false,
		},
		{
			input:    "hello-",
			expected: true,
		},
		{
			input:    "malcolm-in-the-middle",
			expected: true,
		},
		{
			input:    strings.Repeat("s", 149),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 150),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 151),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := BackupPolicyName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
