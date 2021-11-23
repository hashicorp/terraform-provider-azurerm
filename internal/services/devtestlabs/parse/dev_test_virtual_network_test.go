package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = DevTestVirtualNetworkId{}

func TestDevTestVirtualNetworkIDFormatter(t *testing.T) {
	actual := NewDevTestVirtualNetworkID("12345678-1234-9876-4563-123456789012", "group1", "lab1", "vnet1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualNetworks/vnet1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestDevTestVirtualNetworkID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *DevTestVirtualNetworkId
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
			// missing LabName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/",
			Error: true,
		},

		{
			// missing value for LabName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/",
			Error: true,
		},

		{
			// missing VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/",
			Error: true,
		},

		{
			// missing value for VirtualNetworkName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualNetworks/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.DevTestLab/labs/lab1/virtualNetworks/vnet1",
			Expected: &DevTestVirtualNetworkId{
				SubscriptionId:     "12345678-1234-9876-4563-123456789012",
				ResourceGroup:      "group1",
				LabName:            "lab1",
				VirtualNetworkName: "vnet1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/GROUP1/PROVIDERS/MICROSOFT.DEVTESTLAB/LABS/LAB1/VIRTUALNETWORKS/VNET1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := DevTestVirtualNetworkID(v.Input)
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
		if actual.LabName != v.Expected.LabName {
			t.Fatalf("Expected %q but got %q for LabName", v.Expected.LabName, actual.LabName)
		}
		if actual.VirtualNetworkName != v.Expected.VirtualNetworkName {
			t.Fatalf("Expected %q but got %q for VirtualNetworkName", v.Expected.VirtualNetworkName, actual.VirtualNetworkName)
		}
	}
}
