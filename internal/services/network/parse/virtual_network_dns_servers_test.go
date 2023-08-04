// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VirtualNetworkDnsServersId{}

func TestVirtualNetworkDnsServersIDFormatter(t *testing.T) {
	actual := NewVirtualNetworkDnsServersID("12345678-1234-9876-4563-123456789012", "resGroup1", "network1", "default").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualNetworkDnsServersID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualNetworkDnsServersId
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
			// missing VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/",
			Error: true,
		},

		{
			// missing DnsServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/",
			Error: true,
		},

		{
			// missing value for DnsServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default",
			Expected: &VirtualNetworkDnsServersId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				VirtualNetworkName: "network1",
				DnsServerName:      "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/VIRTUALNETWORKS/NETWORK1/DNSSERVERS/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualNetworkDnsServersID(v.Input)
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
		if actual.VirtualNetworkName != v.Expected.VirtualNetworkName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkName", v.Expected.VirtualNetworkName, actual.VirtualNetworkName)
		}
		if actual.DnsServerName != v.Expected.DnsServerName {
			t.Fatalf("Expected %q but got %q for DnsServerName", v.Expected.DnsServerName, actual.DnsServerName)
		}
	}
}

func TestVirtualNetworkDnsServersIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualNetworkDnsServersId
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
			// missing VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/",
			Error: true,
		},

		{
			// missing DnsServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/",
			Error: true,
		},

		{
			// missing value for DnsServerName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualNetworks/network1/dnsServers/default",
			Expected: &VirtualNetworkDnsServersId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				VirtualNetworkName: "network1",
				DnsServerName:      "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/virtualnetworks/network1/dnsservers/default",
			Expected: &VirtualNetworkDnsServersId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				VirtualNetworkName: "network1",
				DnsServerName:      "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/VIRTUALNETWORKS/network1/DNSSERVERS/default",
			Expected: &VirtualNetworkDnsServersId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				VirtualNetworkName: "network1",
				DnsServerName:      "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/ViRtUaLnEtWoRkS/network1/DnSsErVeRs/default",
			Expected: &VirtualNetworkDnsServersId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "resGroup1",
				VirtualNetworkName: "network1",
				DnsServerName:      "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualNetworkDnsServersIDInsensitively(v.Input)
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
		if actual.VirtualNetworkName != v.Expected.VirtualNetworkName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkName", v.Expected.VirtualNetworkName, actual.VirtualNetworkName)
		}
		if actual.DnsServerName != v.Expected.DnsServerName {
			t.Fatalf("Expected %q but got %q for DnsServerName", v.Expected.DnsServerName, actual.DnsServerName)
		}
	}
}
