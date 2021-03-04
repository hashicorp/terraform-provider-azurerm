package validate

import "testing"

func TestMicrosoftCustomerAccountBillingScopeID(t *testing.T) {
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
			// missing BillingProfileName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/",
			Valid: false,
		},

		{
			// missing value for BillingProfileName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/",
			Valid: false,
		},

		{
			// missing InvoiceSectionName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/",
			Valid: false,
		},

		{
			// missing value for InvoiceSectionName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/invoiceSections/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/invoiceSections/MTT4-OBS7-PJA-TGB",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/E879CF0F-2B4D-5431-109A-F72FC9868693:024CABF4-7321-4CF9-BE59-DF0C77CA51DE_2019-05-31/BILLINGPROFILES/PE2Q-NOIT-BG7-TGB/INVOICESECTIONS/MTT4-OBS7-PJA-TGB",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := MicrosoftCustomerAccountBillingScopeID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
