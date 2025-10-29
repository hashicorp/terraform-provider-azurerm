// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestManagedRedisAccessPolicyName(t *testing.T) {
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
			// Valid default policy
			Input: "default",
			Valid: true,
		},
		{
			// Valid with hyphens
			Input: "my-policy-name",
			Valid: true,
		},
		{
			// Valid with spaces
			Input: "Data Contributor",
			Valid: true,
		},
		{
			// Valid with spaces
			Input: "Data Owner",
			Valid: true,
		},
		{
			// Valid alphanumeric
			Input: "policy123",
			Valid: true,
		},
		{
			// Invalid - empty
			Input: "",
			Valid: false,
		},
		{
			// Invalid - starts with hyphen
			Input: "-invalid",
			Valid: false,
		},
		{
			// Invalid - ends with hyphen
			Input: "invalid-",
			Valid: false,
		},
		{
			// Invalid - starts with space
			Input: " invalid policy",
			Valid: false,
		},
		{
			// Invalid - ends with space
			Input: "invalid policy ",
			Valid: false,
		},
		{
			// Invalid - contains special characters
			Input: "invalid@policy",
			Valid: false,
		},
		{
			// Invalid - consecutive hyphens
			Input: "invalid--policy",
			Valid: false,
		},
		{
			// Valid - max length (64 chars)
			Input: "a123456789012345678901234567890123456789012345678901234567890123",
			Valid: true,
		},
		{
			// Invalid - too long (65 chars)
			Input: "a1234567890123456789012345678901234567890123456789012345678901234",
			Valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagedRedisAccessPolicyName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
