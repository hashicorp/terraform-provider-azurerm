package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.Id = PublicMaintenanceConfigurationId{}

func TestPublicMaintenanceConfigurationIDFormatter(t *testing.T) {
	actual := NewPublicMaintenanceConfigurationID("12345678-1234-9876-4563-123456789012", "publicMaintenanceConfiguration1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/publicMaintenanceConfiguration1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestPublicMaintenanceConfigurationID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *PublicMaintenanceConfigurationId
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
			// missing Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Maintenance/",
			Error: true,
		},

		{
			// missing value for Name
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/",
			Error: true,
		},

		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/publicMaintenanceConfiguration1",
			Expected: &PublicMaintenanceConfigurationId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				Name:           "publicMaintenanceConfiguration1",
			},
		},

		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/PROVIDERS/MICROSOFT.MAINTENANCE/PUBLICMAINTENANCECONFIGURATIONS/PUBLICMAINTENANCECONFIGURATION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := PublicMaintenanceConfigurationID(v.Input)
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
		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}
	}
}
