// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestDashboardName(t *testing.T) {
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
			input:    "hello#",
			expected: false,
		},
		{
			input:    "hello-",
			expected: true,
		},
		{
			input:    "hello-world",
			expected: true,
		},
		{
			input:    "hello9",
			expected: true,
		},
		{
			input:    strings.Repeat("s", 159),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 160),
			expected: true,
		},
		{
			input:    strings.Repeat("s", 161),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.input)

		_, errors := DashboardName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %t but got %t", v.expected, actual)
		}
	}
}
