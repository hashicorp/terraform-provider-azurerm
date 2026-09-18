// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validation

import (
	"testing"
)

func TestIsBase64String(t *testing.T) {
	cases := []struct {
		Input  string
		Errors int
	}{
		{
			Input:  "",
			Errors: 1,
		},
		{
			Input:  "aGVsbG8td29ybGQ=",
			Errors: 0,
		},
		{
			Input:  "hello-world",
			Errors: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Input, func(t *testing.T) {
			if _, errors := IsBase64String(tc.Input, "base64"); len(errors) != tc.Errors {
				t.Fatalf("Expected Base64 string to have %d not %d errors for %q: %v", tc.Errors, len(errors), tc.Input, errors)
			}
		})
	}
}
