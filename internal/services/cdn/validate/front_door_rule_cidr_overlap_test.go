// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import "testing"

func TestFrontDoorRuleCidrOverlap(t *testing.T) {
	cases := []struct {
		Input      []interface{}
		HasOverlap bool
	}{
		{
			// Not expected type
			Input:      []interface{}{192},
			HasOverlap: true,
		},

		{
			// Not expected type
			Input:      []interface{}{"192.168.0.1/24", 192},
			HasOverlap: true,
		},

		{
			// Single IPv4 CIDR
			Input:      []interface{}{"192.168.0.1/24"},
			HasOverlap: false,
		},

		{
			// Different IPv4 CIDRs
			Input:      []interface{}{"192.168.0.1/24", "192.168.1.1/24"},
			HasOverlap: false,
		},

		{
			// Single IPv6 CIDR
			Input:      []interface{}{"::FFFF:C0A8:0001/24"},
			HasOverlap: false,
		},

		{
			// Different IPv6 CIDRs
			Input:      []interface{}{"::FFFF:C0A8:0001/24", "2001:db8::/24"},
			HasOverlap: false,
		},

		{
			// Duplicate IPv4 CIDR
			Input:      []interface{}{"192.168.0.1/24", "192.168.0.1/24"},
			HasOverlap: true,
		},

		{
			// IPv4 Overlap least specific first
			Input:      []interface{}{"192.168.0.1/24", "192.168.0.1/26"},
			HasOverlap: true,
		},

		{
			// IPv4 Overlap least specific last
			Input:      []interface{}{"192.168.0.1/26", "192.168.0.1/24"},
			HasOverlap: true,
		},

		{
			// IPv6 CIDR least specific first
			Input:      []interface{}{"2001:db8::/24", "2001:db8::/32"},
			HasOverlap: true,
		},

		{
			// IPv6 CIDR least specific last
			Input:      []interface{}{"2001:db8::/32", "2001:db8::/24"},
			HasOverlap: true,
		},

		{
			// IPv6 CIDR compressed vs shortened
			Input:      []interface{}{"::ffff:ac10:0/32", "0:0:0:0:0:ffff:ac10:0/24"},
			HasOverlap: true,
		},

		{
			// IPv6 CIDR compressed vs expanded
			Input:      []interface{}{"::ffff:ac10:0/32", "0000:0000:0000:0000:0000:ffff:ac10:0000/24"},
			HasOverlap: true,
		},

		{
			// IPv6 CIDR shortened vs expanded
			Input:      []interface{}{"0:0:0:0:0:ffff:ac10:0/32", "0000:0000:0000:0000:0000:ffff:ac10:0000/24"},
			HasOverlap: true,
		},

		{
			// IPv4 and IPv6 least specific first
			Input:      []interface{}{"192.168.0.1/24", "2001:db8::/24", "2001:db8::/32"},
			HasOverlap: true,
		},

		{
			// IPv4 and IPv6 least specific last
			Input:      []interface{}{"192.168.0.1/24", "2001:db8::/32", "2001:db8::/24"},
			HasOverlap: true,
		},

		{
			// IPv4 overlap
			Input:      []interface{}{"172.16.0.0/22", "172.16.2.0/23", "172.16.3.0/24"},
			HasOverlap: true,
		},

		{
			// IPv4 overlap
			Input:      []interface{}{"172.16.4.0/22", "172.16.7.0/24", "172.16.6.0/23"},
			HasOverlap: true,
		},

		{
			// IPv6 overlap shortened, compressed, expanded
			Input:      []interface{}{"0:0:0:0:0:ffff:ac10:0/22", "::ffff:ac10:200/23", "0000:0000:0000:0000:0000:ffff:ac10:0300/24"},
			HasOverlap: true,
		},

		{
			// IPv6 overlap, expanded, compressed, shortened
			Input:      []interface{}{"0000:0000:0000:0000:0000:ffff:ac10:0400/22", "::ffff:ac10:700/24", "0:0:0:0:0:ffff:ac10:0600/23"},
			HasOverlap: true,
		},
	}

	for _, tc := range cases {
		_, errors := FrontDoorRuleCidrOverlap(tc.Input, "match_values")
		HasOverlap := len(errors) > 0

		if tc.HasOverlap != HasOverlap {
			t.Fatalf("[DEBUG] Testing Value %s, Expected HasOverlap to be %t but got %t", tc.Input, tc.HasOverlap, HasOverlap)
		} else {
			if len(errors) > 0 {
				t.Logf("[DEBUG] Testing Value %s Error: %+v", tc.Input, errors)
			} else {
				t.Logf("[DEBUG] Testing Value %s", tc.Input)
			}
		}
	}
}
