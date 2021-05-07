package parse

import (
	"testing"
)

func TestCostManagementExportId(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Expected *CostManagementExportId
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
			Name:     "Missing Cost Management Export Value",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.CostManagement/exports/",
			Expected: nil,
		},
		{
			Name:  "Cost Management Export ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1",
			},
		},
		{
			Name:  "Cost Management Export Subscription ID",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "Cost Management Export Management Group ID",
			Input: "/providers/Microsoft.Management/managementGroups/TestMG/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/providers/Microsoft.Management/managementGroups/TestMG",
			},
		},
		{
			Name:  "Cost Management Export Billing Account ID",
			Input: "/providers/Microsoft.Billing/billingAccounts/123456/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/providers/Microsoft.Billing/billingAccounts/123456",
			},
		},
		{
			Name:  "Cost Management Export Billing Department ID",
			Input: "/providers/Microsoft.Billing/billingAccounts/12/departments/1234/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/providers/Microsoft.Billing/billingAccounts/12/departments/1234",
			},
		},
		{
			Name:  "Cost Management Export Enrollment Account ID",
			Input: "/providers/Microsoft.Billing/billingAccounts/100/enrollmentAccounts/456/providers/Microsoft.CostManagement/exports/Export1",
			Expected: &CostManagementExportId{
				Name:       "Export1",
				ResourceId: "/providers/Microsoft.Billing/billingAccounts/100/enrollmentAccounts/456",
			},
		},
		{
			Name:     "Wrong Casing",
			Input:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.CostManagement/Exports/Service1",
			Expected: nil,
		},
		{
			Name:     "Wrong Casing",
			Input:    "/providers/Microsoft.Billing/billingAccounts/100/enrollmentaccounts/456/providers/Microsoft.CostManagement/exports/Export1",
			Expected: nil,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := CostManagementExportID(v.Input)
		if err != nil {
			if v.Expected == nil {
				continue
			}

			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceId != v.Expected.ResourceId {
			t.Fatalf("Expected %q but got %q for Resource Group", v.Expected.ResourceId, actual.ResourceId)
		}
	}
}
