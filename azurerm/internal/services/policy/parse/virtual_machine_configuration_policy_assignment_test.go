package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VirtualMachineConfigurationPolicyAssignmentId{}

func TestVirtualMachineConfigurationPolicyAssignmentIDFormatter(t *testing.T) {
	actual := NewVirtualMachineConfigurationPolicyAssignmentID("12345678-1234-9876-4563-123456789012", "resGroup1", "vm1", "assignment1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestVirtualMachineConfigurationPolicyAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualMachineConfigurationPolicyAssignmentId
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
			// missing VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/",
			Error: true,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},

		{
			// missing GuestConfigurationAssignmentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/",
			Error: true,
		},

		{
			// missing value for GuestConfigurationAssignmentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1",
			Expected: &VirtualMachineConfigurationPolicyAssignmentId{
				SubscriptionId:                   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                    "resGroup1",
				VirtualMachineName:               "vm1",
				GuestConfigurationAssignmentName: "assignment1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINES/VM1/PROVIDERS/MICROSOFT.GUESTCONFIGURATION/GUESTCONFIGURATIONASSIGNMENTS/ASSIGNMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualMachineConfigurationPolicyAssignmentID(v.Input)
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
		if actual.VirtualMachineName != v.Expected.VirtualMachineName {
			t.Fatalf("Expected %q but got %q for VirtualMachineName", v.Expected.VirtualMachineName, actual.VirtualMachineName)
		}
		if actual.GuestConfigurationAssignmentName != v.Expected.GuestConfigurationAssignmentName {
			t.Fatalf("Expected %q but got %q for GuestConfigurationAssignmentName", v.Expected.GuestConfigurationAssignmentName, actual.GuestConfigurationAssignmentName)
		}
	}
}

func TestVirtualMachineConfigurationPolicyAssignmentIDInsensitively(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *VirtualMachineConfigurationPolicyAssignmentId
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
			// missing VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/",
			Error: true,
		},

		{
			// missing value for VirtualMachineName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},

		{
			// missing GuestConfigurationAssignmentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/",
			Error: true,
		},

		{
			// missing value for GuestConfigurationAssignmentName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/assignment1",
			Expected: &VirtualMachineConfigurationPolicyAssignmentId{
				SubscriptionId:                   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                    "resGroup1",
				VirtualMachineName:               "vm1",
				GuestConfigurationAssignmentName: "assignment1",
			},
		},

		{
			// lower-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/virtualmachines/vm1/providers/Microsoft.GuestConfiguration/guestconfigurationassignments/assignment1",
			Expected: &VirtualMachineConfigurationPolicyAssignmentId{
				SubscriptionId:                   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                    "resGroup1",
				VirtualMachineName:               "vm1",
				GuestConfigurationAssignmentName: "assignment1",
			},
		},

		{
			// upper-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/VIRTUALMACHINES/vm1/providers/Microsoft.GuestConfiguration/GUESTCONFIGURATIONASSIGNMENTS/assignment1",
			Expected: &VirtualMachineConfigurationPolicyAssignmentId{
				SubscriptionId:                   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                    "resGroup1",
				VirtualMachineName:               "vm1",
				GuestConfigurationAssignmentName: "assignment1",
			},
		},

		{
			// mixed-cased segment names
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Compute/ViRtUaLmAcHiNeS/vm1/providers/Microsoft.GuestConfiguration/GuEsTcOnFiGuRaTiOnAsSiGnMeNtS/assignment1",
			Expected: &VirtualMachineConfigurationPolicyAssignmentId{
				SubscriptionId:                   "12345678-1234-9876-4563-123456789012",
				ResourceGroup:                    "resGroup1",
				VirtualMachineName:               "vm1",
				GuestConfigurationAssignmentName: "assignment1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := VirtualMachineConfigurationPolicyAssignmentIDInsensitively(v.Input)
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
		if actual.VirtualMachineName != v.Expected.VirtualMachineName {
			t.Fatalf("Expected %q but got %q for VirtualMachineName", v.Expected.VirtualMachineName, actual.VirtualMachineName)
		}
		if actual.GuestConfigurationAssignmentName != v.Expected.GuestConfigurationAssignmentName {
			t.Fatalf("Expected %q but got %q for GuestConfigurationAssignmentName", v.Expected.GuestConfigurationAssignmentName, actual.GuestConfigurationAssignmentName)
		}
	}
}
