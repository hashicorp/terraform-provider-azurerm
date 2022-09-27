package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = AutomanageConfigurationProfileHCIAssignmentId{}

func TestAutomanageConfigurationProfileHCIAssignmentIDFormatter(t *testing.T) {
	actual := NewAutomanageConfigurationProfileHCIAssignmentID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "cluster1", "configurationProfileAssignment1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/cluster1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomanageConfigurationProfileHCIAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomanageConfigurationProfileHCIAssignmentId
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
			// missing clusters
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/",
			Error: true,
		},
		{
			// missing value for clusters
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/",
			Error: true,
		},
		{
			// missing configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/cluster1/providers/Microsoft.Automanage/",
			Error: true,
		},
		{
			// missing value for configurationProfileAssignments
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/cluster1/providers/Microsoft.Automanage/configurationProfileAssignments/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.AzureStackHci/clusters/cluster1/providers/Microsoft.Automanage/configurationProfileAssignments/configurationProfileAssignment1",
			Expected: &AutomanageConfigurationProfileHCIAssignmentId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				ClusterName:    "cluster1",
				Name:           "configurationProfileAssignment1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.AZURESTACKHCI/CLUSTERS/CLUSTER1/PROVIDERS/MICROSOFT.AUTOMANAGE/CONFIGURATIONPROFILEASSIGNMENTS/CONFIGURATIONPROFILEASSIGNMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomanageConfigurationProfileHCIAssignmentID(v.Input)
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

		if actual.ClusterName != v.Expected.ClusterName {
			t.Fatalf("Expected %q but got %q for ClusterName", v.Expected.ClusterName, actual.ClusterName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
