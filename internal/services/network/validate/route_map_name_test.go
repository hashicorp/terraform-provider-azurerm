// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestValidateRouteMapName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "he.l-l_o_",
			ExpectError: false,
		},
		{
			Input:       "8he.l-8l_o_8",
			ExpectError: false,
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
		_, errors := RouteMapName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Route Map Name to trigger a validation error for '%s'", tc.Input)
		}
	}
}
