package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = MicrosoftCustomerAccountBillingScopeId{}

func TestMicrosoftCustomerAccountBillingScopeIDFormatter(t *testing.T) {
	actual := NewMCABillingScopeID("e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31", "PE2Q-NOIT-BG7-TGB", "MTT4-OBS7-PJA-TGB").ID()
	expected := "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/invoiceSections/MTT4-OBS7-PJA-TGB"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestMicrosoftCustomerAccountBillingScopeID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *MicrosoftCustomerAccountBillingScopeId
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
			// missing BillingProfileName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/",
			Error: true,
		},

		{
			// missing value for BillingProfileName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/",
			Error: true,
		},

		{
			// missing InvoiceSectionName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/",
			Error: true,
		},

		{
			// missing value for InvoiceSectionName
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/invoiceSections/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/billingAccounts/e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31/billingProfiles/PE2Q-NOIT-BG7-TGB/invoiceSections/MTT4-OBS7-PJA-TGB",
			Expected: &MicrosoftCustomerAccountBillingScopeId{
				BillingAccountName: "e879cf0f-2b4d-5431-109a-f72fc9868693:024cabf4-7321-4cf9-be59-df0c77ca51de_2019-05-31",
				BillingProfileName: "PE2Q-NOIT-BG7-TGB",
				InvoiceSectionName: "MTT4-OBS7-PJA-TGB",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/BILLINGACCOUNTS/E879CF0F-2B4D-5431-109A-F72FC9868693:024CABF4-7321-4CF9-BE59-DF0C77CA51DE_2019-05-31/BILLINGPROFILES/PE2Q-NOIT-BG7-TGB/INVOICESECTIONS/MTT4-OBS7-PJA-TGB",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := MicrosoftCustomerAccountBillingScopeID(v.Input)
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
		if actual.BillingProfileName != v.Expected.BillingProfileName {
			t.Fatalf("Expected %q but got %q for BillingProfileName", v.Expected.BillingProfileName, actual.BillingProfileName)
		}
		if actual.InvoiceSectionName != v.Expected.InvoiceSectionName {
			t.Fatalf("Expected %q but got %q for InvoiceSectionName", v.Expected.InvoiceSectionName, actual.InvoiceSectionName)
		}
	}
}
