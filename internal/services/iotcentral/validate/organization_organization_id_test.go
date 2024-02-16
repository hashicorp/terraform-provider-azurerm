// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestOrganizationOrganizationID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "-invalid-start",
			Valid: false,
		},
		{
			Input: "12345678901234567890123456789012345678901234567891",
			Valid: false,
		},
		{
			Input: "1234567890123456789012345678901234567890123456789",
			Valid: true,
		},
		{
			Input: "@",
			Valid: false,
		},
		{
			Input: "#!",
			Valid: false,
		},
		{
			Input: "$$",
			Valid: false,
		},
		{
			Input: "...",
			Valid: false,
		},
		{
			Input: "v",
			Valid: false,
		},
		{
			Input: "1",
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
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := OrganizationOrganizationID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t, for input %s", tc.Valid, valid, tc.Input)
		}
	}
}
