// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudDevToolPortalId{}

func TestSpringCloudDevToolPortalIDFormatter(t *testing.T) {
	actual := NewSpringCloudDevToolPortalID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "service1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/DevToolPortals/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudDevToolPortalID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudDevToolPortalId
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
			// missing SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/",
			Error: true,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/",
			Error: true,
		},

		{
			// missing DevToolPortalName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/",
			Error: true,
		},

		{
			// missing value for DevToolPortalName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/DevToolPortals/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/Spring/service1/DevToolPortals/default",
			Expected: &SpringCloudDevToolPortalId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resourceGroup1",
				SpringName:        "service1",
				DevToolPortalName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SERVICE1/DEVTOOLPORTALS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudDevToolPortalID(v.Input)
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
		if actual.SpringName != v.Expected.SpringName {
			t.Fatalf("Expected %q but got %q for SpringName", v.Expected.SpringName, actual.SpringName)
		}
		if actual.DevToolPortalName != v.Expected.DevToolPortalName {
			t.Fatalf("Expected %q but got %q for DevToolPortalName", v.Expected.DevToolPortalName, actual.DevToolPortalName)
		}
	}
}
