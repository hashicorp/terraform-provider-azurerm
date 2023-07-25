// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateSupportInfoEmail(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "test.terraform.com",
			ExpectError: true,
		},
		{
			Input:       "test@terraform.com",
			ExpectError: false,
		},
	}

	for _, tc := range cases {
		_, errors := SupportInfoEmail(tc.Input, "email")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Support Info Email to trigger a validation error for '%s'", tc.Input)
		}
	}
}
