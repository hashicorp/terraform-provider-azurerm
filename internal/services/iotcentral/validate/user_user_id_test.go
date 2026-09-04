// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestUserUserID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "-invalid-start",
			Valid: false,
		},
		{
			Input: "@",
			Valid: false,
		},
		{
			Input: "#",
			Valid: false,
		},
		{
			Input: "$",
			Valid: false,
		},
		{
			Input: "...",
			Valid: false,
		},
		{
			Input: "valid-string1",
			Valid: true,
		},
		{
			Input: "validstring2",
			Valid: true,
		},
		{
			Input: "v",
			Valid: true,
		},
		{
			Input: "1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := UserUserID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t, for input %s", tc.Valid, valid, tc.Input)
		}
	}
}
