package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type ContainerRegistryReplicationResource struct {
}

func TestAccContainerRegistryReplication_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_replication", "test")
	r := ContainerRegistryReplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryReplication_basicUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_replication", "test")
	r := ContainerRegistryReplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.%").HasValue("2"),
				check.That(data.ResourceName).Key("tags.label").HasValue("test1"),
				check.That(data.ResourceName).Key("tags.ENV").HasValue("prod"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryReplication_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_replication", "test")
	r := ContainerRegistryReplicationResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_container_registry_replication"),
		},
	})
}

func (ContainerRegistryReplicationResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrReplicationTest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Premium"
}

resource "azurerm_container_registry_replication" "test" {
  name                = "testReplication%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  tags = {
    label = "test"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (ContainerRegistryReplicationResource) basicUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "rg" {
  name     = "acctestRG-aks-%d"
  location = "%s"
}

resource "azurerm_container_registry" "acr" {
  name                = "acrReplicationTest%d"
  resource_group_name = azurerm_resource_group.rg.name
  location            = "%s"
  sku                 = "Premium"
}

resource "azurerm_container_registry_replication" "test" {
  name                = "testReplication%d"
  resource_group_name = azurerm_resource_group.rg.name
  registry_name       = azurerm_container_registry.acr.name
  location            = "%s"

  tags = {
    label = "test1"
    ENV   = "prod"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.Locations.Secondary)
}

func (r ContainerRegistryReplicationResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_replication" "import" {
  name                = azurerm_container_registry_replication.test.name
  resource_group_name = azurerm_container_registry_replication.test.resource_group_name
  registry_name       = azurerm_container_registry_replication.test.registry_name
  location            = azurerm_container_registry_replication.test.location
}
`, r.basic(data))
}

func (t ContainerRegistryReplicationResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	registryName := id.Path["registries"]
	name := id.Path["replications"]

	resp, err := clients.Containers.ReplicationsClient.Get(ctx, resourceGroup, registryName, name)
	if err != nil {
		return nil, fmt.Errorf("reading Container Registry Replication (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
