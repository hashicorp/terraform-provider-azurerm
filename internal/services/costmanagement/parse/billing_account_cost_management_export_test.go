// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = BillingAccountCostManagementExportId{}

func TestBillingAccountCostManagementExportIDFormatter(t *testing.T) {
	actual := NewBillingAccountCostManagementExportID("12345678", "export1").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/exports/export1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBillingAccountCostManagementExportID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BillingAccountCostManagementExportId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing BillingAccountName
			Input: "/providers/Microsoft.Billing/",
			Error: true,
		},

		{
			// missing value for BillingAccountName
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Error: true,
		},

		{
			// missing ExportName
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/",
			Error: true,
		},

		{
			// missing value for ExportName
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/exports/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/providers/Microsoft.CostManagement/exports/export1",
			Expected: &BillingAccountCostManagementExportId{
				BillingAccountName: "12345678",
				ExportName:         "export1",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/12345678/PROVIDERS/MICROSOFT.COSTMANAGEMENT/EXPORTS/EXPORT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BillingAccountCostManagementExportID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.BillingAccountName != v.Expected.BillingAccountName {
			t.Fatalf("Expected %q but got %q for BillingAccountName", v.Expected.BillingAccountName, actual.BillingAccountName)
		}
		if actual.ExportName != v.Expected.ExportName {
			t.Fatalf("Expected %q but got %q for ExportName", v.Expected.ExportName, actual.ExportName)
		}
	}
}
