// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestStatusCodeRange(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "10",
		},
		{
			Input: "999",
		},
		{
			Input: "100-600",
		},
		{
			Input: "300-200",
		},
		{
			Input: "200",
			Valid: true,
		},
		{
			Input: "300-302",
			Valid: true,
		},
		{
			Input: "100-599",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StatusCodeRange(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
