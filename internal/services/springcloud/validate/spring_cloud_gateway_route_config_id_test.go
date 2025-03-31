// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestSpringCloudGatewayRouteConfigID(t *testing.T) {
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Valid: false,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Valid: false,
		},

		{
			// missing SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/",
			Valid: false,
		},

		{
			// missing value for SpringName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/",
			Valid: false,
		},

		{
			// missing GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/",
			Valid: false,
		},

		{
			// missing value for GatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/",
			Valid: false,
		},

		{
			// missing RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/",
			Valid: false,
		},

		{
			// missing value for RouteConfigName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AppPlatform/spring/service1/gateways/gateway1/routeConfigs/routeConfig1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.APPPLATFORM/SPRING/SERVICE1/GATEWAYS/GATEWAY1/ROUTECONFIGS/ROUTECONFIG1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := SpringCloudGatewayRouteConfigID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
