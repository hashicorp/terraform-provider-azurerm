// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudBuildPackBindingId{}

func TestSpringCloudBuildPackBindingIDFormatter(t *testing.T) {
	actual := NewSpringCloudBuildPackBindingID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "service1", "buildService1", "builder1", "buildPackBinding1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/buildPackBinding1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudBuildPackBindingID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudBuildPackBindingId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing BuildServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for BuildServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/",
			Error: true,
		},

		{
			// missing BuilderName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/",
			Error: true,
		},

		{
			// missing value for BuilderName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/",
			Error: true,
		},

		{
			// missing BuildPackBindingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/",
			Error: true,
		},

		{
			// missing value for BuildPackBindingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/buildPackBinding1",
			Expected: &SpringCloudBuildPackBindingId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resourceGroup1",
				SpringName:           "service1",
				BuildServiceName:     "buildService1",
				BuilderName:          "builder1",
				BuildPackBindingName: "buildPackBinding1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SERVICE1/BUILDSERVICES/BUILDSERVICE1/BUILDERS/BUILDER1/BUILDPACKBINDINGS/BUILDPACKBINDING1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudBuildPackBindingID(v.Input)
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
		if actual.BuildServiceName != v.Expected.BuildServiceName {
			t.Fatalf("Expected %q but got %q for BuildServiceName", v.Expected.BuildServiceName, actual.BuildServiceName)
		}
		if actual.BuilderName != v.Expected.BuilderName {
			t.Fatalf("Expected %q but got %q for BuilderName", v.Expected.BuilderName, actual.BuilderName)
		}
		if actual.BuildPackBindingName != v.Expected.BuildPackBindingName {
			t.Fatalf("Expected %q but got %q for BuildPackBindingName", v.Expected.BuildPackBindingName, actual.BuildPackBindingName)
		}
	}
}

func TestSpringCloudBuildPackBindingIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudBuildPackBindingId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/",
			Error: true,
		},

		{
			// missing BuildServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for BuildServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/",
			Error: true,
		},

		{
			// missing BuilderName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/",
			Error: true,
		},

		{
			// missing value for BuilderName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/",
			Error: true,
		},

		{
			// missing BuildPackBindingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/",
			Error: true,
		},

		{
			// missing value for BuildPackBindingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildServices/buildService1/builders/builder1/buildPackBindings/buildPackBinding1",
			Expected: &SpringCloudBuildPackBindingId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resourceGroup1",
				SpringName:           "service1",
				BuildServiceName:     "buildService1",
				BuilderName:          "builder1",
				BuildPackBindingName: "buildPackBinding1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/buildservices/buildService1/builders/builder1/buildpackbindings/buildPackBinding1",
			Expected: &SpringCloudBuildPackBindingId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resourceGroup1",
				SpringName:           "service1",
				BuildServiceName:     "buildService1",
				BuilderName:          "builder1",
				BuildPackBindingName: "buildPackBinding1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SPRING/service1/BUILDSERVICES/buildService1/BUILDERS/builder1/BUILDPACKBINDINGS/buildPackBinding1",
			Expected: &SpringCloudBuildPackBindingId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resourceGroup1",
				SpringName:           "service1",
				BuildServiceName:     "buildService1",
				BuilderName:          "builder1",
				BuildPackBindingName: "buildPackBinding1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SpRiNg/service1/BuIlDsErViCeS/buildService1/BuIlDeRs/builder1/BuIlDpAcKbInDiNgS/buildPackBinding1",
			Expected: &SpringCloudBuildPackBindingId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resourceGroup1",
				SpringName:           "service1",
				BuildServiceName:     "buildService1",
				BuilderName:          "builder1",
				BuildPackBindingName: "buildPackBinding1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudBuildPackBindingIDInsensitively(v.Input)
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
		if actual.BuildServiceName != v.Expected.BuildServiceName {
			t.Fatalf("Expected %q but got %q for BuildServiceName", v.Expected.BuildServiceName, actual.BuildServiceName)
		}
		if actual.BuilderName != v.Expected.BuilderName {
			t.Fatalf("Expected %q but got %q for BuilderName", v.Expected.BuilderName, actual.BuilderName)
		}
		if actual.BuildPackBindingName != v.Expected.BuildPackBindingName {
			t.Fatalf("Expected %q but got %q for BuildPackBindingName", v.Expected.BuildPackBindingName, actual.BuildPackBindingName)
		}
	}
}
