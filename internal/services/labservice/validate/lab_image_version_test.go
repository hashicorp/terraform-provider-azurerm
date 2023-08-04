// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabImageVersion(t *testing.T) {
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
			Input:       strings.Repeat("s", 9),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 10),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 11),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabImageVersion(tc.Input, "version")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Image Version to trigger a validation error for '%s'", tc.Input)
		}
	}
}
