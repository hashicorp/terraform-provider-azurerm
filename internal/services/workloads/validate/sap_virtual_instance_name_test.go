// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestSAPVirtualInstanceName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "x01",
			Valid: false,
		},
		{
			Input: "xx",
			Valid: false,
		},
		{
			Input: "x0222",
			Valid: false,
		},
		{
			Input: "X01",
			Valid: true,
		},
		{
			Input: "XAA",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SAPVirtualInstanceName(tc.Input, "name")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
