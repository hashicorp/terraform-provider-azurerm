// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestVirtualMachineScaleSetPublicIPAddressID(t *testing.T) {
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
			// missing VirtualMachineScaleSetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/",
			Valid: false,
		},

		{
			// missing value for VirtualMachineScaleSetName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/",
			Valid: false,
		},

		{
			// missing VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/",
			Valid: false,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/",
			Valid: false,
		},

		{
			// missing NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/",
			Valid: false,
		},

		{
			// missing value for NetworkInterfaceName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/",
			Valid: false,
		},

		{
			// missing IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/",
			Valid: false,
		},

		{
			// missing value for IpConfigurationName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/",
			Valid: false,
		},

		{
			// missing PublicIPAddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/",
			Valid: false,
		},

		{
			// missing value for PublicIPAddressName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/",
			Valid: false,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/scaleSet1/virtualMachines/virtualMachine1/networkInterfaces/networkInterface1/ipConfigurations/ipConfiguration1/publicIPAddresses/publicIpAddress1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS/SCALESET1/VIRTUALMACHINES/VIRTUALMACHINE1/NETWORKINTERFACES/NETWORKINTERFACE1/IPCONFIGURATIONS/IPCONFIGURATION1/PUBLICIPADDRESSES/PUBLICIPADDRESS1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := VirtualMachineScaleSetPublicIPAddressID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
