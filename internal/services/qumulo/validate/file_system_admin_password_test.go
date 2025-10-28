// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestFileSystemName(t *testing.T) {
	cases := []struct {
		Value    string
		ErrCount int
	}{
		{
			Value:    "a",
			ErrCount: 1,
		},
		{
			Value:    "abB1091",
			ErrCount: 1,
		},
		{
			Value:    "abB10919",
			ErrCount: 0,
		},
		{
			Value:    "hello-world!",
			ErrCount: 1,
		},
		{
			Value:    "hello_world",
			ErrCount: 1,
		},
		{
			Value:    "-helloWorld",
			ErrCount: 0,
		},
		{
			Value:    "helloWorld-",
			ErrCount: 0,
		},
		{
			Value:    "hello@world",
			ErrCount: 1,
		},
		{
			Value:    "hello@World",
			ErrCount: 0,
		},
		{
			Value:    "1234567890123456",
			ErrCount: 1,
		},
		{
			Value:    "123456789012345-",
			ErrCount: 1,
		},
		{
			Value:    "123456789012345-A",
			ErrCount: 0,
		},
	}

	for _, tc := range cases {
		_, errors := ValidatePasswordComplexity(tc.Value, "adminPassword")
		if len(errors) != tc.ErrCount {
			t.Fatalf("Expected ValidatePasswordComplexity to have %d not %d errors for %q", tc.ErrCount, len(errors), tc.Value)
		}
	}
}
