// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"testing"
)

func TestValidateVaultName(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "hi",
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
			Input:       "hello-world-21",
			ExpectError: false,
		},
		{
			Input:       "hello_world_21",
			ExpectError: true,
		},
		{
			Input:       "Hello-World",
			ExpectError: false,
		},
		{
			Input:       "20202020",
			ExpectError: false,
		},
		{
			Input:       "ABC123!@£",
			ExpectError: true,
		},
		{
			Input:       "abcdefghijklmnopqrstuvwxyz",
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := VaultName(tc.Input, "")

		hasError := len(errors) > 0
		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the Key Vault Name to trigger a validation error for '%s'", tc.Input)
		}

		if !tc.ExpectError && hasError {
			t.Fatalf("Encountered unexpected validation error for Key Vault Name '%s'", tc.Input)
		}
	}
}
