// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateLabSkuName(t *testing.T) {
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
			Input:       strings.Repeat("s", 199),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 200),
			ExpectError: false,
		},
		{
			Input:       strings.Repeat("s", 201),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := LabSkuName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Lab Sku Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
