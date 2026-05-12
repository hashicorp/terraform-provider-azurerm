// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import "testing"

func TestNumberOfIPAddressesDecreased(t *testing.T) {
	cases := []struct {
		name     string
		existing string
		expanded string
		expected bool
	}{
		{
			name:     "decrease from larger to smaller",
			existing: "128",
			expanded: "32",
			expected: true,
		},
		{
			name:     "increase from smaller to larger",
			existing: "32",
			expanded: "128",
			expected: false,
		},
		{
			name:     "equal values",
			existing: "64",
			expanded: "64",
			expected: false,
		},
		{
			name:     "lexicographic trap - increase looks like decrease",
			existing: "9",
			expanded: "10",
			expected: false,
		},
		{
			name:     "lexicographic trap - decrease looks like increase",
			existing: "10",
			expanded: "9",
			expected: true,
		},
		{
			name:     "large IPv6 address counts",
			existing: "340282366920938463463374607431768211456",
			expanded: "340282366920938463463374607431768211455",
			expected: true,
		},
		{
			name:     "large IPv6 address counts - increase",
			existing: "340282366920938463463374607431768211455",
			expanded: "340282366920938463463374607431768211456",
			expected: false,
		},
		{
			name:     "unparseable existing value",
			existing: "abc",
			expanded: "128",
			expected: false,
		},
		{
			name:     "unparseable expanded value",
			existing: "128",
			expanded: "xyz",
			expected: false,
		},
		{
			name:     "both unparseable",
			existing: "foo",
			expanded: "bar",
			expected: false,
		},
		{
			name:     "empty strings",
			existing: "",
			expanded: "",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := numberOfIPAddressesDecreased(tc.existing, tc.expanded)
			if result != tc.expected {
				t.Errorf("numberOfIPAddressesDecreased(%q, %q) = %v, want %v", tc.existing, tc.expanded, result, tc.expected)
			}
		})
	}
}
