// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VpnGatewayNatRuleId{}

func TestVpnGatewayNatRuleIDFormatter(t *testing.T) {
	actual := NewVpnGatewayNatRuleID("12345678-1234-9876-4563-123456789012", "resGroup1", "vpnGateway1", "natRule1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/natRules/natRule1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVpnGatewayNatRuleID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VpnGatewayNatRuleId
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
			// missing VpnGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VpnGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/",
			Error: true,
		},

		{
			// missing NatRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/",
			Error: true,
		},

		{
			// missing value for NatRuleName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/natRules/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/vpnGateways/vpnGateway1/natRules/natRule1",
			Expected: &VpnGatewayNatRuleId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resGroup1",
				VpnGatewayName: "vpnGateway1",
				NatRuleName:    "natRule1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/VPNGATEWAYS/VPNGATEWAY1/NATRULES/NATRULE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VpnGatewayNatRuleID(v.Input)
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
		if actual.VpnGatewayName != v.Expected.VpnGatewayName {
			t.Fatalf("Expected %q but got %q for VpnGatewayName", v.Expected.VpnGatewayName, actual.VpnGatewayName)
		}
		if actual.NatRuleName != v.Expected.NatRuleName {
			t.Fatalf("Expected %q but got %q for NatRuleName", v.Expected.NatRuleName, actual.NatRuleName)
		}
	}
}
