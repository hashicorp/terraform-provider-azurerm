// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = SpringCloudGatewayRouteConfigId{}

func TestSpringCloudGatewayRouteConfigIDFormatter(t *testing.T) {
	actual := NewSpringCloudGatewayRouteConfigID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "service1", "gateway1", "routeConfig1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/routeConfig1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSpringCloudGatewayRouteConfigID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudGatewayRouteConfigId
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
			// missing GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/",
			Error: true,
		},

		{
			// missing RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/",
			Error: true,
		},

		{
			// missing value for RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/routeConfig1",
			Expected: &SpringCloudGatewayRouteConfigId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "resourceGroup1",
				SpringName:      "service1",
				GatewayName:     "gateway1",
				RouteConfigName: "routeConfig1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SERVICE1/GATEWAYS/GATEWAY1/ROUTECONFIGS/ROUTECONFIG1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudGatewayRouteConfigID(v.Input)
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
		if actual.GatewayName != v.Expected.GatewayName {
			t.Fatalf("Expected %q but got %q for GatewayName", v.Expected.GatewayName, actual.GatewayName)
		}
		if actual.RouteConfigName != v.Expected.RouteConfigName {
			t.Fatalf("Expected %q but got %q for RouteConfigName", v.Expected.RouteConfigName, actual.RouteConfigName)
		}
	}
}

func TestSpringCloudGatewayRouteConfigIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SpringCloudGatewayRouteConfigId
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
			// missing GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Error: true,
		},

		{
			// missing value for GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/",
			Error: true,
		},

		{
			// missing RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/",
			Error: true,
		},

		{
			// missing value for RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/routeConfig1",
			Expected: &SpringCloudGatewayRouteConfigId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "resourceGroup1",
				SpringName:      "service1",
				GatewayName:     "gateway1",
				RouteConfigName: "routeConfig1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeconfigs/routeConfig1",
			Expected: &SpringCloudGatewayRouteConfigId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "resourceGroup1",
				SpringName:      "service1",
				GatewayName:     "gateway1",
				RouteConfigName: "routeConfig1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SPRING/service1/GATEWAYS/gateway1/ROUTECONFIGS/routeConfig1",
			Expected: &SpringCloudGatewayRouteConfigId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "resourceGroup1",
				SpringName:      "service1",
				GatewayName:     "gateway1",
				RouteConfigName: "routeConfig1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/SpRiNg/service1/GaTeWaYs/gateway1/RoUtEcOnFiGs/routeConfig1",
			Expected: &SpringCloudGatewayRouteConfigId{
				SubscriptionId:  "12345678-1234-9876-4563-123456789012",
				ResourceGroup:   "resourceGroup1",
				SpringName:      "service1",
				GatewayName:     "gateway1",
				RouteConfigName: "routeConfig1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SpringCloudGatewayRouteConfigIDInsensitively(v.Input)
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
		if actual.GatewayName != v.Expected.GatewayName {
			t.Fatalf("Expected %q but got %q for GatewayName", v.Expected.GatewayName, actual.GatewayName)
		}
		if actual.RouteConfigName != v.Expected.RouteConfigName {
			t.Fatalf("Expected %q but got %q for RouteConfigName", v.Expected.RouteConfigName, actual.RouteConfigName)
		}
	}
}
