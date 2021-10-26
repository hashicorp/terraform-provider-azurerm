package validate

import "testing"

func TestBillingProfileID(t *testing.T) {
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
			// missing Name
			Input: "/providers/Microsoft.Billing/billingAccounts/123456/",
			Valid: false,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/123456/billingProfiles/123456",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/123456/BILLINGPROFILES/123456",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := BillingProfileID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
