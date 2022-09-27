package parse

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	"testing"
)

var _ resourceid.Formatter = AutomanageConfigurationProfilesVersionId{}

func TestAutomanageConfigurationProfilesVersionIDFormatter(t *testing.T) {
	actual := NewAutomanageConfigurationProfilesVersionID("12345678-1234-9876-4563-123456789012", "resourceGroup1", "configurationProfile1", "version1").ID()
	expected := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1/versions/version1"
	if actual != expected {
		t.Fatalf("Expected %q but got %q", expected, actual)
	}
}

func TestAutomanageConfigurationProfilesVersionID(t *testing.T) {
	testData := []struct {
		Input    string
		Error    bool
		Expected *AutomanageConfigurationProfilesVersionId
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
			// missing versions
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1/",
			Error: true,
		},
		{
			// missing value for versions
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1/versions/",
			Error: true,
		},
		{
			// valid
			Input: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resourceGroup1/providers/Microsoft.Automanage/configurationProfiles/configurationProfile1/versions/version1",
			Expected: &AutomanageConfigurationProfilesVersionId{
				SubscriptionId:           "12345678-1234-9876-4563-123456789012",
				ResourceGroup:            "resourceGroup1",
				ConfigurationProfileName: "configurationProfile1",
				Name:                     "version1",
			},
		},
		{
			// upper-cased
			Input: "/SUBSCRIPTIONS/12345678-1234-9876-4563-123456789012/RESOURCEGROUPS/RESOURCEGROUP1/PROVIDERS/MICROSOFT.AUTOMANAGE/CONFIGURATIONPROFILES/CONFIGURATIONPROFILE1/VERSIONS/VERSION1",
			Error: true,
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing %q", v.Input)

		actual, err := AutomanageConfigurationProfilesVersionID(v.Input)
		if err != nil {
			if v.Error {
				continue
			}
			t.Fatalf("Expected a value but got an error: %s", err)
		}

		if actual.SubscriptionId != v.Expected.SubscriptionId {
			t.Fatalf("Expected %q but got %q for SubscriptionId", v.Expected.SubscriptionId, actual.SubscriptionId)
		}

		if actual.ConfigurationProfileName != v.Expected.ConfigurationProfileName {
			t.Fatalf("Expected %q but got %q for ConfigurationProfileName", v.Expected.ConfigurationProfileName, actual.ConfigurationProfileName)
		}

		if actual.Name != v.Expected.Name {
			t.Fatalf("Expected %q but got %q for Name", v.Expected.Name, actual.Name)
		}

		if actual.ResourceGroup != v.Expected.ResourceGroup {
			t.Fatalf("Expected %q but got %q for ResourceGroup", v.Expected.ResourceGroup, actual.ResourceGroup)
		}
	}
}
