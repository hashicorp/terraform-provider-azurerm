// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SlotVirtualNetworkSwiftConnectionId{}

func TestSlotVirtualNetworkSwiftConnectionIDFormatter(t *testing.T) {
	actual := NewSlotVirtualNetworkSwiftConnectionID("12345678-1234-9876-4563-123456789012", "resGroup1", "site1", "slot1", "virtualNetwork").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/config/virtualNetwork"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSlotVirtualNetworkSwiftConnectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SlotVirtualNetworkSwiftConnectionId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},

		{
			// missing SiteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/",
			Error: true,
		},

		{
			// missing value for SiteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/",
			Error: true,
		},

		{
			// missing SlotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/",
			Error: true,
		},

		{
			// missing value for SlotName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/",
			Error: true,
		},

		{
			// missing ConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/",
			Error: true,
		},

		{
			// missing value for ConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/config/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Web/sites/site1/slots/slot1/config/virtualNetwork",
			Expected: &SlotVirtualNetworkSwiftConnectionId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				SiteName:       "site1",
				SlotName:       "slot1",
				ConfigName:     "virtualNetwork",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.WEB/SITES/SITE1/SLOTS/SLOT1/CONFIG/VIRTUALNETWORK",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SlotVirtualNetworkSwiftConnectionID(v.Input)
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
		if actual.SiteName != v.Expected.SiteName {
			t.Fatalf("Expected %q but got %q for SiteName", v.Expected.SiteName, actual.SiteName)
		}
		if actual.SlotName != v.Expected.SlotName {
			t.Fatalf("Expected %q but got %q for SlotName", v.Expected.SlotName, actual.SlotName)
		}
		if actual.ConfigName != v.Expected.ConfigName {
			t.Fatalf("Expected %q but got %q for ConfigName", v.Expected.ConfigName, actual.ConfigName)
		}
	}
}
