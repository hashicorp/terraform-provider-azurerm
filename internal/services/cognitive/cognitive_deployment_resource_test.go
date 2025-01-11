// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitive_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2024-10-01/deployments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CognitiveDeploymentTestResource struct{}

func TestAccCognitiveDeployment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := CognitiveDeploymentTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveDeployment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")

	r := CognitiveDeploymentTestResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccCognitiveDeployment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := CognitiveDeploymentTestResource{}
	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("dynamic_throttling_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccCognitiveDeployment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_cognitive_deployment", "test")
	r := CognitiveDeploymentTestResource{}

	data.ResourceSequentialTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.versionUpgradeOption(data, "OnceNewDefaultVersionAvailable"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.versionUpgradeOption(data, "OnceCurrentVersionExpired"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.updateVersion(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r CognitiveDeploymentTestResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := deployments.ParseDeploymentID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Cognitive.DeploymentsClient
	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return utils.Bool(resp.Model != nil), nil
}

func (r CognitiveDeploymentTestResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
resource "azurerm_cognitive_account" "test" {
  name                = "acctest-ca-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  kind                = "OpenAI"
  sku_name            = "S0"
}
`, data.RandomInteger, data.Locations.Secondary, data.RandomInteger)
}

func (r CognitiveDeploymentTestResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  model {
    format = "OpenAI"
    name   = "text-embedding-ada-002"
  }
  sku {
    name = "Standard"
  }
  lifecycle {
    ignore_changes = [model.0.version]
  }
}
`, template, data.RandomInteger)
}

func (r CognitiveDeploymentTestResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "import" {
  name                 = azurerm_cognitive_deployment.test.name
  cognitive_account_id = azurerm_cognitive_account.test.id
  model {
    format  = "OpenAI"
    name    = "text-embedding-ada-002"
    version = "2"
  }
  sku {
    name = "Standard"
  }
}
`, config)
}

func (r CognitiveDeploymentTestResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "test" {
  name                       = "acctest-cd-%d"
  cognitive_account_id       = azurerm_cognitive_account.test.id
  dynamic_throttling_enabled = true
  model {
    format  = "OpenAI"
    name    = "text-embedding-ada-002"
    version = "2"
  }
  sku {
    name = "Standard"
  }
  rai_policy_name        = "Microsoft.DefaultV2"
  version_upgrade_option = "OnceNewDefaultVersionAvailable"
}
`, template, data.RandomInteger)
}

func (r CognitiveDeploymentTestResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  rai_policy_name      = "Microsoft.Default"
  model {
    format  = "OpenAI"
    name    = "text-embedding-ada-002"
    version = "2"
  }
  sku {
    name     = "Standard"
    capacity = 2
  }
}
`, template, data.RandomInteger)
}

func (r CognitiveDeploymentTestResource) updateVersion(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "test" {
  name                 = "acctest-cd-%d"
  cognitive_account_id = azurerm_cognitive_account.test.id
  rai_policy_name      = "Microsoft.Default"
  model {
    format  = "OpenAI"
    name    = "text-embedding-ada-002"
    version = "1"
  }
  sku {
    name     = "Standard"
    capacity = 2
  }
}
`, template, data.RandomInteger)
}

func (r CognitiveDeploymentTestResource) versionUpgradeOption(data acceptance.TestData, versionUpgradeOption string) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_cognitive_deployment" "test" {
  name                   = "acctest-cd-%d"
  cognitive_account_id   = azurerm_cognitive_account.test.id
  rai_policy_name        = "Microsoft.Default"
  version_upgrade_option = "%s"
  model {
    format  = "OpenAI"
    name    = "text-embedding-ada-002"
    version = "1"
  }
  sku {
    name     = "Standard"
    capacity = 2
  }
}
`, template, data.RandomInteger, versionUpgradeOption)
}
