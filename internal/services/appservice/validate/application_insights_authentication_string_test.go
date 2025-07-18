// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestApplicationInsightsAuthenticationString(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "Notvalid",
		},
		{
			Input: "Notvalid",
		},
		{
			Input: "Authorization=NAD",
		},
		{
			Input: "Authorization=AAD",
			Valid: true,
		},
		{
			Input: "Authorization=AAD;",
		},
		{
			Input: "Authorization=AAD;Garbage=1234567890",
		},
		{
			Input: "Authorization=AAD;ClientId=ee470da-669-4810-acd4-af4a3410a203",
		},
		{
			Input: "Authorization=AAD;ClientId=eee470da-6969-4810-acd4-af4a3410a203",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ApplicationInsightsAuthenticationString(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
