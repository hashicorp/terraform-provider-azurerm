// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VirtualMachineScaleSetPublicIPAddressId{}

func TestVirtualMachineScaleSetPublicIPAddressIDFormatter(t *testing.T) {
	actual := NewVirtualMachineScaleSetPublicIPAddressID("12345678-1234-9876-4563-123456789012", "resGroup1", "scaleSet1", "virtualMachine1", "networkInterface1", "ipConfiguration1", "publicIpAddress1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/publicIpAddress1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualMachineScaleSetPublicIPAddressID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualMachineScaleSetPublicIPAddressId
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
			// missing VirtualMachineScaleSetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/",
			Error: true,
		},

		{
			// missing value for VirtualMachineScaleSetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/",
			Error: true,
		},

		{
			// missing VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/",
			Error: true,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/",
			Error: true,
		},

		{
			// missing NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/",
			Error: true,
		},

		{
			// missing value for NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/",
			Error: true,
		},

		{
			// missing IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/",
			Error: true,
		},

		{
			// missing value for IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/",
			Error: true,
		},

		{
			// missing PublicIPAddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/",
			Error: true,
		},

		{
			// missing value for PublicIPAddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/publicIpAddress1",
			Expected: &VirtualMachineScaleSetPublicIPAddressId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "scaleSet1",
				VirtualMachineName:         "virtualMachine1",
				NetworkInterfaceName:       "networkInterface1",
				IpConfigurationName:        "ipConfiguration1",
				PublicIPAddressName:        "publicIpAddress1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS/SCALESET1/VIRTUALMACHINES/VIRTUALMACHINE1/NETWORKINTERFACES/NETWORKINTERFACE1/IPCONFIGURATIONS/IPCONFIGURATION1/PUBLICIPADDRESSES/PUBLICIPADDRESS1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualMachineScaleSetPublicIPAddressID(v.Input)
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
		if actual.VirtualMachineScaleSetName != v.Expected.VirtualMachineScaleSetName {
			t.Fatalf("Expected %q but got %q for VirtualMachineScaleSetName", v.Expected.VirtualMachineScaleSetName, actual.VirtualMachineScaleSetName)
		}
		if actual.VirtualMachineName != v.Expected.VirtualMachineName {
			t.Fatalf("Expected %q but got %q for VirtualMachineName", v.Expected.VirtualMachineName, actual.VirtualMachineName)
		}
		if actual.NetworkInterfaceName != v.Expected.NetworkInterfaceName {
			t.Fatalf("Expected %q but got %q for NetworkInterfaceName", v.Expected.NetworkInterfaceName, actual.NetworkInterfaceName)
		}
		if actual.IpConfigurationName != v.Expected.IpConfigurationName {
			t.Fatalf("Expected %q but got %q for IpConfigurationName", v.Expected.IpConfigurationName, actual.IpConfigurationName)
		}
		if actual.PublicIPAddressName != v.Expected.PublicIPAddressName {
			t.Fatalf("Expected %q but got %q for PublicIPAddressName", v.Expected.PublicIPAddressName, actual.PublicIPAddressName)
		}
	}
}
