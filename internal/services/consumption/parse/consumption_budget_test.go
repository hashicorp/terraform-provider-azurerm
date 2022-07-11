package parse

import (
	"testing"
)

func TestCostManagementExportID(t *testing.T) {
	testData := []struct {
		Name     string
		Input    string
		Error    bool
		Expected *ConsumptionBudgetId
	}{
		{
			Name:  "empty",
			Input: "",
			Error: true,
		},
		{
			Name:  "resource group consumption budget",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Consumption/budgets/budget1",
			Expected: &ConsumptionBudgetId{
				Name:  "budget1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo",
			},
		},
		{
			Name:  "resource group consumption budget but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/foo/providers/Microsoft.Consumption/budgets/",
			Error: true,
		},
		{
			Name:  "subscription consumption budget",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/budget1",
			Expected: &ConsumptionBudgetId{
				Name:  "budget1",
				Scope: "/subscriptions/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "subscription consumption budget but no name",
			Input: "/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/",
			Error: true,
		},
		{
			Name:  "management group consumption budget",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/budget1",
			Expected: &ConsumptionBudgetId{
				Name:  "budget1",
				Scope: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000",
			},
		},
		{
			Name:  "management group consumption budget but no name",
			Input: "/providers/Microsoft.Management/managementGroups/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Name)

		actual, err := ConsumptionBudgetID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expected a value but got an error: %+v", err)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q", v.Expected.Name, actual.Name)
		}

		if v.Expected.Scope != actual.Scope {
			t.Fatalf("Expected %+v but got %+v", v.Expected.Scope, actual.Scope)
		}
	}
}
