// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestHubRouteTableName(t *testing.T) {
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
			Input:       "hello-world",
			ExpectError: false,
		},
		{
			Input:       "test<",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := HubRouteTableName(tc.Input, "name")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Virtual Hub Route Table Name to trigger a validation error for '%s'", tc.Input)
		}

		if !tc.ExpectError && hasError {
			t.Fatalf("Encountered unexpected validation error for Virtual Hub Route Table Name '%s'", tc.Input)
		}
	}
}
