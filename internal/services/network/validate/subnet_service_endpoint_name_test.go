// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestSubnetServiceEndpointName(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			input:    "",
			expected: false,
		},
		{
			input:    "Microsoft.Storage",
			expected: true,
		},
		{
			input:    "Microsoft.Storage.Global",
			expected: true,
		},
		{
			input:    "Microsoft.CognitiveServices",
			expected: true,
		},
		{
			input:    "Microsoft.KeyVault",
			expected: true,
		},
		{
			input:    "microsoft.storage",
			expected: false,
		},
		{
			input:    "Microsoft.Test",
			expected: false,
		},
	}

	validateFunc := SubnetServiceEndpointName()

	for _, tc := range tests {
		_, errs := validateFunc(tc.input, "service")
		ok := len(errs) == 0
		if ok != tc.expected {
			t.Fatalf("input %q: expected %t got %t", tc.input, tc.expected, ok)
		}
	}
}
