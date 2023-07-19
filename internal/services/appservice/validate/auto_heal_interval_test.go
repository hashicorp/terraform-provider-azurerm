// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestAutoHealInterval(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "100:00:00",
		},
		{
			Input: "10:70:00",
		},
		{
			Input: "10:00:61",
		},
		{
			Input: "00:00:00",
			Valid: true,
		},
		{
			Input: "07:45:11",
			Valid: true,
		},
		{
			Input: "99:00:00",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := TimeInterval(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
