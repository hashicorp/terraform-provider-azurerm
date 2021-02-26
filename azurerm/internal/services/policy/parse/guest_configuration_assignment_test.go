package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = GuestConfigurationAssignmentId{}

func TestGuestConfigurationAssignmentIDFormatter(t *testing.T) {
	id := NewGuestConfigurationAssignmentID("12345678-1234-5678-1234-123456789012", "resourceGroup1", "vm1", "guestConfigurationAssignment1")
	actual := id.ID()
	expected := "/subscriptions/12345678-1234-5678-1234-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/guestConfigurationAssignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestGuestConfigurationAssignmentID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *GuestConfigurationAssignmentId
	}{
		{
			Name:     "Empty",
			Input:    "",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Segment",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000",
			Expected: nil,
		},
		{
			Name:     "No Resource Groups Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/",
			Expected: nil,
		},
		{
			Name:     "Resource Group ID",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/",
			Expected: nil,
		},
		{
			Name:     "Missing GuestConfigurationAssignment Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments",
			Expected: nil,
		},
		{
			Name:  "guestConfiguration Assignment ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/guestConfigurationAssignments/guestConfigurationAssignment1",
			Expected: &GuestConfigurationAssignmentId{
				ResourceGroup: "resourceGroup1",
				VMName:        "vm1",
				Name:          "guestConfigurationAssignment1",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.GuestConfiguration/GuestConfigurationAssignments/guestConfigurationAssignment1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q..", v.Name)

		actual, err := GuestConfigurationAssignmentID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.VMName != v.Expected.VMName {
			t.Fatalf("Expected %q but got %q for VMName", v.Expected.VMName, actual.VMName)
		}
	}
}
