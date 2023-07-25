// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudServiceRegistryId{}

func TestSpringCloudServiceRegistryIDFormatter(t *testing.T) {
	actual := NewSpringCloudServiceRegistryID("12345678-1234-9876-4563-123456789012", "resGroup1", "spring1", "serviceRegistry1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceRegistries/serviceRegistry1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudServiceRegistryID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudServiceRegistryId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/",
			Error: true,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing ServiceRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ServiceRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceRegistries/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceRegistries/serviceRegistry1",
			Expected: &SpringCloudServiceRegistryId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "resGroup1",
				SpringName:          "spring1",
				ServiceRegistryName: "serviceRegistry1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SPRING1/SERVICEREGISTRIES/SERVICEREGISTRY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudServiceRegistryID(v.Input)
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
		if actual.ServiceRegistryName != v.Expected.ServiceRegistryName {
			t.Fatalf("Expected %q but got %q for ServiceRegistryName", v.Expected.ServiceRegistryName, actual.ServiceRegistryName)
		}
	}
}

func TestSpringCloudServiceRegistryIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudServiceRegistryId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/",
			Error: true,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing ServiceRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ServiceRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceRegistries/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceRegistries/serviceRegistry1",
			Expected: &SpringCloudServiceRegistryId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "resGroup1",
				SpringName:          "spring1",
				ServiceRegistryName: "serviceRegistry1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/serviceregistries/serviceRegistry1",
			Expected: &SpringCloudServiceRegistryId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "resGroup1",
				SpringName:          "spring1",
				ServiceRegistryName: "serviceRegistry1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SPRING/spring1/SERVICEREGISTRIES/serviceRegistry1",
			Expected: &SpringCloudServiceRegistryId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "resGroup1",
				SpringName:          "spring1",
				ServiceRegistryName: "serviceRegistry1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SpRiNg/spring1/SeRvIcErEgIsTrIeS/serviceRegistry1",
			Expected: &SpringCloudServiceRegistryId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				ResourceGroup:       "resGroup1",
				SpringName:          "spring1",
				ServiceRegistryName: "serviceRegistry1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudServiceRegistryIDInsensitively(v.Input)
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
		if actual.ServiceRegistryName != v.Expected.ServiceRegistryName {
			t.Fatalf("Expected %q but got %q for ServiceRegistryName", v.Expected.ServiceRegistryName, actual.ServiceRegistryName)
		}
	}
}
