// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VirtualNetworkGatewayIpConfigurationId{}

func TestVirtualNetworkGatewayIpConfigurationIDFormatter(t *testing.T) {
	actual := NewVirtualNetworkGatewayIpConfigurationID("12345678-1234-9876-4563-123456789012", "resGroup1", "gw1", "cfg1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualNetworkGatewayIpConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualNetworkGatewayIpConfigurationId
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
			// missing VirtualNetworkGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/",
			Error: true,
		},

		{
			// missing IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/",
			Error: true,
		},

		{
			// missing value for IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1",
			Expected: &VirtualNetworkGatewayIpConfigurationId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				VirtualNetworkGatewayName: "gw1",
				IpConfigurationName:       "cfg1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/VIRTUALNETWORKGATEWAYS/GW1/IPCONFIGURATIONS/CFG1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualNetworkGatewayIpConfigurationID(v.Input)
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
		if actual.VirtualNetworkGatewayName != v.Expected.VirtualNetworkGatewayName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkGatewayName", v.Expected.VirtualNetworkGatewayName, actual.VirtualNetworkGatewayName)
		}
		if actual.IpConfigurationName != v.Expected.IpConfigurationName {
			t.Fatalf("Expected %q but got %q for IpConfigurationName", v.Expected.IpConfigurationName, actual.IpConfigurationName)
		}
	}
}

func TestVirtualNetworkGatewayIpConfigurationIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualNetworkGatewayIpConfigurationId
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
			// missing VirtualNetworkGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkGatewayName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/",
			Error: true,
		},

		{
			// missing IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/",
			Error: true,
		},

		{
			// missing value for IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworkGateways/gw1/ipConfigurations/cfg1",
			Expected: &VirtualNetworkGatewayIpConfigurationId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				VirtualNetworkGatewayName: "gw1",
				IpConfigurationName:       "cfg1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualnetworkgateways/gw1/ipconfigurations/cfg1",
			Expected: &VirtualNetworkGatewayIpConfigurationId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				VirtualNetworkGatewayName: "gw1",
				IpConfigurationName:       "cfg1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/VIRTUALNETWORKGATEWAYS/gw1/IPCONFIGURATIONS/cfg1",
			Expected: &VirtualNetworkGatewayIpConfigurationId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				VirtualNetworkGatewayName: "gw1",
				IpConfigurationName:       "cfg1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ViRtUaLnEtWoRkGaTeWaYs/gw1/IpCoNfIgUrAtIoNs/cfg1",
			Expected: &VirtualNetworkGatewayIpConfigurationId{
				SubscriptionId:            "12345678-1234-9876-4563-123456789012",
				ResourceGroup:             "resGroup1",
				VirtualNetworkGatewayName: "gw1",
				IpConfigurationName:       "cfg1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualNetworkGatewayIpConfigurationIDInsensitively(v.Input)
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
		if actual.VirtualNetworkGatewayName != v.Expected.VirtualNetworkGatewayName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkGatewayName", v.Expected.VirtualNetworkGatewayName, actual.VirtualNetworkGatewayName)
		}
		if actual.IpConfigurationName != v.Expected.IpConfigurationName {
			t.Fatalf("Expected %q but got %q for IpConfigurationName", v.Expected.IpConfigurationName, actual.IpConfigurationName)
		}
	}
}
