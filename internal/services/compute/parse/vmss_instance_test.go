// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = VMSSInstanceId{}

func TestVMSSInstanceIDFormatter(t *testing.T) {
	actual := NewVMSSInstanceID("12345678-1234-9876-4563-123456789012", "resGroup1", "vmss1", "vm1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/vm1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVMSSInstanceID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VMSSInstanceId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/",
			Error: true,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/vm1",
			Expected: &VMSSInstanceId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "vmss1",
				VirtualMachineName:         "vm1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINESCALESETS/VMSS1/VIRTUALMACHINES/VM1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VMSSInstanceID(v.Input)
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
	}
}

func TestVMSSInstanceIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VMSSInstanceId
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
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/",
			Error: true,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/virtualMachines/vm1",
			Expected: &VMSSInstanceId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "vmss1",
				VirtualMachineName:         "vm1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualmachinescalesets/vmss1/virtualmachines/vm1",
			Expected: &VMSSInstanceId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "vmss1",
				VirtualMachineName:         "vm1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/VIRTUALMACHINESCALESETS/vmss1/VIRTUALMACHINES/vm1",
			Expected: &VMSSInstanceId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "vmss1",
				VirtualMachineName:         "vm1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/ViRtUaLmAcHiNeScAlEsEtS/vmss1/ViRtUaLmAcHiNeS/vm1",
			Expected: &VMSSInstanceId{
				SubscriptionId:             "12345678-1234-9876-4563-123456789012",
				ResourceGroup:              "resGroup1",
				VirtualMachineScaleSetName: "vmss1",
				VirtualMachineName:         "vm1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VMSSInstanceIDInsensitively(v.Input)
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
	}
}
