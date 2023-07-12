// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateSupportInfoPhone(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "+1-555-555-5555-5555",
			ExpectError: true,
		},
		{
			Input:       "+1-555-555-5555",
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		_, errors := SupportInfoEmail(tc.Input, "phone")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Support Info Phone to trigger a validation error for '%s'", tc.Input)
		}
	}
}
