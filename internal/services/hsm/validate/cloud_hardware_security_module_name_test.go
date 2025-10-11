// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateCloudHsmClusterName(t *testing.T) {
	testCases := []struct {
		Input    string
		Expected bool
	}{
		{
			// Valid: 3 characters, alphanumeric
			Input:    "abc",
			Expected: true,
		},
		{
			// Valid: 23 characters, alphanumeric
			Input:    strings.Repeat("a", 23),
			Expected: true,
		},
		{
			// Valid: with hyphens
			Input:    "hello-world",
			Expected: true,
		},
		{
			// Valid: starts with digit
			Input:    "1hello-world",
			Expected: true,
		},
		{
			// Valid: ends with digit
			Input:    "hello-world9",
			Expected: true,
		},
		{
			// Valid: contains digits and hyphens
			Input:    "test1-cluster2",
			Expected: true,
		},
		{
			// Invalid: too short (2 characters)
			Input:    "ab",
			Expected: false,
		},
		{
			// Invalid: too long (24 characters)
			Input:    strings.Repeat("a", 24),
			Expected: false,
		},
		{
			// Invalid: starts with hyphen
			Input:    "-hello-world",
			Expected: false,
		},
		{
			// Invalid: ends with hyphen
			Input:    "hello-world-",
			Expected: false,
		},
		{
			// Invalid: consecutive hyphens
			Input:    "hello--world",
			Expected: false,
		},
		{
			// Invalid: multiple consecutive hyphens
			Input:    "hello---world",
			Expected: false,
		},
		{
			// Invalid: contains special characters
			Input:    "hello_world",
			Expected: false,
		},
		{
			// Invalid: contains spaces
			Input:    "hello world",
			Expected: false,
		},
		{
			// Invalid: contains uppercase (should be case-insensitive, but let's test)
			Input:    "Hello-World",
			Expected: true,
		},
		{
			// Invalid: only hyphens
			Input:    "---",
			Expected: false,
		},
		{
			// Valid: minimum length with hyphen
			Input:    "a-b",
			Expected: true,
		},
		{
			// Valid: maximum length with hyphens
			Input:    "a" + strings.Repeat("-b", 10) + "c",
			Expected: true,
		},
	}

	for _, v := range testCases {
		_, errors := ValidateCloudHsmClusterName(v.Input, "name")
		result := len(errors) == 0
		if result != v.Expected {
			t.Fatalf("Expected the result to be %t for input %q but got %t (and %d errors)", v.Expected, v.Input, result, len(errors))
		}
	}
}
