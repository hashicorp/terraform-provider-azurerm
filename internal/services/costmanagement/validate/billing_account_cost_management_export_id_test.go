// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import "testing"

func TestBillingAccountCostManagementExportID(t *testing.T) {
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
			// missing BillingAccountName
			Input: "/providers/Microsoft.Billing/",
			Valid: false,
		},

		{
			// missing value for BillingAccountName
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Valid: false,
		},

		{
			// missing ExportName
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/",
			Valid: false,
		},

		{
			// missing value for ExportName
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/exports/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/exports/export1",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/12345678/PROVIDERS/MICROSOFT.COSTMANAGEMENT/EXPORTS/EXPORT1",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := BillingAccountCostManagementExportID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
