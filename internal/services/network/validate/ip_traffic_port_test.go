// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestIpTrafficPort(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "*",
			ExpectError: false,
		},
		{
			Input:       "hello",
			ExpectError: true,
		},
		{
			Input:       "-1",
			ExpectError: true,
		},
		{
			Input:       "0",
			ExpectError: false,
		},
		{
			Input:       "1",
			ExpectError: false,
		},
		{
			Input:       "65535",
			ExpectError: false,
		},
		{
			Input:       "65536",
			ExpectError: true,
		},
		{
			Input:       "100-120",
			ExpectError: false,
		},
		{
			Input:       "1-1",
			ExpectError: true,
		},
		{
			Input:       "-1-88",
			ExpectError: true,
		},
		{
			Input:       "0-65535",
			ExpectError: false,
		},
		{
			Input:       "0-65537",
			ExpectError: true,
		},
		{
			Input:       "-65535",
			ExpectError: true,
		},
		{
			Input:       "4-2",
			ExpectError: true,
		},
		{
			Input:       "2--10",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := IpTrafficPort(tc.Input, "port")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the IP Traffic Port to trigger a validation error for '%s'", tc.Input)
		}
	}
}
