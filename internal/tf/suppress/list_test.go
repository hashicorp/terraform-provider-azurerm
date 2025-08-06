// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package suppress

import "testing"

func TestStringSlicesAreEqual(t *testing.T) {
	cases := []struct {
		Name     string
		SliceA   []string
		SliceB   []string
		Suppress bool
	}{
		{
			Name:     "empty",
			SliceA:   []string{""},
			SliceB:   []string{""},
			Suppress: true,
		},
		{
			Name:     "same",
			SliceA:   []string{"value1", "value2"},
			SliceB:   []string{"value1", "value2"},
			Suppress: true,
		},
		{
			Name:     "different_order",
			SliceA:   []string{"value1", "value2"},
			SliceB:   []string{"value2", "value1"},
			Suppress: true,
		},
		{
			Name:     "different_values",
			SliceA:   []string{"value1", "value2"},
			SliceB:   []string{"value1", "value2", "value3"},
			Suppress: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if stringSlicesAreEqual(tc.SliceA, tc.SliceB) != tc.Suppress {
				t.Fatalf("Expected stringSlicesAreEqual to return %t for '%q' == '%q'", tc.Suppress, tc.SliceA, tc.SliceB)
			}
		})
	}
}
