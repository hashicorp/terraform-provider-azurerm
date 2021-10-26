package parse

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = BillingProfileId{}

func TestBillingProfileIDFormatter(t *testing.T) {
	actual := NewBillingProfileID("123456", "123456").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/123456/billingProfiles/123456"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBillingProfileID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BillingProfileId
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
			// missing Name
			Input: "/providers/Microsoft.Billing/billingAccounts/123456/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/123456/billingProfiles/123456",
			Expected: &BillingProfileId{
				BillingAccountName: "123456",
				Name:               "123456",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/123456/BILLINGPROFILES/123456",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BillingProfileID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
