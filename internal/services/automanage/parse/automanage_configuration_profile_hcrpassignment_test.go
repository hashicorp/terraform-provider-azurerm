package parse

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = AutomanageConfigurationProfileHCRPAssignmentId{}

func TestAutomanageConfigurationProfileHCRPAssignmentIDFormatter(t *testing.T) {
	actual := NewAutomanageConfigurationProfileHCRPAssignmentID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "machine1", "configurationProfileAssignment1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomanageConfigurationProfileHCRPAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomanageConfigurationProfileHCRPAssignmentId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// missing subscriptions
			Input: "/",
			Error: true,
		},
		{
			// missing value for subscriptions
			Input: "/subscriptions/",
			Error: true,
		},
		{
			// missing resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},
		{
			// missing value for resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},
		{
			// missing machines
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/",
			Error: true,
		},
		{
			// missing value for machines
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/",
			Error: true,
		},
		{
			// missing configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Automanage/",
			Error: true,
		},
		{
			// missing value for configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Automanage/configurationProfileAssignments/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.HybridCompute/machines/machine1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1",
			Expected: &AutomanageConfigurationProfileHCRPAssignmentId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				MachineName:    "machine1",
				Name:           "configurationProfileAssignment1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.HYBRIDCOMPUTE/MACHINES/MACHINE1/PROVIDERS/MICROSOFT.AUTOMANAGE/CONFIGURATIONPROFILEASSIGNMENTS/CONFIGURATIONPROFILEASSIGNMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomanageConfigurationProfileHCRPAssignmentID(v.Input)
		if err != nil {
			if v.Error {
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

		if actual.MachineName != v.Expected.MachineName {
			t.Fatalf("Expected %q but got %q for MachineName", v.Expected.MachineName, actual.MachineName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
