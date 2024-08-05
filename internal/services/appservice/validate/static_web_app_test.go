// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestStaticWebAppPassword(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{
		{
			Input: "",
		},
		{
			Input: "short",
		},
		{
			Input: "nocapitallettersornumbers",
		},
		{
			Input: "NOLOWERCASELETTERSNUMBERSORSYMBOLS",
		},
		{
			Input: "CapitalLettersButNoNumbersOrSymbols",
		},
		{
			Input: "CapitalLettersWith0123ButNoSymbols",
		},
		{
			Input: "SECUREPASSWORD&1",
		},
		{
			Input: "SecUrePa$$Word1",
			Valid: true,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := StaticWebAppPassword(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
