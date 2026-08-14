// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package kusto

import "testing"

func TestTrustedExternalTenantsEqual(t *testing.T) {
	cases := []struct {
		name     string
		a        []interface{}
		b        []interface{}
		expected bool
	}{
		{
			name:     "both empty",
			a:        []interface{}{},
			b:        []interface{}{},
			expected: true,
		},
		{
			name:     "same order",
			a:        []interface{}{"tenant-a", "tenant-b", "tenant-c"},
			b:        []interface{}{"tenant-a", "tenant-b", "tenant-c"},
			expected: true,
		},
		{
			name:     "different order same set",
			a:        []interface{}{"tenant-a", "tenant-b", "tenant-c"},
			b:        []interface{}{"tenant-c", "tenant-a", "tenant-b"},
			expected: true,
		},
		{
			name:     "different lengths",
			a:        []interface{}{"tenant-a", "tenant-b"},
			b:        []interface{}{"tenant-a", "tenant-b", "tenant-c"},
			expected: false,
		},
		{
			name:     "disjoint sets",
			a:        []interface{}{"tenant-a", "tenant-b"},
			b:        []interface{}{"tenant-c", "tenant-d"},
			expected: false,
		},
		{
			name:     "same length different members",
			a:        []interface{}{"tenant-a", "tenant-b", "tenant-c"},
			b:        []interface{}{"tenant-a", "tenant-b", "tenant-d"},
			expected: false,
		},
		{
			name:     "wildcard and empty string reordered",
			a:        []interface{}{"*", ""},
			b:        []interface{}{"", "*"},
			expected: true,
		},
		{
			name:     "duplicate collapses to same set",
			a:        []interface{}{"tenant-a"},
			b:        []interface{}{"tenant-a", "tenant-a"},
			expected: true,
		},
		{
			name:     "duplicates reordered same set",
			a:        []interface{}{"tenant-a", "tenant-b", "tenant-a"},
			b:        []interface{}{"tenant-b", "tenant-a", "tenant-b"},
			expected: true,
		},
		{
			name:     "duplicate masking missing member",
			a:        []interface{}{"tenant-a", "tenant-a"},
			b:        []interface{}{"tenant-a", "tenant-b"},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := trustedExternalTenantsEqual(tc.a, tc.b); got != tc.expected {
				t.Fatalf("trustedExternalTenantsEqual(%v, %v) = %t, want %t", tc.a, tc.b, got, tc.expected)
			}
		})
	}
}
