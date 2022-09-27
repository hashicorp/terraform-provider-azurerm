package automanage_test

import (
	"testing"
)

func TestAccAutomanageConfigurationProfileAssignment_Createorupdateconfigurationprofileassignment_mockserver(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration_profile_assignment", "test")
	r := AutomanageConfigurationProfileAssignmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.Createorupdateconfigurationprofileassignment_mockserver(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r AutomanageConfigurationProfileAssignmentResource) Createorupdateconfigurationprofileassignment_mockserver() string {
	return `
provider "azurerm" {
  features {}
}

resource "azurerm_automanage_configuration_profile_assignment" "test" {
  name = "default"
  resource_group_name = "myResourceGroupName"
  vm_name = "myVMName"
  configuration_profile = "/providers/Microsoft.Automanage/bestPractices/AzureBestPracticesProduction"
}
`
}
