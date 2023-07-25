// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudConfigurationServiceId{}

func TestSpringCloudConfigurationServiceIDFormatter(t *testing.T) {
	actual := NewSpringCloudConfigurationServiceID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "service1", "configurationService1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/configurationService1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudConfigurationServiceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudConfigurationServiceId
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
			// missing ConfigurationServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for ConfigurationServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/configurationService1",
			Expected: &SpringCloudConfigurationServiceId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				SpringName:               "service1",
				ConfigurationServiceName: "configurationService1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SERVICE1/CONFIGURATIONSERVICES/CONFIGURATIONSERVICE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudConfigurationServiceID(v.Input)
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
		if actual.ConfigurationServiceName != v.Expected.ConfigurationServiceName {
			t.Fatalf("Expected %q but got %q for ConfigurationServiceName", v.Expected.ConfigurationServiceName, actual.ConfigurationServiceName)
		}
	}
}

func TestSpringCloudConfigurationServiceIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudConfigurationServiceId
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
			// missing ConfigurationServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for ConfigurationServiceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationServices/configurationService1",
			Expected: &SpringCloudConfigurationServiceId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				SpringName:               "service1",
				ConfigurationServiceName: "configurationService1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/configurationservices/configurationService1",
			Expected: &SpringCloudConfigurationServiceId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				SpringName:               "service1",
				ConfigurationServiceName: "configurationService1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SPRING/service1/CONFIGURATIONSERVICES/configurationService1",
			Expected: &SpringCloudConfigurationServiceId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				SpringName:               "service1",
				ConfigurationServiceName: "configurationService1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SpRiNg/service1/CoNfIgUrAtIoNsErViCeS/configurationService1",
			Expected: &SpringCloudConfigurationServiceId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				SpringName:               "service1",
				ConfigurationServiceName: "configurationService1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudConfigurationServiceIDInsensitively(v.Input)
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
		if actual.ConfigurationServiceName != v.Expected.ConfigurationServiceName {
			t.Fatalf("Expected %q but got %q for ConfigurationServiceName", v.Expected.ConfigurationServiceName, actual.ConfigurationServiceName)
		}
	}
}
