package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestManagementGroupCostManagementExportID(t *testing.T) {
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
			// missing ExportName
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/",
			Valid: false,
		},

		{
			// missing value for ExportName
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/exports/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/exports/export1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1/PROVIDERS/MICROSOFT.COSTMANAGEMENT/EXPORTS/EXPORT1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := ManagementGroupCostManagementExportID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
