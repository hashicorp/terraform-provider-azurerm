package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VirtualMachineScaleSetExtensionId{}

func TestVirtualMachineScaleSetExtensionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	vmssId := NewVirtualMachineScaleSetId(subscriptionId, "group1", "vmss1")
	actual := NewVirtualMachineScaleSetExtensionId(vmssId, "extension1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/virtualMachineScaleSets/vmss1/extensions/extension1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseVirtualMachineScaleSetExtensionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualMachineScaleSetExtensionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			Expected: nil,
		},
		{
			Name:     "No Virtual Machine Scale Set Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/",
			Expected: nil,
		},
		{
			Name:     "No Extensions Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1",
			Expected: nil,
		},
		{
			Name:     "No Extensions Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1/extensions/",
			Expected: nil,
		},
		{
			Name:  "Completed",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/virtualMachineScaleSets/machine1/extensions/extension1",
			Expected: &VirtualMachineScaleSetExtensionId{
				SubscriptionId:             "00000000-0000-0000-0000-000000000000",
				ResourceGroup:              "foo",
				VirtualMachineScaleSetName: "machine1",
				Name:                       "extension1",
			},
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VirtualMachineScaleSetExtensionID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
