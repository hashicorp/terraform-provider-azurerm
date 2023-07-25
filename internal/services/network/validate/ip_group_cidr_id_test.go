// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestIpGroupCidrID(t *testing.T) {
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
			// missing SubscriptionId
			Input: "/",
			Valid: false,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Valid: false,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/",
			Valid: false,
		},

		{
			// missing IpGroupName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/",
			Valid: false,
		},

		{
			// missing value for IpGroupName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/",
			Valid: false,
		},

		{
			// missing CidrName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/",
			Valid: false,
		},

		{
			// missing value for CidrName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/127.0.0.1_32",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/4FE7B06F-0EB5-489A-B053-F2AFBDD8A4F3/RESOURCEGROUPS/RG-CORE-HUB-NETWORKING-IP-GROUPS/PROVIDERS/MICROSOFT.NETWORK/IPGROUPS/MY-IPS/CIDRS/127.0.0.1_32",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := IpGroupCidrID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
