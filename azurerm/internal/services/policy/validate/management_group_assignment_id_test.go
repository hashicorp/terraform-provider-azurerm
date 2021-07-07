package validate

import "testing"

func TestManagementGroupAssignmentID(t *testing.T) {
	cases := []struct {
		Input string
		Valid bool
	}{

		{
			// empty
			Input: "",
			Valid: false,
		},

		{
			// missing ManagementGroupName
			Input: "/providers/Microsoft.Management/",
			Valid: false,
		},

		{
			// missing value for ManagementGroupName
			Input: "/providers/Microsoft.Management/managementGroups/",
			Valid: false,
		},

		{
			// missing PolicyAssignmentName
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/",
			Valid: false,
		},

		{
			// missing value for PolicyAssignmentName
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/policyAssignments/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/managementGroup1/providers/Microsoft.Authorization/policyAssignments/assignment1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/MANAGEMENTGROUP1/PROVIDERS/MICROSOFT.AUTHORIZATION/POLICYASSIGNMENTS/ASSIGNMENT1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagementGroupAssignmentID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
