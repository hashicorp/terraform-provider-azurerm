// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"testing"
)

func TestSubnetNamesChangePeeringKind(t *testing.T) {
	names := func(v ...string) []interface{} {
		out := make([]interface{}, 0, len(v))
		for _, s := range v {
			out = append(out, s)
		}
		return out
	}

	testData := []struct {
		name     string
		old      []interface{}
		new      []interface{}
		expected bool
	}{
		{
			name:     "complete virtual network peering becomes a subnet peering",
			old:      names(),
			new:      names("subnet1"),
			expected: true,
		},
		{
			name:     "subnet peering becomes a complete virtual network peering",
			old:      names("subnet1"),
			new:      names(),
			expected: true,
		},
		{
			name:     "a subnet is added to an existing subnet peering",
			old:      names("subnet1"),
			new:      names("subnet1", "subnet2"),
			expected: false,
		},
		{
			name:     "a subnet peering swaps which subnet it covers",
			old:      names("subnet1"),
			new:      names("subnet2"),
			expected: false,
		},
		{
			name:     "no subnet names before or after",
			old:      names(),
			new:      names(),
			expected: false,
		},
	}

	for _, v := range testData {
		t.Run(v.name, func(t *testing.T) {
			if actual := subnetNamesChangePeeringKind(v.old, v.new); actual != v.expected {
				t.Fatalf("expected %t but got %t for %v -> %v", v.expected, actual, v.old, v.new)
			}
		})
	}
}
