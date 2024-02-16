// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestUserType(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "invalid",
			ExpectError: true,
		},
		{
			Input:       "123",
			ExpectError: true,
		},
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "Group",
			ExpectError: false,
		},
		{
			Input:       "ServicePrincipal",
			ExpectError: false,
		},
		{
			Input:       "Email",
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		warnings, err := UserType(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
