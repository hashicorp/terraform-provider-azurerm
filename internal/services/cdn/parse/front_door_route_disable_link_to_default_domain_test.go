// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = FrontDoorRouteDisableLinkToDefaultDomainId{}

func TestFrontDoorRouteDisableLinkToDefaultDomainIDFormatter(t *testing.T) {
	actual := NewFrontDoorRouteDisableLinkToDefaultDomainID("12345678-1234-9876-4563-123456789012", "resGroup1", "profile1", "endpoint1", "route1", "disableLinkToDefaultDomain1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/disableLinkToDefaultDomain/disableLinkToDefaultDomain1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestFrontDoorRouteDisableLinkToDefaultDomainID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *FrontDoorRouteDisableLinkToDefaultDomainId
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
			// missing ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/",
			Error: true,
		},

		{
			// missing value for ProfileName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/",
			Error: true,
		},

		{
			// missing AfdEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/",
			Error: true,
		},

		{
			// missing value for AfdEndpointName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/",
			Error: true,
		},

		{
			// missing RouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/",
			Error: true,
		},

		{
			// missing value for RouteName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/",
			Error: true,
		},

		{
			// missing DisableLinkToDefaultDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/",
			Error: true,
		},

		{
			// missing value for DisableLinkToDefaultDomainName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/disableLinkToDefaultDomain/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Cdn/profiles/profile1/afdEndpoints/endpoint1/routes/route1/disableLinkToDefaultDomain/disableLinkToDefaultDomain1",
			Expected: &FrontDoorRouteDisableLinkToDefaultDomainId{
				SubscriptionId:                 "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                  "resGroup1",
				ProfileName:                    "profile1",
				AfdEndpointName:                "endpoint1",
				RouteName:                      "route1",
				DisableLinkToDefaultDomainName: "disableLinkToDefaultDomain1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.CDN/PROFILES/PROFILE1/AFDENDPOINTS/ENDPOINT1/ROUTES/ROUTE1/DISABLELINKTODEFAULTDOMAIN/DISABLELINKTODEFAULTDOMAIN1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := FrontDoorRouteDisableLinkToDefaultDomainID(v.Input)
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
		if actual.ProfileName != v.Expected.ProfileName {
			t.Fatalf("Expected %q but got %q for ProfileName", v.Expected.ProfileName, actual.ProfileName)
		}
		if actual.AfdEndpointName != v.Expected.AfdEndpointName {
			t.Fatalf("Expected %q but got %q for AfdEndpointName", v.Expected.AfdEndpointName, actual.AfdEndpointName)
		}
		if actual.RouteName != v.Expected.RouteName {
			t.Fatalf("Expected %q but got %q for RouteName", v.Expected.RouteName, actual.RouteName)
		}
		if actual.DisableLinkToDefaultDomainName != v.Expected.DisableLinkToDefaultDomainName {
			t.Fatalf("Expected %q but got %q for DisableLinkToDefaultDomainName", v.Expected.DisableLinkToDefaultDomainName, actual.DisableLinkToDefaultDomainName)
		}
	}
}
