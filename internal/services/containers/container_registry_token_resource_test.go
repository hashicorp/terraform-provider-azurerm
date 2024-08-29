// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/tokens"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryTokenResource struct{}

func TestAccContainerRegistryToken_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")
	r := ContainerRegistryTokenResource{}

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

func TestAccContainerRegistryToken_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")
	r := ContainerRegistryTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry_token"),
		},
	})
}

func TestAccContainerRegistryToken_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")
	r := ContainerRegistryTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryToken_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_token", "test")
	r := ContainerRegistryTokenResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data, true),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data, false),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (t ContainerRegistryTokenResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := tokens.ParseTokenID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.ContainerRegistryClient_v2023_06_01_preview.Tokens.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ContainerRegistryTokenResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}

# use system wide scope map for tests
data "azurerm_container_registry_scope_map" "pull_repos" {
  name                    = "_repositories_pull"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testtoken-%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ContainerRegistryTokenResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_token" "import" {
  name                    = azurerm_container_registry_token.test.name
  resource_group_name     = azurerm_container_registry_token.test.resource_group_name
  container_registry_name = azurerm_container_registry_token.test.container_registry_name
  scope_map_id            = azurerm_container_registry_token.test.scope_map_id
}
`, r.basic(data))
}

func (r ContainerRegistryTokenResource) complete(data acceptance.TestData, status bool) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "testacccr%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  admin_enabled       = false
  sku                 = "Premium"

  tags = {
    environment = "production"
  }
}

# use system wide scope map for tests
data "azurerm_container_registry_scope_map" "pull_repos" {
  name                    = "_repositories_pull"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testtoken%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  scope_map_id            = data.azurerm_container_registry_scope_map.pull_repos.id
  enabled                 = %t
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, status)
}
