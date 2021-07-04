package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ManagementGroupAssignmentId{}

func TestManagementGroupAssignmentIDFormatter(t *testing.T) {
	actual := NewManagementGroupAssignmentID("managementGroup1", "assignment1").ID()
	expected := "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/policyAssignments/assignment1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagementGroupAssignmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagementGroupAssignmentId
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
			// missing PolicyAssignmentName
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/",
			Error: true,
		},

		{
			// missing value for PolicyAssignmentName
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/policyAssignments/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Expected: &ManagementGroupAssignmentId{
				ManagementGroupName:  "managementGroup1",
				PolicyAssignmentName: "assignment1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/MANAGEMENTGROUP1/PROVIDERS/MICROSOFT.AUTHORIZATION/POLICYASSIGNMENTS/ASSIGNMENT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagementGroupAssignmentID(v.Input)
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
		if actual.PolicyAssignmentName != v.Expected.PolicyAssignmentName {
			t.Fatalf("Expected %q but got %q for PolicyAssignmentName", v.Expected.PolicyAssignmentName, actual.PolicyAssignmentName)
		}
	}
}
