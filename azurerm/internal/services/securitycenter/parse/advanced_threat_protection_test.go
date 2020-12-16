package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = AdvancedThreatProtectionId{}

func TestAdvancedThreatProtectionIDFormatter(t *testing.T) {
	actual := NewAdvancedThreatProtectionID("12345678-1234-9876-4563-123456789012", "setting1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/advancedThreatProtectionSettings/setting1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAdvancedThreatProtectionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AdvancedThreatProtectionId
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
			// missing AdvancedThreatProtectionSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/",
			Error: true,
		},

		{
			// missing value for AdvancedThreatProtectionSettingName
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/advancedThreatProtectionSettings/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Security/advancedThreatProtectionSettings/setting1",
			Expected: &AdvancedThreatProtectionId{
				SubscriptionId:                      "12345678-1234-9876-4563-123456789012",
				AdvancedThreatProtectionSettingName: "setting1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.SECURITY/ADVANCEDTHREATPROTECTIONSETTINGS/SETTING1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AdvancedThreatProtectionID(v.Input)
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
		if actual.AdvancedThreatProtectionSettingName != v.Expected.AdvancedThreatProtectionSettingName {
			t.Fatalf("Expected %q but got %q for AdvancedThreatProtectionSettingName", v.Expected.AdvancedThreatProtectionSettingName, actual.AdvancedThreatProtectionSettingName)
		}
	}
}
