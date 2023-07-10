// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestDomainServiceName(t *testing.T) {
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
			// one length subdomain of invalid char
			Input: "-.com",
			Valid: false,
		},

		{
			// two length subdomain of invalid char (prefix)
			Input: "-a.com",
			Valid: false,
		},

		{
			// two length subdomain of invalid char (suffix)
			Input: "a-.com",
			Valid: false,
		},

		{
			// subdomain longer than 30
			Input: strings.Repeat("a", 31) + ".com",
			Valid: false,
		},

		{
			// three length subdomain
			Input: "aaa.com",
			Valid: true,
		},

		{
			// three length subdomain with dash
			Input: "a-a.com",
			Valid: true,
		},

		{
			// wide length subdomain
			Input: strings.Repeat("a", 30) + ".com",
			Valid: true,
		},

		{
			// one length subdomain
			Input: "a.com",
			Valid: true,
		},

		{
			// two length subdomain
			Input: "aa.com",
			Valid: true,
		},

		{
			// three length subdomain
			Input: "aaa.com",
			Valid: true,
		},

		{
			// three length subdomain with dash
			Input: "a-a.com",
			Valid: true,
		},

		{
			// wide length subdomain
			Input: strings.Repeat("a", 29) + ".com",
			Valid: true,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := DomainServiceName(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
