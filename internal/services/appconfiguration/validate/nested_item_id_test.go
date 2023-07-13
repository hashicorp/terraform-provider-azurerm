// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestNestedItemId(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/?label=testLabel",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=testLabel",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/test+%2F123?label=test%2B%2F123",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?b=",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a&b=c",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a&label=b",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=a%2Cb%2Cc",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=%00",
			ExpectError: true,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/testKey?label=%2500",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/.appconfig.featureflag%2FtestKey?label=testLabel",
			ExpectError: false,
		},
		{
			Input:       "https://testappconf1.azconfig.io/kv/.appconfig.featureflag/testKey?label=testLabel",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		warnings, err := NestedItemId(tc.Input, "example")
		if err != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, err)
			}

			return
		}

		if tc.ExpectError && len(warnings) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(warnings) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
