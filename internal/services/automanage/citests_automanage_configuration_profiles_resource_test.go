package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfile_Createorupdateconfigurationprofile_ci(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile", "test")
	r := AutomanageConfigurationProfileResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.Createorupdateconfigurationprofile_ci(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileResource) Createorupdateconfigurationprofile_ci(data acceptance.TestData) string {
	template := r.citemplate(data)
	return fmt.Sprintf(`
%s

resource "azurerm_automanage_configuration_profile" "test" {
  name = "acctest-acp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location = azurerm_resource_group.test.location
  configuration = "{\"Antimalware/Enable\":false,\"AzureSecurityCenter/Enable\":true,\"Backup/Enable\":false,\"BootDiagnostics/Enable\":true,\"ChangeTrackingAndInventory/Enable\":true,\"GuestConfiguration/Enable\":true,\"LogAnalytics/Enable\":true,\"UpdateManagement/Enable\":true,\"VMInsights/Enable\":true}"
  tags = {
    ENV = "Test"
  }
}
`, template, data.RandomInteger)
}

func (r AutomanageConfigurationProfileResource) citemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-automanage-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}
