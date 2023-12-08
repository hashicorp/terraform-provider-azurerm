// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestPolicyGroupName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "a",
			ExpectError: false,
		},
		{
			Input:       "a_",
			ExpectError: false,
		},
		{
			Input:       "aaaa.-_a",
			ExpectError: false,
		},
		{
			Input:       "_",
			ExpectError: true,
		},
		{
			Input:       strings.Repeat("s", 79),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 80),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 81),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := PolicyGroupName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Policy Group Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
