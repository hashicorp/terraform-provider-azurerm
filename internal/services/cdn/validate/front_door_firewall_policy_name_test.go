// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorFirewallPolicyName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Empty
			Input: "",
			Valid: false,
		},

		{
			// Min Len
			Input: "A",
			Valid: true,
		},

		{
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// Invalid Character
			Input: "1%A",
			Valid: false,
		},

		{
			// Start With Hyphen
			Input: "-1",
			Valid: false,
		},

		{
			// End With Hyphen
			Input: "1-",
			Valid: false,
		},

		{
			// Start With Letter, End With Letter
			Input: "AA",
			Valid: true,
		},

		{
			// Start With Number, End With Number
			Input: "11",
			Valid: false,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorFirewallPolicyName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Testing value %q, Expected %t but got %t", tc.Input, tc.Valid, valid)
		}
	}
}
