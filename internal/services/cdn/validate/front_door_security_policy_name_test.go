// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorSecurityPolicyName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "A",
			Valid: true,
		},
		{
			Input: "1",
			Valid: true,
		},
		{
			Input: "A-1",
			Valid: true,
		},
		{
			Input: "-A1",
			Valid: false,
		},
		{
			Input: "A1-",
			Valid: false,
		},
		{
			Input: "A_1",
			Valid: false,
		},
		{
			Input: "A 1",
			Valid: false,
		},
		{
			Input: "A!1",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := FrontDoorSecurityPolicyName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Testing value %q, Expected %t but got %t", tc.Input, tc.Valid, valid)
		}
	}
}
