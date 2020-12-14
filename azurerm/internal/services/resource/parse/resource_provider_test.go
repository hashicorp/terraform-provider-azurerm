package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

var _ resourceid.Formatter = ResourceProviderId{}

func TestResourceProviderIDFormatter(t *testing.T) {
	actual := NewResourceProviderID("12345678-1234-9876-4563-123456789012", "Instruments.Didgeridoo").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Instruments.Didgeridoo"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestResourceProviderID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *ResourceProviderId
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
			// missing Providers
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},

		{
			// missing value for Providers
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Instruments.Didgeridoo",
			Expected: &ResourceProviderId{
				SubscriptionId:   "12345678-1234-9876-4563-123456789012",
				ResourceProvider: "Instruments.Didgeridoo",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/INSTRUMENTS.DIDGERIDOO",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := ResourceProviderID(v.Input)
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
		if actual.ResourceProvider != v.Expected.ResourceProvider {
			t.Fatalf("Expected %q but got %q for ResourceProvider", v.Expected.ResourceProvider, actual.ResourceProvider)
		}
	}
}
