package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ManagementGroupTemplateDeploymentId{}

func TestManagementGroupTemplateDeploymentIDFormatter(t *testing.T) {
	actual := NewManagementGroupTemplateDeploymentID("my-management-group-id", "deploy1").ID()
	expected := "/providers/Microsoft.Management/managementGroups/my-management-group-id/providers/Microsoft.Resources/deployments/deploy1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagementGroupTemplateDeploymentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagementGroupTemplateDeploymentId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing ManagementGroupName
			Input: "/providers/Microsoft.Management/",
			Error: true,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/Microsoft.Management/managementGroups/",
			Error: true,
		},

		{
			// missing DeploymentName
			Input: "/providers/Microsoft.Management/managementGroups/my-management-group-id/providers/Microsoft.Resources/",
			Error: true,
		},

		{
			// missing value for DeploymentName
			Input: "/providers/Microsoft.Management/managementGroups/my-management-group-id/providers/Microsoft.Resources/deployments/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/my-management-group-id/providers/Microsoft.Resources/deployments/deploy1",
			Expected: &ManagementGroupTemplateDeploymentId{
				ManagementGroupName: "my-management-group-id",
				DeploymentName:      "deploy1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/MY-MANAGEMENT-GROUP-ID/PROVIDERS/MICROSOFT.RESOURCES/DEPLOYMENTS/DEPLOY1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagementGroupTemplateDeploymentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.ManagementGroupName != v.Expected.ManagementGroupName {
			t.Fatalf("Expected %q but got %q for ManagementGroupName", v.Expected.ManagementGroupName, actual.ManagementGroupName)
		}
		if actual.DeploymentName != v.Expected.DeploymentName {
			t.Fatalf("Expected %q but got %q for DeploymentName", v.Expected.DeploymentName, actual.DeploymentName)
		}
	}
}
