// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorRuleCidrIsValid(t *testing.T) {
	cases := []struct {
		Input interface{}
		Valid bool
	}{

		{
			// IPv6 IPv4 literal
			Input: "::FFFF:192.168.0.1/24",
			Valid: true,
		},

		{
			// IPv4 Basic
			Input: "192.168.0.1/24",
			Valid: true,
		},

		{
			// IPv6 basic
			Input: "::FFFF:C0A8:1/24",
			Valid: true,
		},

		{
			// IPv6 compressed
			Input: "::FFFF:C0A8:0001/24",
			Valid: true,
		},

		{
			// IPv6 shortened
			Input: "0:0:0:0:0:FFFF:C0A8:1/24",
			Valid: true,
		},

		{
			// IPv6 shortened lower case
			Input: "0:0:0:0:0:ffff:c0a8:1/24",
			Valid: true,
		},

		{
			// IPv6 expanded
			Input: "0000:0000:0000:0000:0000:FFFF:C0A8:1/24",
			Valid: true,
		},

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// not expected type
			Input: 192,
			Valid: false,
		},

		{
			// IPv4 with Port
			Input: "192.168.0.1:80/24",
			Valid: false,
		},

		{
			// IPV4 with missing Port with Port delimiter
			Input: "192.168.0.1:/24",
			Valid: false,
		},

		{
			// IPv4 with invalid port
			Input: "192.168.0.1:foo/24",
			Valid: false,
		},

		{
			// IPv6 basic too many compressed zero octets
			Input: ":::FFFF:C0A8:1/24",
			Valid: false,
		},

		{
			// IPv6 with zone info
			Input: "::FFFF:C0A8:1%1/24",
			Valid: false,
		},

		{
			// IPv6 with invalid zone info
			Input: "::FFFF:C0A8:foo%bar/24",
			Valid: false,
		},

		{
			// IPv6 invalid IPv4 literal
			Input: "::FFFF:192.168.0.256/24",
			Valid: false,
		},

		{
			// IPv6 with port info
			Input: "[::FFFF:C0A8:1]:80/24",
			Valid: false,
		},

		{
			// IPv6 missing port info
			Input: "[::FFFF:C0A8:1]:/24",
			Valid: false,
		},

		{
			// IPv6 with invalid port info
			Input: "[::FFFF:C0A8:1]:foo/24",
			Valid: false,
		},

		{
			// IPv6 with zone and port info
			Input: "[::FFFF:C0A8:1%1]:80/24",
			Valid: false,
		},
	}

	for _, tc := range cases {
		_, errors := FrontDoorRuleCidrIsValid(tc.Input, "match_values")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("[DEBUG] Testing Value %s, Expected %t but got %t, Error: %+v", tc.Input, tc.Valid, valid, errors)
		} else {
			if !valid {
				t.Logf("[DEBUG] Testing Value %s, Error: %+v", tc.Input, errors)
			} else {
				t.Logf("[DEBUG] Testing Value %s", tc.Input)
			}
		}
	}
}
