// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabTitle(t *testing.T) {
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
			Input:       strings.Repeat("s", 99),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 100),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 101),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabTitle(tc.Input, "title")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Title to trigger a validation error for '%s'", tc.Input)
		}
	}
}
