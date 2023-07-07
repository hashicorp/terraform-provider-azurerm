// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestConfidentialLedgerName(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// too long
			Input: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Valid: false,
		},

		{
			// starts with -
			Input: "-my-ledger",
			Valid: false,
		},

		{
			// ends with -
			Input: "my-ledger-",
			Valid: false,
		},

		{
			// invalid characters
			Input: "my.ledger",
			Valid: false,
		},

		{
			// valid
			Input: "aA0-bB1-cC2-dD3-eE4-fF5-gG6-hH7i",
			Valid: true,
		},

		{
			// all lowercase characters valid
			Input: "abcdefghijklmnopqrstuvwxyz",
			Valid: true,
		},

		{
			// all numbers are valid
			Input: "0123456789",
			Valid: true,
		},

		{
			// all uppercase characters valid
			Input: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ConfidentialLedgerName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Logf("[DEBUG] Errors: %v", errors)
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
