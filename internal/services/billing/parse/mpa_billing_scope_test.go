// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package parse

import "testing"

func TestMicrosoftPartnerAccountBillingScopeIDFormatter(t *testing.T) {
	actual := NewMPABillingScopeID("e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31", "2281f543-7321-4cf9-1e23-edb4Oc31a31c").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/customers/2281f543-7321-4cf9-1e23-edb4Oc31a31c"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestMicrosoftPartnerAccountBillingScopeID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MicrosoftPartnerAccountBillingScopeId
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
			// missing CustomerName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/",
			Error: true,
		},

		{
			// missing value for CustomerName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/customers/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/customers/2281f543-7321-4cf9-1e23-edb4Oc31a31c",
			Expected: &MicrosoftPartnerAccountBillingScopeId{
				BillingAccountName: "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31",
				CustomerName:       "2281f543-7321-4cf9-1e23-edb4Oc31a31c",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/E879CF0F-2B4D-5431-109A-F72FC9868693:024CABF4-7321-4CF9-BE59-DF0C77CA51DE_2019-05-31/CUSTOMERS/2281F543-7321-4CF9-1E23-EDB4OC31A31C",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := MicrosoftPartnerAccountBillingScopeID(v.Input)
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
		if actual.CustomerName != v.Expected.CustomerName {
			t.Fatalf("Expected %q but got %q for CustomerName", v.Expected.CustomerName, actual.CustomerName)
		}
	}
}
