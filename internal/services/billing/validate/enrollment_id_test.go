package validate

import "testing"

func TestEnrollmentID(t *testing.T) {
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
			// missing EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/",
			Valid: false,
		},

		{
			// missing value for EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/enrollmentAccounts/",
			Valid: false,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/enrollmentAccounts/12345678-1234-9876-4563-123456789012",
			Valid: true,
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/ENROLLMENTACCOUNTS/12345678-1234-9876-4563-123456789012",
			Valid: false,
		},
	}
	for _, tc := range cases {
		t.Logf("[DEBUG] Testing Value %s", tc.Input)
		_, errors := EnrollmentID(tc.Input, "test")
		valid := len(errors) == 0

		if tc.Valid != valid {
			t.Fatalf("Expected %t but got %t", tc.Valid, valid)
		}
	}
}
