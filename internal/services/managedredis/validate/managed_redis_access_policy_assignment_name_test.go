// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestManagedRedisAccessPolicyAssignmentName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			// Valid single character
			Input: "a",
			Valid: true,
		},
		{
			// Invalid test name from acceptance tests (contains hyphens)
			Input: "acctest-assignment-251029123644929571",
			Valid: false,
		},
		{
			// Invalid with hyphens
			Input: "my-assignment-name",
			Valid: false,
		},
		{
			// Valid alphanumeric
			Input: "assignment123",
			Valid: true,
		},
		{
			// Invalid - empty
			Input: "",
			Valid: false,
		},
		{
			// Invalid - contains special characters
			Input: "invalid@name",
			Valid: false,
		},
		{
			// Invalid - contains spaces
			Input: "invalid name",
			Valid: false,
		},
		{
			// Valid - max length (60 chars)
			Input: "a12345678901234567890123456789012345678901234567890123456789",
			Valid: true,
		},
		{
			// Invalid - too long (61 chars)
			Input: "a1234567890123456789012345678901234567890123456789012345678901",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagedRedisAccessPolicyAssignmentName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
