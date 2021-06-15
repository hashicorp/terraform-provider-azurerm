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

type TenantTemplateDeploymentResource struct {
}

func TestAccTenantTemplateDeployment_empty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_tenant_template_deployment", "test")
	if data.Client().IsServicePrincipal {
		t.Skip("Skipping due to permissions unavailable on tenant scope")
	}
	r := TenantTemplateDeploymentResource{}

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

func (t TenantTemplateDeploymentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.SubscriptionTemplateDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Resource.DeploymentsClient.GetAtTenantScope(ctx, id.DeploymentName)
	if err != nil {
		return nil, fmt.Errorf("reading Tenant Template Deployment (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}

func (TenantTemplateDeploymentResource) emptyConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_tenant_template_deployment" "test" {
  name     = "acctestTenantDeploy-%d"
  location = %q

  template_content = <<TEMPLATE
{
  "$schema": "https://pluginsdk.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
  "contentVersion": "1.0.0.0",
  "parameters": {},
  "variables": {},
  "resources": []
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary)
}

func (TenantTemplateDeploymentResource) emptyWithTagsConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_tenant_template_deployment" "test" {
  name     = "acctestsubdeploy-%d"
  location = %q
  tags = {
    Hello = "World"
  }

  template_content = <<TEMPLATE
{
 "$schema": "https://pluginsdk.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
 "contentVersion": "1.0.0.0",
 "parameters": {},
 "variables": {},
 "resources": []
}
TEMPLATE
}
`, data.RandomInteger, data.Locations.Primary)
}
