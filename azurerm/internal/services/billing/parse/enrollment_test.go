package parse

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = EnrollmentId{}

func TestEnrollmentIDFormatter(t *testing.T) {
	actual := NewEnrollmentID("12345678-1234-9876-4563-123456789012").ID()
	expected := "/providers/Microsoft.Billing/enrollmentAccounts/12345678-1234-9876-4563-123456789012"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestEnrollmentID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *EnrollmentId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/",
			Error: true,
		},

		{
			// missing value for EnrollmentAccountName
			Input: "/providers/Microsoft.Billing/enrollmentAccounts/",
			Error: true,
		},

		{
			// valid
			Input: "/providers/Microsoft.Billing/enrollmentAccounts/12345678-1234-9876-4563-123456789012",
			Expected: &EnrollmentId{
				EnrollmentAccountName: "12345678-1234-9876-4563-123456789012",
			},
		},

		{
			// upper-cased
			Input: "/PROVIDERS/MICROSOFT.BILLING/ENROLLMENTACCOUNTS/12345678-1234-9876-4563-123456789012",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := EnrollmentID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.EnrollmentAccountName != v.Expected.EnrollmentAccountName {
			t.Fatalf("Expected %q but got %q for EnrollmentAccountName", v.Expected.EnrollmentAccountName, actual.EnrollmentAccountName)
		}
	}
}
