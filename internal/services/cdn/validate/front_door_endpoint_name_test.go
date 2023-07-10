// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorEndpointName(t *testing.T) {
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
			// Max Len
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIE",
			Valid: true,
		},

		{
			// Too Long
			Input: "AAAAAAAAAAAAAHHHHHHHHHHHHHHHHIIIIIIIIIIIIIIIIIE",
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
			// Too Short
			Input: "1",
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
			Valid: true,
		},

		{
			// Start With Letter, End With Number
			Input: "A1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter
			Input: "1A",
			Valid: true,
		},

		{
			// Start With Letter, End With Letter and Hyphen Separator
			Input: "A-A",
			Valid: true,
		},

		{
			// Start With Number, End With Number and Hyphen Separator
			Input: "1-1",
			Valid: true,
		},

		{
			// Start With Letter, End With Number and Hyphen Separator
			Input: "A-1",
			Valid: true,
		},

		{
			// Start With Number, End With Letter and Hyphen Separator
			Input: "1-A",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorEndpointName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
