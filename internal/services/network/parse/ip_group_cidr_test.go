// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = IpGroupCidrId{}

func TestIpGroupCidrIDFormatter(t *testing.T) {
	actual := NewIpGroupCidrID("4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3", "rg-core-hub-networking-ip-groups", "my-ips", "127.0.0.1_32").ID()
	expected := "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/127.0.0.1_32"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestIpGroupCidrID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *IpGroupCidrId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing ResourceGroup
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/",
			Error: true,
		},

		{
			// missing IpGroupName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for IpGroupName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/",
			Error: true,
		},

		{
			// missing CidrName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/",
			Error: true,
		},

		{
			// missing value for CidrName
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3/resourceGroups/rg-core-hub-networking-ip-groups/providers/Microsoft.Network/ipGroups/my-ips/cidrs/127.0.0.1_32",
			Expected: &IpGroupCidrId{
				SubscriptionId: "4fe7b06f-0eb5-489a-b053-f2afbdd8a4f3",
				ResourceGroup:  "rg-core-hub-networking-ip-groups",
				IpGroupName:    "my-ips",
				CidrName:       "127.0.0.1_32",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/4FE7B06F-0EB5-489A-B053-F2AFBDD8A4F3/RESOURCEGROUPS/RG-CORE-HUB-NETWORKING-IP-GROUPS/PROVIDERS/MICROSOFT.NETWORK/IPGROUPS/MY-IPS/CIDRS/127.0.0.1_32",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := IpGroupCidrID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
		if actual.IpGroupName != v.Expected.IpGroupName {
			t.Fatalf("Expected %q but got %q for IpGroupName", v.Expected.IpGroupName, actual.IpGroupName)
		}
		if actual.CidrName != v.Expected.CidrName {
			t.Fatalf("Expected %q but got %q for CidrName", v.Expected.CidrName, actual.CidrName)
		}
	}
}
