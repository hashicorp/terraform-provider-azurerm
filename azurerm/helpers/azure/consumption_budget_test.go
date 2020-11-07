package azure

import (
	"reflect"
	"testing"
)

func TestParseAzureConsumptionBudgetID(t *testing.T) {
	testCases := []struct {
		id                          string
		expectedConsumptionBudgetID *ConsumptionBudgetID
		expectError                 bool
	}{
		{
			// Empty string
			"",
			nil,
			true,
		},
		{
			// Missing provider
			"/subscriptions/00000000-0000-0000-0000-000000000000",
			nil,
			true,
		},
		{
			// Missing scope
			"/providers/Microsoft.Consumption/budgets/my-budget",
			nil,
			true,
		},
		{
			// Missing budget name
			"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets",
			nil,
			true,
		},
		{
			// Subscription
			"/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/subscriptions/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
		{
			// Resource Group
			"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/my-resource-group",
			},
			false,
		},
		{
			// Management Group
			"/providers/Microsoft.Management/managementGroups/my-management-group/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Management/managementGroups/my-management-group",
			},
			false,
		},
		{
			// Billing Account
			"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
		{
			// Billing profile
			"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/billingProfiles/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/billingProfiles/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
		{
			// Billing account invoice section
			"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/invoiceSection/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/invoiceSection/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
		{
			// Billing account department
			"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/departments/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/departments/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
		{
			// Billing account enrollment account
			"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/enrollmentAccounts/00000000-0000-0000-0000-000000000000/providers/Microsoft.Consumption/budgets/my-budget",
			&ConsumptionBudgetID{
				"my-budget",
				"/providers/Microsoft.Billing/billingAccounts/00000000-0000-0000-0000-000000000000/enrollmentAccounts/00000000-0000-0000-0000-000000000000",
			},
			false,
		},
	}

	for _, test := range testCases {
		t.Logf("[DEBUG] Testing %q", test.id)
		parsed, err := ParseAzureConsumptionBudgetID(test.id)
		if test.expectError && err != nil {
			continue
		}
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(test.expectedConsumptionBudgetID, parsed) {
			t.Fatalf("Unexpected consumption budget ID:\nExpected: %+v\nGot:      %+v\n", test.expectedConsumptionBudgetID, parsed)
		}
	}
}
