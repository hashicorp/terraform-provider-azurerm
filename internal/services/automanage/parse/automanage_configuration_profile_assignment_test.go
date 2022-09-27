package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = AutomanageConfigurationProfileAssignmentId{}

func TestAutomanageConfigurationProfileAssignmentIDFormatter(t *testing.T) {
	actual := NewAutomanageConfigurationProfileAssignmentID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "vm1", "configurationProfileAssignment1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomanageConfigurationProfileAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomanageConfigurationProfileAssignmentId
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
			// missing virtualMachines
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/",
			Error: true,
		},
		{
			// missing value for virtualMachines
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/",
			Error: true,
		},
		{
			// missing configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Automanage/",
			Error: true,
		},
		{
			// missing value for configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Automanage/configurationProfileAssignments/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Compute/virtualMachines/vm1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1",
			Expected: &AutomanageConfigurationProfileAssignmentId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				VMName:         "vm1",
				Name:           "configurationProfileAssignment1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.COMPUTE/VIRTUALMACHINES/VM1/PROVIDERS/MICROSOFT.AUTOMANAGE/CONFIGURATIONPROFILEASSIGNMENTS/CONFIGURATIONPROFILEASSIGNMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomanageConfigurationProfileAssignmentID(v.Input)
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

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.VMName != v.Expected.VMName {
			t.Fatalf("Expected %q but got %q for VMName", v.Expected.VMName, actual.VMName)
		}
	}
}
