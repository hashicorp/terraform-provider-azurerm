package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = VirtualMachineExtensionId{}

func TestVirtualMachineExtensionIDFormatter(t *testing.T) {
	subscriptionId := "12345678-1234-5678-1234-123456789012"
	vmId := NewVirtualMachineId(subscriptionId, "group1", "vm1")
	actual := NewVirtualMachineExtensionId(vmId, "extension1").ID(subscriptionId)
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/group1/providers/Microsoft.Compute/virtualMachines/vm1/extensions/extension1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestParseVirtualMachineExtensionID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *VirtualMachineExtensionId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No virtual machine segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/",
			Expected: nil,
		},
		{
			Name:     "No extension name",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/extension/",
			Expected: nil,
		},
		{
			Name:     "Case incorrect in path element",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/microsoft.compute/virtualMachines/machine1/Extensions/extName",
			Expected: nil,
		},
		{
			Name:  "Valid",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/myGroup1/providers/Microsoft.Compute/virtualMachines/machine1/extensions/extName",
			Expected: &VirtualMachineExtensionId{
				SubscriptionId: "00000000-0000-0000-0000-000000000000",
				ResourceGroup:  "myGroup1",
				Name:           "extName",
				VirtualMachine: "machine1",
			},
		},
	}
	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := VirtualMachineExtensionID(v.Input)
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

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for Subscription Id", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
	}
}
