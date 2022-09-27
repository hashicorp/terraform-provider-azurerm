package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfileHCIAssignment_CreateorupdateaHCIconfigurationprofileassignment_mockserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_hciassignment", "test")
	r := AutomanageConfigurationProfileHCIAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.CreateorupdateaHCIconfigurationprofileassignment_mockserver(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileHCIAssignmentResource) CreateorupdateaHCIconfigurationprofileassignment_mockserver() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_automanage_configuration_profile_hciassignment" "test" {
  name = "default"
  resource_group_name = "myResourceGroupName"
  cluster_name = "myClusterName"
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
`
}
