// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerregistry/2021-08-01-preview/connectedregistries"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerConnectedRegistryResource struct{}

func TestAccContainerConnectedRegistry_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

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

func TestAccContainerConnectedRegistry_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

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

func TestAccContainerConnectedRegistry_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

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

func TestAccContainerConnectedRegistry_mirror(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.mirror(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerConnectedRegistry_registry(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.registry(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerConnectedRegistry_cascaded(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.cascaded(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerConnectedRegistry_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_connected_registry", "test")
	r := ContainerConnectedRegistryResource{}

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

func (r ContainerConnectedRegistryResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Containers.ContainerRegistryClient_v2021_08_01_preview.ConnectedRegistries

	id, err := connectedregistries.ParseConnectedRegistryID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (r ContainerConnectedRegistryResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_connected_registry" "test" {
  name                  = "testacccrc%d"
  container_registry_id = azurerm_container_registry.test.id
  sync_token_id         = azurerm_container_registry_token.test.id
}
`, template, data.RandomInteger)
}

func (r ContainerConnectedRegistryResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%[1]s

resource "azurerm_container_registry_scope_map" "client" {
  name                    = "testacccrclient%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/delete",
    "repositories/hello-world/content/read",
    "repositories/hello-world/content/write",
    "repositories/hello-world/metadata/read",
    "repositories/hello-world/metadata/write",
  ]
}

resource "azurerm_container_registry_token" "client" {
  name                    = "testacccrtokenclient%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.client.id
}

resource "azurerm_container_connected_registry" "test" {
  name                  = "testacccrc%[2]d"
  container_registry_id = azurerm_container_registry.test.id
  sync_token_id         = azurerm_container_registry_token.test.id
  notification {
    name   = "hello-world"
    tag    = "latest"
    action = "*"
  }
  log_level         = "Debug"
  audit_log_enabled = true
  client_token_ids  = [azurerm_container_registry_token.client.id]

  # This is necessary to make the Terraform apply order works correctly.
  # Without CBD: azurerm_container_registry_token.client (destroy) -> azurerm_container_connected_registry.test (update)
  # 			 (the 1st step wil fail as the token is under used by the connected registry)
  # With CBD   : azurerm_container_connected_registry.test (update) -> azurerm_container_registry_token.client (destroy) 
  lifecycle {
    create_before_destroy = true
  }
}
`, template, data.RandomInteger)
}

func (r ContainerConnectedRegistryResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_connected_registry" "import" {
  name                  = azurerm_container_connected_registry.test.name
  container_registry_id = azurerm_container_connected_registry.test.container_registry_id
  sync_token_id         = azurerm_container_connected_registry.test.sync_token_id
}
`, template)
}

func (r ContainerConnectedRegistryResource) cascaded(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                  = "testacccr%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  sku                   = "Premium"
  data_endpoint_enabled = true
}

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testacccr%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/delete",
    "repositories/hello-world/content/read",
    "repositories/hello-world/content/write",
    "repositories/hello-world/metadata/read",
    "repositories/hello-world/metadata/write",
    "gateway/testacccrc%[2]d/config/read",
    "gateway/testacccrc%[2]d/config/write",
    "gateway/testacccrc%[2]d/message/read",
    "gateway/testacccrc%[2]d/message/write",
    "gateway/testacccrcchild%[2]d/config/read",
    "gateway/testacccrcchild%[2]d/config/write",
    "gateway/testacccrcchild%[2]d/message/read",
    "gateway/testacccrcchild%[2]d/message/write",
  ]
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testacccrtoken%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.test.id
}

resource "azurerm_container_connected_registry" "test" {
  name                  = "testacccrc%[2]d"
  container_registry_id = azurerm_container_registry.test.id
  sync_token_id         = azurerm_container_registry_token.test.id
}

resource "azurerm_container_registry_scope_map" "child" {
  name                    = "testacccrchild%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/read",
    "repositories/hello-world/metadata/read",
    "gateway/testacccrcchild%[2]d/config/read",
    "gateway/testacccrcchild%[2]d/config/write",
    "gateway/testacccrcchild%[2]d/message/read",
    "gateway/testacccrcchild%[2]d/message/write",
  ]
}

resource "azurerm_container_registry_token" "child" {
  name                    = "testacccrtokenchild%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.child.id
}

resource "azurerm_container_connected_registry" "child" {
  name                  = "testacccrcchild%[2]d"
  container_registry_id = azurerm_container_registry.test.id
  parent_registry_id    = azurerm_container_connected_registry.test.id
  sync_token_id         = azurerm_container_registry_token.child.id
  mode                  = "ReadOnly"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r ContainerConnectedRegistryResource) mirror(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                  = "testacccr%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  sku                   = "Premium"
  data_endpoint_enabled = true
}

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testacccr%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/read",
    "repositories/hello-world/metadata/read",
    "gateway/testacccrc%[2]d/config/read",
    "gateway/testacccrc%[2]d/config/write",
    "gateway/testacccrc%[2]d/message/read",
    "gateway/testacccrc%[2]d/message/write",
  ]
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testacccrtoken%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.test.id
}

resource "azurerm_container_connected_registry" "test" {
  name                  = "testacccrc%[2]d"
  container_registry_id = azurerm_container_registry.test.id
  sync_token_id         = azurerm_container_registry_token.test.id
  mode                  = "Mirror"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r ContainerConnectedRegistryResource) registry(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[2]d"
  location = "%[1]s"
}

resource "azurerm_container_registry" "test" {
  name                  = "testacccr%[2]d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  sku                   = "Premium"
  data_endpoint_enabled = true
}

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testacccr%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/delete",
    "repositories/hello-world/content/read",
    "repositories/hello-world/content/write",
    "repositories/hello-world/metadata/read",
    "repositories/hello-world/metadata/write",
    "gateway/testacccrc%[2]d/config/read",
    "gateway/testacccrc%[2]d/config/write",
    "gateway/testacccrc%[2]d/message/read",
    "gateway/testacccrc%[2]d/message/write",
  ]
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testacccrtoken%[2]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.test.id
}

resource "azurerm_container_connected_registry" "test" {
  name                  = "testacccrc%[2]d"
  container_registry_id = azurerm_container_registry.test.id
  sync_token_id         = azurerm_container_registry_token.test.id
  mode                  = "Registry"
}
`, data.Locations.Primary, data.RandomInteger)
}

func (r ContainerConnectedRegistryResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acr-%[1]d"
  location = "%[2]s"
}

resource "azurerm_container_registry" "test" {
  name                  = "testacccr%[1]d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = azurerm_resource_group.test.location
  sku                   = "Premium"
  data_endpoint_enabled = true
}

resource "azurerm_container_registry_scope_map" "test" {
  name                    = "testacccr%[1]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  actions = [
    "repositories/hello-world/content/delete",
    "repositories/hello-world/content/read",
    "repositories/hello-world/content/write",
    "repositories/hello-world/metadata/read",
    "repositories/hello-world/metadata/write",
    "gateway/testacccrc%[1]d/config/read",
    "gateway/testacccrc%[1]d/config/write",
    "gateway/testacccrc%[1]d/message/read",
    "gateway/testacccrc%[1]d/message/write",
  ]
}

resource "azurerm_container_registry_token" "test" {
  name                    = "testacccrtoken%[1]d"
  container_registry_name = azurerm_container_registry.test.name
  resource_group_name     = azurerm_container_registry.test.resource_group_name
  scope_map_id            = azurerm_container_registry_scope_map.test.id
}
`, data.RandomInteger, data.Locations.Primary)
}
