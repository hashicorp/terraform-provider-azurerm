// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabUsername(t *testing.T) {
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
			Input:       strings.Repeat("s", 19),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 20),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 21),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabUsername(tc.Input, "username")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Username to trigger a validation error for '%s'", tc.Input)
		}
	}
}
