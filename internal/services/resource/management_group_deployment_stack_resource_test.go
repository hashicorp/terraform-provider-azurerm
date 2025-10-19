// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatmanagementgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagementGroupDeploymentStackResource struct{}

func TestAccManagementGroupDeploymentStack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_deployment_stack", "test")
	r := ManagementGroupDeploymentStackResource{}

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

func TestAccManagementGroupDeploymentStack_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_deployment_stack", "test")
	r := ManagementGroupDeploymentStackResource{}

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

func TestAccManagementGroupDeploymentStack_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_deployment_stack", "test")
	r := ManagementGroupDeploymentStackResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccManagementGroupDeploymentStack_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_management_group_deployment_stack", "test")
	r := ManagementGroupDeploymentStackResource{}

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
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r ManagementGroupDeploymentStackResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deploymentstacksatmanagementgroup.ParseProviders2DeploymentStackID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.DeploymentStacksManagementGroupClient.DeploymentStacksGetAtManagementGroup(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ManagementGroupDeploymentStackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%[1]d"
}

resource "azurerm_management_group_deployment_stack" "test" {
  name                = "acctestds-%[1]d"
  location            = %[2]q
  management_group_id = azurerm_management_group.test.id

  template_content = jsonencode({
    "$schema"        = "https://schema.management.azure.com/schemas/2019-08-01/managementGroupDeploymentTemplate.json#"
    "contentVersion" = "1.0.0.0"
    "resources"      = []
  })

  action_on_unmanage {
    resources = "delete"
  }

  deny_settings {
    mode = "none"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r ManagementGroupDeploymentStackResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_management_group_deployment_stack" "import" {
  name                = azurerm_management_group_deployment_stack.test.name
  location            = azurerm_management_group_deployment_stack.test.location
  management_group_id = azurerm_management_group_deployment_stack.test.management_group_id

  template_content = azurerm_management_group_deployment_stack.test.template_content

  action_on_unmanage {
    resources = "delete"
  }

  deny_settings {
    mode = "none"
  }
}
`, r.basic(data))
}

func (r ManagementGroupDeploymentStackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-%[1]d"
}

resource "azurerm_management_group_deployment_stack" "test" {
  name                       = "acctestds-%[1]d"
  location                   = %[2]q
  management_group_id        = azurerm_management_group.test.id
  description                = "Test deployment stack"
  deployment_subscription_id = data.azurerm_client_config.current.subscription_id

  template_content = jsonencode({
    "$schema"        = "https://schema.management.azure.com/schemas/2019-08-01/managementGroupDeploymentTemplate.json#"
    "contentVersion" = "1.0.0.0"
    "parameters" = {
      "policyName" = {
        "type" = "string"
      }
    }
    "resources" = [
      {
        "type"       = "Microsoft.Authorization/policyDefinitions"
        "apiVersion" = "2021-06-01"
        "name"       = "[parameters('policyName')]"
        "properties" = {
          "policyType"  = "Custom"
          "mode"        = "All"
          "displayName" = "Test Policy"
          "policyRule" = {
            "if" = {
              "field"  = "type"
              "equals" = "Microsoft.Storage/storageAccounts"
            }
            "then" = {
              "effect" = "audit"
            }
          }
        }
      }
    ]
  })

  parameters_content = jsonencode({
    "policyName" = {
      "value" = "acctest-policy-%[1]d"
    }
  })

  action_on_unmanage {
    resources         = "delete"
    resource_groups   = "delete"
    management_groups = "delete"
  }

  deny_settings {
    mode                  = "denyWriteAndDelete"
    apply_to_child_scopes = true
    excluded_actions      = ["Microsoft.Authorization/policyDefinitions/delete"]
  }

  bypass_stack_out_of_sync_error = true

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
