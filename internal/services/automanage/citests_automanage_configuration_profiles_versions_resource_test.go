package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfilesVersion_Createorupdateconfigurationprofileversion_ci(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.Createorupdateconfigurationprofileversion_ci(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfilesVersionResource) Createorupdateconfigurationprofileversion_ci(data acceptance.TestData) string {
	template := r.citemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profiles_version" "test" {
  name = "acctest-acpv-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration_profile_name = azurerm_automanage_configuration_profile.test.name
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r AutomanageConfigurationProfilesVersionResource) citemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}

resource "azurerm_automanage_configuration_profile" "test" {
  name = "acctest-acp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    ENV = "Test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
