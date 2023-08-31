// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestDatabricksVirtualNetworkPeeringName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// minimum length
			Input: "AA",
			Valid: true,
		},

		{
			// Valid ending
			Input: "Aa",
			Valid: true,
		},

		{
			// Valid ending
			Input: "A9",
			Valid: true,
		},

		{
			// Valid ending
			Input: "A_",
			Valid: true,
		},

		{
			// maximum length
			Input: "AAAAAAAAHHHHHHIIIIIIIIEEEEEEEEAAAAAAAHHHHHHHHHHIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// too short
			Input: "A",
			Valid: false,
		},

		{
			// too long
			Input: "AAAAAAAAHHHHHHIIIIIIIIEEEEEEEEAAAAAAAHHHHHHHHHHIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEE_",
			Valid: false,
		},

		{
			// invalid prefix character
			Input: "-AA",
			Valid: false,
		},

		{
			// invalid prefix character
			Input: "_AA",
			Valid: false,
		},

		{
			// invalid suffix character
			Input: "A-",
			Valid: false,
		},

		{
			// valid prefix character
			Input: "1A",
			Valid: true,
		},

		{
			// valid suffix character
			Input: "A1",
			Valid: true,
		},

		{
			// invalid substring character
			Input: "1A*A1",
			Valid: false,
		},

		{
			// valid substring characters
			Input: "1A.-_A1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := DatabricksVirtualNetworkPeeringName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
