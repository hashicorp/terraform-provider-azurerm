// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestEndpointName(t *testing.T) {
	testData := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "1",
			expected: true,
		},
		{
			input:    "a",
			expected: true,
		},
		{
			input:    "endpoint-1",
			expected: true,
		},
		{
			input:    "endpoint_1",
			expected: true,
		},
		{
			input:    "-endpoint",
			expected: false,
		},
		{
			input:    "_endpoint",
			expected: false,
		},
		{
			input:    "endpoint.1",
			expected: false,
		},
		{
			input:    strings.Repeat("a", 64),
			expected: true,
		},
		{
			input:    strings.Repeat("a", 65),
			expected: false,
		},
	}

	for _, v := range testData {
		_, errors := EndpointName(v.input, "name")
		actual := len(errors) == 0
		if v.expected != actual {
			t.Fatalf("Expected %q to validate as %t but got %t", v.input, v.expected, actual)
		}
	}
}
