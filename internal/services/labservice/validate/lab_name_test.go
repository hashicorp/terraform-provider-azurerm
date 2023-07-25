// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateLabName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hello",
			ExpectError: false,
		},
		{
			Input:       "h.e-l_lo2",
			ExpectError: false,
		},
		{
			Input:       "h#e-l_lo2",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
