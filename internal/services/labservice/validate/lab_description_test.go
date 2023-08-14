// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabDescription(t *testing.T) {
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
			Input:       strings.Repeat("s", 499),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 500),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 501),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabDescription(tc.Input, "description")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Description to trigger a validation error for '%s'", tc.Input)
		}
	}
}
