// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = NetworkInterfaceIpConfigurationId{}

func TestNetworkInterfaceIpConfigurationIDFormatter(t *testing.T) {
	actual := NewNetworkInterfaceIpConfigurationID("12345678-1234-9876-4563-123456789012", "resGroup1", "networkInterface1", "config1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/ipConfigurations/config1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestNetworkInterfaceIpConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *NetworkInterfaceIpConfigurationId
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
			// missing NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/",
			Error: true,
		},

		{
			// missing value for NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/",
			Error: true,
		},

		{
			// missing IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/",
			Error: true,
		},

		{
			// missing value for IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/ipConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Network/networkInterfaces/networkInterface1/ipConfigurations/config1",
			Expected: &NetworkInterfaceIpConfigurationId{
				SubscriptionId:       "12345678-1234-9876-4563-123456789012",
				ResourceGroup:        "resGroup1",
				NetworkInterfaceName: "networkInterface1",
				IpConfigurationName:  "config1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.NETWORK/NETWORKINTERFACES/NETWORKINTERFACE1/IPCONFIGURATIONS/CONFIG1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := NetworkInterfaceIpConfigurationID(v.Input)
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
		if actual.NetworkInterfaceName != v.Expected.NetworkInterfaceName {
			t.Fatalf("Expected %q but got %q for NetworkInterfaceName", v.Expected.NetworkInterfaceName, actual.NetworkInterfaceName)
		}
		if actual.IpConfigurationName != v.Expected.IpConfigurationName {
			t.Fatalf("Expected %q but got %q for IpConfigurationName", v.Expected.IpConfigurationName, actual.IpConfigurationName)
		}
	}
}
