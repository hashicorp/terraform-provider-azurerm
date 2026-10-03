// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = StackHCIVirtualMachineId{}

func TestStackHCIVirtualMachineIDFormatter(t *testing.T) {
	actual := NewStackHCIVirtualMachineID("00000000-0000-0000-0000-000000000000", "resourceGroup1", "machine1", "default").ID()
	expected := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestStackHCIVirtualMachineID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StackHCIVirtualMachineId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},

		{
			// missing MachineName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/",
			Error: true,
		},

		{
			// missing value for MachineName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/",
			Error: true,
		},

		{
			// missing VirtualMachineInstanceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/",
			Error: true,
		},

		{
			// missing value for VirtualMachineInstanceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default",
			Expected: &StackHCIVirtualMachineId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "resourceGroup1",
				MachineName:                "machine1",
				VirtualMachineInstanceName: "default",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/00000000-0000-0000-0000-000000000000/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.HYBRIDCOMPUTE/MACHINES/MACHINE1/PROVIDERS/MICROSOFT.AZURESTACKHCI/VIRTUALMACHINEINSTANCES/DEFAULT",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StackHCIVirtualMachineID(v.Input)
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
		if actual.MachineName != v.Expected.MachineName {
			t.Fatalf("Expected %q but got %q for MachineName", v.Expected.MachineName, actual.MachineName)
		}
		if actual.VirtualMachineInstanceName != v.Expected.VirtualMachineInstanceName {
			t.Fatalf("Expected %q but got %q for VirtualMachineInstanceName", v.Expected.VirtualMachineInstanceName, actual.VirtualMachineInstanceName)
		}
	}
}

func TestStackHCIVirtualMachineIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *StackHCIVirtualMachineId
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
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/",
			Error: true,
		},

		{
			// missing value for ResourceGroup
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Error: true,
		},

		{
			// missing MachineName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/",
			Error: true,
		},

		{
			// missing value for MachineName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/",
			Error: true,
		},

		{
			// missing VirtualMachineInstanceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/",
			Error: true,
		},

		{
			// missing value for VirtualMachineInstanceName
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default",
			Expected: &StackHCIVirtualMachineId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "resourceGroup1",
				MachineName:                "machine1",
				VirtualMachineInstanceName: "default",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.AzureStackHCI/virtualmachineinstances/default",
			Expected: &StackHCIVirtualMachineId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "resourceGroup1",
				MachineName:                "machine1",
				VirtualMachineInstanceName: "default",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/MACHINES/machine1/providers/Microsoft.AzureStackHCI/VIRTUALMACHINEINSTANCES/default",
			Expected: &StackHCIVirtualMachineId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "resourceGroup1",
				MachineName:                "machine1",
				VirtualMachineInstanceName: "default",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/MaChInEs/machine1/providers/Microsoft.AzureStackHCI/ViRtUaLmAcHiNeInStAnCeS/default",
			Expected: &StackHCIVirtualMachineId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "resourceGroup1",
				MachineName:                "machine1",
				VirtualMachineInstanceName: "default",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := StackHCIVirtualMachineIDInsensitively(v.Input)
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
		if actual.MachineName != v.Expected.MachineName {
			t.Fatalf("Expected %q but got %q for MachineName", v.Expected.MachineName, actual.MachineName)
		}
		if actual.VirtualMachineInstanceName != v.Expected.VirtualMachineInstanceName {
			t.Fatalf("Expected %q but got %q for VirtualMachineInstanceName", v.Expected.VirtualMachineInstanceName, actual.VirtualMachineInstanceName)
		}
	}
}
