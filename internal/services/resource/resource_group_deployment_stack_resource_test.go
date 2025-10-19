// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatresourcegroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ResourceGroupDeploymentStackResource struct{}

func TestAccResourceGroupDeploymentStack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_deployment_stack", "test")
	r := ResourceGroupDeploymentStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccResourceGroupDeploymentStack_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_deployment_stack", "test")
	r := ResourceGroupDeploymentStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccResourceGroupDeploymentStack_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_deployment_stack", "test")
	r := ResourceGroupDeploymentStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("bypass_stack_out_of_sync_error"),
	})
}

func TestAccResourceGroupDeploymentStack_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_resource_group_deployment_stack", "test")
	r := ResourceGroupDeploymentStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep("bypass_stack_out_of_sync_error"),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ResourceGroupDeploymentStackResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deploymentstacksatresourcegroup.ParseProviderDeploymentStackID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.DeploymentStacksResourceGroupClient.DeploymentStacksGetAtResourceGroup(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ResourceGroupDeploymentStackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource_group_deployment_stack" "test" {
  name                = "acctest-ds-%d"
  resource_group_name = azurerm_resource_group.test.name

  template_content = jsonencode({
    "$schema"      = "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#"
    contentVersion = "1.0.0.0"
    resources      = []
  })

  action_on_unmanage {
    resources = "detach"
  }

  deny_settings {
    mode = "none"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r ResourceGroupDeploymentStackResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_resource_group_deployment_stack" "import" {
  name                = azurerm_resource_group_deployment_stack.test.name
  resource_group_name = azurerm_resource_group_deployment_stack.test.resource_group_name

  template_content = azurerm_resource_group_deployment_stack.test.template_content

  action_on_unmanage {
    resources = "detach"
  }

  deny_settings {
    mode = "none"
  }
}
`, r.basic(data))
}

func (r ResourceGroupDeploymentStackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_resource_group_deployment_stack" "test" {
  name                = "acctest-ds-%d"
  resource_group_name = azurerm_resource_group.test.name
  description         = "Test deployment stack"

  template_content = jsonencode({
    "$schema"      = "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#"
    contentVersion = "1.0.0.0"
    parameters = {
      storageAccountName = {
        type = "string"
      }
    }
    resources = [
      {
        type       = "Microsoft.Storage/storageAccounts"
        apiVersion = "2023-01-01"
        name       = "[parameters('storageAccountName')]"
        location   = "[resourceGroup().location]"
        sku = {
          name = "Standard_LRS"
        }
        kind = "StorageV2"
      }
    ]
  })

  parameters_content = jsonencode({
    storageAccountName = {
      value = "acctest%s"
    }
  })

  action_on_unmanage {
    resources        = "delete"
    resource_groups  = "detach"
    management_groups = "detach"
  }

  deny_settings {
    mode                   = "denyDelete"
    apply_to_child_scopes = false
  }

  bypass_stack_out_of_sync_error = true

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomString)
}
