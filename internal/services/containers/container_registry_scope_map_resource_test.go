// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2023-06-01-preview/scopemaps"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryScopeMapResource struct{}

func TestAccContainerRegistryScopeMap_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")
	r := ContainerRegistryScopeMapResource{}

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

func TestAccContainerRegistryScopeMap_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")
	r := ContainerRegistryScopeMapResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("actions.0").HasValue("repositories/testrepo/content/read"),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry_scope_map"),
		},
	})
}

func TestAccContainerRegistryScopeMap_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")
	r := ContainerRegistryScopeMapResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("actions.0").HasValue("repositories/testrepo/content/read"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccontainerRegistryScopeMap_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_scope_map", "test")
	r := ContainerRegistryScopeMapResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("actions.0").HasValue("repositories/testrepo/content/read"),
			),
		},
		data.ImportStep(),
		{
			Config: r.completeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("actions.#").HasValue("2"),
				check.That(data.ResourceName).Key("actions.0").HasValue("repositories/testrepo/content/read"),
				check.That(data.ResourceName).Key("actions.1").HasValue("repositories/testrepo/content/delete"),
			),
		},
		data.ImportStep(),
	})
}

func (ContainerRegistryScopeMapResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scopemaps.ParseScopeMapID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Containers.ContainerRegistryClient_v2023_06_01_preview.ScopeMaps.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (ContainerRegistryScopeMapResource) basic(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (r ContainerRegistryScopeMapResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_scope_map" "import" {
  name                    = azurerm_container_registry_scope_map.test.name
  resource_group_name     = azurerm_container_registry_scope_map.test.resource_group_name
  container_registry_name = azurerm_container_registry_scope_map.test.container_registry_name
  actions                 = azurerm_container_registry_scope_map.test.actions
}
`, r.basic(data))
}

func (ContainerRegistryScopeMapResource) complete(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  description             = "An example scope map"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (ContainerRegistryScopeMapResource) completeUpdated(data acceptance.TestData) string {
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

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testscopemap%d"
  description             = "An example scope map"
  resource_group_name     = azurerm_resource_group.test.name
  container_registry_name = azurerm_container_registry.test.name
  actions                 = ["repositories/testrepo/content/read", "repositories/testrepo/content/delete"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
