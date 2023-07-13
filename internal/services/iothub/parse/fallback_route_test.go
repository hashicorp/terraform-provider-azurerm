// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FallbackRouteId{}

func TestFallbackRouteIDFormatter(t *testing.T) {
	actual := NewFallbackRouteID("12345678-1234-9876-4563-123456789012", "resGroup1", "hub1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFallbackRouteID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FallbackRouteId
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
			// missing IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/",
			Error: true,
		},

		{
			// missing value for IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/",
			Error: true,
		},

		{
			// missing FallbackRouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/",
			Error: true,
		},

		{
			// missing value for FallbackRouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/default",
			Expected: &FallbackRouteId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				IotHubName:        "hub1",
				FallbackRouteName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.DEVICES/IOTHUBS/HUB1/FALLBACKROUTE/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FallbackRouteID(v.Input)
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
		if actual.IotHubName != v.Expected.IotHubName {
			t.Fatalf("Expected %q but got %q for IotHubName", v.Expected.IotHubName, actual.IotHubName)
		}
		if actual.FallbackRouteName != v.Expected.FallbackRouteName {
			t.Fatalf("Expected %q but got %q for FallbackRouteName", v.Expected.FallbackRouteName, actual.FallbackRouteName)
		}
	}
}

func TestFallbackRouteIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FallbackRouteId
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
			// missing IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/",
			Error: true,
		},

		{
			// missing value for IotHubName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/",
			Error: true,
		},

		{
			// missing FallbackRouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/",
			Error: true,
		},

		{
			// missing value for FallbackRouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iotHubs/hub1/fallbackRoute/default",
			Expected: &FallbackRouteId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				IotHubName:        "hub1",
				FallbackRouteName: "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/iothubs/hub1/fallbackroute/default",
			Expected: &FallbackRouteId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				IotHubName:        "hub1",
				FallbackRouteName: "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IOTHUBS/hub1/FALLBACKROUTE/default",
			Expected: &FallbackRouteId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				IotHubName:        "hub1",
				FallbackRouteName: "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Devices/IoThUbS/hub1/FaLlBaCkRoUtE/default",
			Expected: &FallbackRouteId{
				SubscriptionId:    "12345678-1234-9876-4563-123456789012",
				ResourceGroup:     "resGroup1",
				IotHubName:        "hub1",
				FallbackRouteName: "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FallbackRouteIDInsensitively(v.Input)
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
		if actual.IotHubName != v.Expected.IotHubName {
			t.Fatalf("Expected %q but got %q for IotHubName", v.Expected.IotHubName, actual.IotHubName)
		}
		if actual.FallbackRouteName != v.Expected.FallbackRouteName {
			t.Fatalf("Expected %q but got %q for FallbackRouteName", v.Expected.FallbackRouteName, actual.FallbackRouteName)
		}
	}
}
