package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = EnrollmentBillingScopeId{}

func TestEnrollmentBillingScopeIDFormatter(t *testing.T) {
	actual := NewEnrollmentBillingScopeID("12345678", "123456").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/12345678/enrollmentAccounts/123456"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestEnrollmentBillingScopeID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EnrollmentBillingScopeId
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
			// missing EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/",
			Error: true,
		},

		{
			// missing value for EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/12345678/enrollmentAccounts/123456",
			Expected: &EnrollmentBillingScopeId{
				BillingAccountName:    "12345678",
				EnrollmentAccountName: "123456",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/12345678/ENROLLMENTACCOUNTS/123456",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := EnrollmentBillingScopeID(v.Input)
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
		if actual.EnrollmentAccountName != v.Expected.EnrollmentAccountName {
			t.Fatalf("Expected %q but got %q for EnrollmentAccountName", v.Expected.EnrollmentAccountName, actual.EnrollmentAccountName)
		}
	}
}
