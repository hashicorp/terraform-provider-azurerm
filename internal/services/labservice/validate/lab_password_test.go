// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabPassword(t *testing.T) {
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
			Input:       strings.Repeat("s", 122),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 123),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 124),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabPassword(tc.Input, "password")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Password to trigger a validation error for '%s'", tc.Input)
		}
	}
}
