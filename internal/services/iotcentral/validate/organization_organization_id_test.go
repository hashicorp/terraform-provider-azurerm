// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestOrganizationOrganizationID(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "-invalid-start",
			ExpectError: true,
		},
		{
			Input:       "invalid--hyphen",
			ExpectError: true,
		},
		{
			Input:       "1234567890123456789012345678901234567890123456789",
			ExpectError: true,
		},
		{
			Input:       "valid-string1",
			ExpectError: false,
		},
		{
			Input:       "validstring2",
			ExpectError: false,
		},
		{
			Input:       "v",
			ExpectError: false,
		},
		{
			Input:       "1",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := OrganizationID(tc.Input, "example")
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
