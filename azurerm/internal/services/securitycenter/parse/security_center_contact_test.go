package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = SecurityCenterContactId{}

func TestSecurityCenterContactIDFormatter(t *testing.T) {
	actual := NewSecurityCenterContactID("12345678-1234-9876-4563-123456789012", "contact1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/securityContacts/contact1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestSecurityCenterContactID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *SecurityCenterContactId
	}{

		{
			// empty
			Input: "",
			Error: true,
		},

		{
			// missing SubscriptionId
			Input: "/",
			Error: true,
		},

		{
			// missing value for SubscriptionId
			Input: "/subscriptions/",
			Error: true,
		},

		{
			// missing SecurityContactName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/",
			Error: true,
		},

		{
			// missing value for SecurityContactName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/securityContacts/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/securityContacts/contact1",
			Expected: &SecurityCenterContactId{
				SubscriptionId:      "12345678-1234-9876-4563-123456789012",
				SecurityContactName: "contact1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.SECURITY/SECURITYCONTACTS/CONTACT1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := SecurityCenterContactID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}

			t.Fatalf("Expect a value but got an error: %s", err)
		}
		if v.Error {
			t.Fatal("Expect an error but didn't get one")
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}
		if actual.SecurityContactName != v.Expected.SecurityContactName {
			t.Fatalf("Expected %q but got %q for SecurityContactName", v.Expected.SecurityContactName, actual.SecurityContactName)
		}
	}
}
