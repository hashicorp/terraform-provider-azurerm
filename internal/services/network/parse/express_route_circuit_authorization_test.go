// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = ExpressRouteCircuitAuthorizationId{}

func TestExpressRouteCircuitAuthorizationIDFormatter(t *testing.T) {
	actual := NewExpressRouteCircuitAuthorizationID("12345678-1234-9876-4563-123456789012", "resGroup1", "expressRouteCircuit1", "authorization1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/expressRouteCircuit1/authorizations/authorization1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestExpressRouteCircuitAuthorizationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ExpressRouteCircuitAuthorizationId
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
			// missing ExpressRouteCircuitName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for ExpressRouteCircuitName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/",
			Error: true,
		},

		{
			// missing AuthorizationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/expressRouteCircuit1/",
			Error: true,
		},

		{
			// missing value for AuthorizationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/expressRouteCircuit1/authorizations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/expressRouteCircuits/expressRouteCircuit1/authorizations/authorization1",
			Expected: &ExpressRouteCircuitAuthorizationId{
				SubscriptionId:          "12345678-1234-9876-4563-123456789012",
				ResourceGroup:           "resGroup1",
				ExpressRouteCircuitName: "expressRouteCircuit1",
				AuthorizationName:       "authorization1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/EXPRESSROUTECIRCUITS/EXPRESSROUTECIRCUIT1/AUTHORIZATIONS/AUTHORIZATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ExpressRouteCircuitAuthorizationID(v.Input)
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
		if actual.ExpressRouteCircuitName != v.Expected.ExpressRouteCircuitName {
			t.Fatalf("Expected %q but got %q for ExpressRouteCircuitName", v.Expected.ExpressRouteCircuitName, actual.ExpressRouteCircuitName)
		}
		if actual.AuthorizationName != v.Expected.AuthorizationName {
			t.Fatalf("Expected %q but got %q for AuthorizationName", v.Expected.AuthorizationName, actual.AuthorizationName)
		}
	}
}
