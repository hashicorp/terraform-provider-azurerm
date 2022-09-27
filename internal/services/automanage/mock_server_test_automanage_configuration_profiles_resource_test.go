package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfile_Createorupdateconfigurationprofile_mockserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.Createorupdateconfigurationprofile_mockserver(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileResource) Createorupdateconfigurationprofile_mockserver() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_automanage_configuration_profile" "test" {
  name = "customConfigurationProfile"
  resource_group_name = "myResourceGroupName"
  location = "East US"
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    Organization = "Administration"
  }
}
`
}
