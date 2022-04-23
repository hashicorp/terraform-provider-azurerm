package parse

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ConsumptionBudgetManagementGroupId{}

func TestConsumptionBudgetManagementGroupIDFormatter(t *testing.T) {
	actual := NewConsumptionBudgetManagementGroupID("12345678-1234-9876-4563-123456789012", "budget1").ID()
	expected := "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/budget1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestConsumptionBudgetManagementGroupID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ConsumptionBudgetManagementGroupId
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
			// missing BudgetName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/",
			Error: true,
		},

		{
			// missing value for BudgetName
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Management/managementGroups/12345678-1234-9876-4563-123456789012/providers/Microsoft.Consumption/budgets/budget1",
			Expected: &ConsumptionBudgetManagementGroupId{
				ManagementGroupName: "12345678-1234-9876-4563-123456789012",
				BudgetName:          "budget1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.MANAGEMENT/MANAGEMENTGROUPS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.CONSUMPTION/BUDGETS/BUDGET1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ConsumptionBudgetManagementGroupID(v.Input)
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
		if actual.BudgetName != v.Expected.BudgetName {
			t.Fatalf("Expected %q but got %q for BudgetName", v.Expected.BudgetName, actual.BudgetName)
		}
	}
}
