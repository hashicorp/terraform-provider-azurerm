package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ManagementGroupCostManagementExportId{}

func TestManagementGroupCostManagementExportIDFormatter(t *testing.T) {
	actual := NewManagementGroupCostManagementExportID("group1", "export1").ID()
	expected := "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/exports/export1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestManagementGroupCostManagementExportID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ManagementGroupCostManagementExportId
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
			// missing ExportName
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/",
			Error: true,
		},

		{
			// missing value for ExportName
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/group1/providers/Microsoft.CostManagement/exports/export1",
			Expected: &ManagementGroupCostManagementExportId{
				ManagementGroupName: "group1",
				ExportName:          "export1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/GROUP1/PROVIDERS/MICROSOFT.COSTMANAGEMENT/EXPORTS/EXPORT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ManagementGroupCostManagementExportID(v.Input)
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
		if actual.ExportName != v.Expected.ExportName {
			t.Fatalf("Expected %q but got %q for ExportName", v.Expected.ExportName, actual.ExportName)
		}
	}
}
