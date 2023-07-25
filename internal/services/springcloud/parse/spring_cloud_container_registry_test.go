// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudContainerRegistryId{}

func TestSpringCloudContainerRegistryIDFormatter(t *testing.T) {
	actual := NewSpringCloudContainerRegistryID("12345678-1234-9876-4563-123456789012", "resGroup1", "spring1", "containerRegistry1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerRegistries/containerRegistry1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudContainerRegistryID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudContainerRegistryId
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
			// missing ContainerRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ContainerRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerRegistries/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerRegistries/containerRegistry1",
			Expected: &SpringCloudContainerRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SpringName:            "spring1",
				ContainerRegistryName: "containerRegistry1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SPRING1/CONTAINERREGISTRIES/CONTAINERREGISTRY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudContainerRegistryID(v.Input)
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
		if actual.ContainerRegistryName != v.Expected.ContainerRegistryName {
			t.Fatalf("Expected %q but got %q for ContainerRegistryName", v.Expected.ContainerRegistryName, actual.ContainerRegistryName)
		}
	}
}

func TestSpringCloudContainerRegistryIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudContainerRegistryId
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
			// missing ContainerRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/",
			Error: true,
		},

		{
			// missing value for ContainerRegistryName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerRegistries/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerRegistries/containerRegistry1",
			Expected: &SpringCloudContainerRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SpringName:            "spring1",
				ContainerRegistryName: "containerRegistry1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/spring/spring1/containerregistries/containerRegistry1",
			Expected: &SpringCloudContainerRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SpringName:            "spring1",
				ContainerRegistryName: "containerRegistry1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SPRING/spring1/CONTAINERREGISTRIES/containerRegistry1",
			Expected: &SpringCloudContainerRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SpringName:            "spring1",
				ContainerRegistryName: "containerRegistry1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AppPlatform/SpRiNg/spring1/CoNtAiNeRrEgIsTrIeS/containerRegistry1",
			Expected: &SpringCloudContainerRegistryId{
				SubscriptionId:        "12345678-1234-9876-4563-123456789012",
				ResourceGroup:         "resGroup1",
				SpringName:            "spring1",
				ContainerRegistryName: "containerRegistry1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudContainerRegistryIDInsensitively(v.Input)
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
		if actual.ContainerRegistryName != v.Expected.ContainerRegistryName {
			t.Fatalf("Expected %q but got %q for ContainerRegistryName", v.Expected.ContainerRegistryName, actual.ContainerRegistryName)
		}
	}
}
