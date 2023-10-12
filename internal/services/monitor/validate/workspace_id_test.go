// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestWorkspaceID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// empty
			Input: "",
			Valid: false,
		},
		{
			Input: "12345678-1234-9876-4563-123456789012|123456789012",
			Valid: false,
		},
		{
			Input: "123456789012|12345678-1234-9876-4563-123456789012",
			Valid: false,
		},
		{
			Input: "|12345678-1234-9876-4563-123456789012",
			Valid: false,
		},
		{
			Input: "12345678-1234-9876-4563-123456789012|",
			Valid: false,
		},
		{
			Input: "abcds",
			Valid: false,
		},
		{
			Input: "123456789012",
			Valid: false,
		},
		{
			Input: "12345678-1234-9876-4563-123456789012",
			Valid: false,
		},
		{
			Input: "12345678-1234-9876-4563-123456789012|12345678-dcba-dcba-dcba-098765432109",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := WorkspaceID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
