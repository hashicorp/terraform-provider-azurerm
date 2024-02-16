// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestUserType(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "invalid",
			Valid: false,
		},
		{
			Input: "123",
			Valid: false,
		},
		{
			Input: "",
			Valid: false,
		},
		{
			Input: "Group",
			Valid: true,
		},
		{
			Input: "ServicePrincipal",
			Valid: true,
		},
		{
			Input: "Email",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := UserType(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t, for input %s", tc.Valid, valid, tc.Input)
		}
	}
}
