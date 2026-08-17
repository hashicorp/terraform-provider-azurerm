// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"
)

func TestFrontDoorCustomDomainHostName(t *testing.T) {
	cases := []struct {
		input string
		valid bool
	}{
		{
			input: "",
			valid: false,
		},
		{
			input: "contoso.example.com",
			valid: true,
		},
		{
			input: "a.b",
			valid: true,
		},
		{
			input: "*.contoso.com",
			valid: true,
		},
		{
			input: "*foo",
			valid: false,
		},
		{
			input: "*foo.example.com",
			valid: false,
		},
		{
			input: "localhost",
			valid: false,
		},
		{
			input: "foo.*.contoso.com",
			valid: false,
		},
		{
			input: "-contoso.example.com",
			valid: false,
		},
		{
			input: "contoso-.example.com",
			valid: false,
		},
		{
			input: strings.Repeat("a", 64) + ".example.com",
			valid: false,
		},
		{
			input: strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 62),
			valid: false,
		},
	}

	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.input)
		_, errors := FrontDoorCustomDomainHostName(tc.input, "test")
		valid := len(errors) == 0

		if tc.valid != valid {
			t.Fatalf("expected %t but got %t", tc.valid, valid)
		}
	}
}
