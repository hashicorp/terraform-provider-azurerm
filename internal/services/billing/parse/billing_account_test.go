package parse

import (
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
)

var _ resourceid.Formatter = BillingAccountId{}

func TestBillingAccountIDFormatter(t *testing.T) {
	actual := NewBillingAccountID("123456").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/123456"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestBillingAccountID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *BillingAccountId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing Name
			Input: "/providers/Microsoft.Billing/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/providers/Microsoft.Billing/billingAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/123456",
			Expected: &BillingAccountId{
				Name: "123456",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/123456",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := BillingAccountID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
