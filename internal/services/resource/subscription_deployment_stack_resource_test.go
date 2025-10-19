// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2024-03-01/deploymentstacksatsubscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SubscriptionDeploymentStackResource struct{}

func TestAccSubscriptionDeploymentStack_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_deployment_stack", "test")
	r := SubscriptionDeploymentStackResource{}

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

func TestAccSubscriptionDeploymentStack_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_deployment_stack", "test")
	r := SubscriptionDeploymentStackResource{}

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

func TestAccSubscriptionDeploymentStack_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_deployment_stack", "test")
	r := SubscriptionDeploymentStackResource{}

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

func TestAccSubscriptionDeploymentStack_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_subscription_deployment_stack", "test")
	r := SubscriptionDeploymentStackResource{}

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

func (r SubscriptionDeploymentStackResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deploymentstacksatsubscription.ParseDeploymentStackID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Resource.DeploymentStacksSubscriptionClient.DeploymentStacksGetAtSubscription(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r SubscriptionDeploymentStackResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_subscription_deployment_stack" "test" {
  name     = "acctestds-%d"
  location = %q

  template_content = jsonencode({
    "$schema" = "https://schema.management.azure.com/schemas/2018-05-01/subscriptionDeploymentTemplate.json#"
    "contentVersion" = "1.0.0.0"
    "resources" = []
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

func (r SubscriptionDeploymentStackResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_subscription_deployment_stack" "import" {
  name     = azurerm_subscription_deployment_stack.test.name
  location = azurerm_subscription_deployment_stack.test.location

  template_content = azurerm_subscription_deployment_stack.test.template_content

  action_on_unmanage {
    resources = "delete"
  }

  deny_settings {
    mode = "none"
  }
}
`, r.basic(data))
}

func (r SubscriptionDeploymentStackResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = %[2]q
}

resource "azurerm_subscription_deployment_stack" "test" {
  name                        = "acctestds-%[1]d"
  location                    = %[2]q
  description                 = "Test deployment stack"
  deployment_resource_group_name = azurerm_resource_group.test.name

  template_content = jsonencode({
    "$schema" = "https://schema.management.azure.com/schemas/2018-05-01/subscriptionDeploymentTemplate.json#"
    "contentVersion" = "1.0.0.0"
    "parameters" = {
      "storageAccountName" = {
        "type" = "string"
      }
    }
    "resources" = [
      {
        "type" = "Microsoft.Storage/storageAccounts"
        "apiVersion" = "2021-04-01"
        "name" = "[parameters('storageAccountName')]"
        "location" = %[2]q
        "sku" = {
          "name" = "Standard_LRS"
        }
        "kind" = "StorageV2"
      }
    ]
  })

  parameters_content = jsonencode({
    "storageAccountName" = {
      "value" = "acctest%[3]s"
    }
  })

  action_on_unmanage {
    resources        = "delete"
    resource_groups  = "delete"
    management_groups = "detach"
  }

  deny_settings {
    mode                 = "denyDelete"
    apply_to_child_scopes = true
  }

  bypass_stack_out_of_sync_error = true

  tags = {
    environment = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}
