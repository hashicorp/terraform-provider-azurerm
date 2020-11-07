package azure

import (
	"reflect"
	"testing"
	"time"
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

func TestValidateAzureConsumptionBudgetTimePeriodStartDate(t *testing.T) {
	// Set up time for testing
	now := time.Now()
	validTime := time.Date(
		now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		Input         string
		ExpectError   bool
		ExpectWarning bool
	}{
		{
			Input:         "",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			Input:         "2006-01-02",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Not on the first of a month
			Input:         "2020-11-02T00:00:00Z",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Before June 1, 2017
			Input:         "2000-01-01T00:00:00Z",
			ExpectError:   true,
			ExpectWarning: false,
		},
		{
			// Valid date and time
			Input:         validTime.Format(time.RFC3339),
			ExpectError:   false,
			ExpectWarning: false,
		},
		{
			// More than 12 months in the future
			Input:         validTime.AddDate(2, 0, 0).Format(time.RFC3339),
			ExpectError:   false,
			ExpectWarning: true,
		},
	}

	for _, tc := range cases {
		warnings, errors := ValidateAzureConsumptionBudgetTimePeriodStartDate(tc.Input, "start_date")
		if errors != nil {
			if !tc.ExpectError {
				t.Fatalf("Got error for input %q: %+v", tc.Input, errors)
			}

			return
		}

		if warnings != nil {
			if !tc.ExpectWarning {
				t.Fatalf("Got warnings for input %q: %+v", tc.Input, warnings)
			}

			return
		}

		if tc.ExpectError && len(errors) == 0 {
			t.Fatalf("Got no errors for input %q but expected some", tc.Input)
		} else if !tc.ExpectError && len(errors) > 0 {
			t.Fatalf("Got %d errors for input %q when didn't expect any", len(errors), tc.Input)
		}

		if tc.ExpectWarning && len(warnings) == 0 {
			t.Fatalf("Got no warnings for input %q but expected some", tc.Input)
		} else if !tc.ExpectWarning && len(warnings) > 0 {
			t.Fatalf("Got %d warnings for input %q when didn't expect any", len(warnings), tc.Input)
		}
	}
}
