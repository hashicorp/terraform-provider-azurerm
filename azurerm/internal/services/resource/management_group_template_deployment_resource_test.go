package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/resource/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ManagementGroupTemplateDeploymentResource struct{}

func TestAccManagementGroupTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_template_deployment", "test")
	r := ManagementGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			// set some tags
			Config: r.emptyWithTagsConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupTemplateDeployment_templateSpec(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_template_deployment", "test")
	r := ManagementGroupTemplateDeploymentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.templateSpecVersionConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (ManagementGroupTemplateDeploymentResource) templateSpecVersionConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "TestAcc-Deployment-%[1]d"
}

data "azurerm_template_spec_version" "test" {
  name                = "acctest-standing-data-empty"
  resource_group_name = "standing-data-for-acctest"
  version             = "v1.0.0"
}

resource "azurerm_management_group_template_deployment" "test" {
  name                = "acctestMGdeploy-%[1]d"
  management_group_id = azurerm_management_group.test.id
  location            = %[2]q

  template_spec_version_id = data.azurerm_template_spec_version.test.id

}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (ManagementGroupTemplateDeploymentResource) emptyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "TestAcc-Deployment-%[1]d"
}

resource "azurerm_management_group_template_deployment" "test" {
  name                = "acctestMGdeploy-%[1]d"
  management_group_id = azurerm_management_group.test.id
  location            = %[2]q

  template_content = <<TEMPLATE
{
 "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
 "contentVersion": "1.0.0.0",
 "parameters": {},
 "variables": {},
 "resources": []
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary)
}

func (ManagementGroupTemplateDeploymentResource) emptyWithTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_management_group" "test" {
  name = "TestAcc-Deployment-%[1]d"
}

resource "azurerm_management_group_template_deployment" "test" {
  name                = "acctestMGdeploy-%[1]d"
  management_group_id = azurerm_management_group.test.id
  location            = %q
  tags = {
    Hello = "World"
  }


  template_content = <<TEMPLATE
{
 "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
 "contentVersion": "1.0.0.0",
 "parameters": {},
 "variables": {},
 "resources": []
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t ManagementGroupTemplateDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.ManagementGroupTemplateDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.DeploymentsClient.GetAtManagementGroupScope(ctx, id.ManagementGroupName, id.DeploymentName)
	if err != nil {
		return nil, fmt.Errorf("reading Subscription Template Deployment (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
