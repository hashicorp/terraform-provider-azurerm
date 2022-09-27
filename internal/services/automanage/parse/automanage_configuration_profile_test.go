package parse

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = AutomanageConfigurationProfileId{}

func TestAutomanageConfigurationProfileIDFormatter(t *testing.T) {
	actual := NewAutomanageConfigurationProfileID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "configurationProfile1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomanageConfigurationProfileID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomanageConfigurationProfileId
	}{
		{
			// empty
			Input: "",
			Error: true,
		},
		{
			// missing subscriptions
			Input: "/",
			Error: true,
		},
		{
			// missing value for subscriptions
			Input: "/subscriptions/",
			Error: true,
		},
		{
			// missing resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/",
			Error: true,
		},
		{
			// missing value for resourceGroups
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/",
			Error: true,
		},
		{
			// missing configurationProfiles
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/",
			Error: true,
		},
		{
			// missing value for configurationProfiles
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1",
			Expected: &AutomanageConfigurationProfileId{
				SubscriptionId: "12345678-1234-9876-4563-123456789012",
				ResourceGroup:  "resourceGroup1",
				Name:           "configurationProfile1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.AUTOMANAGE/CONFIGURATIONPROFILES/CONFIGURATIONPROFILE1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomanageConfigurationProfileID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
