// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestResourceDeploymentScriptAzurePowerShellVersion(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			Input: "",
			Valid: false,
		},

		{
			Input: "2",
			Valid: false,
		},

		{
			Input: "1.1.1",
			Valid: false,
		},

		{
			Input: "3.",
			Valid: false,
		},

		{
			Input: "9.7",
			Valid: true,
		},

		{
			Input: "10.3",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ResourceDeploymentScriptAzurePowerShellVersion(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
