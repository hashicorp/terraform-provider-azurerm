package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfileHCRPAssignment_CreateorupdateHCRPconfigurationprofileassignment_mockserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_hcrpassignment", "test")
	r := AutomanageConfigurationProfileHCRPAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.CreateorupdateHCRPconfigurationprofileassignment_mockserver(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileHCRPAssignmentResource) CreateorupdateHCRPconfigurationprofileassignment_mockserver() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_automanage_configuration_profile_hcrpassignment" "test" {
  name = "default"
  resource_group_name = "myResourceGroupName"
  machine_name = "myMachineName"
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
`
}
