package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfilesVersion_Createorupdateconfigurationprofileversion_mockserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profiles_version", "test")
	r := AutomanageConfigurationProfilesVersionResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.Createorupdateconfigurationprofileversion_mockserver(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfilesVersionResource) Createorupdateconfigurationprofileversion_mockserver() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_automanage_configuration_profiles_version" "test" {
  name = "version1"
  resource_group_name = "myResourceGroupName"
  location = "East US"
  configuration_profile_name = "customConfigurationProfile"
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    Organization = "Administration"
  }
}
`
}
