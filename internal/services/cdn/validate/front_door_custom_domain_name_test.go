// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorCustomDomainName(t *testing.T) {
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
			// maximum length
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEAAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEAAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: true,
		},

		{
			// too short
			Input: "A",
			Valid: false,
		},

		{
			// too long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEAAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEAAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE",
			Valid: false,
		},

		{
			// invalid prefix character
			Input: "-AA",
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
			Input: "1A_A1",
			Valid: false,
		},

		{
			// valid substring character
			Input: "1A-A1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorCustomDomainName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
