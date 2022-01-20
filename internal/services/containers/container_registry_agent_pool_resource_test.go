package containers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ContainerRegistryAgentPoolResource struct{}

func TestAccContainerRegistryAgentPool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_agent_pool", "test")
	r := ContainerRegistryAgentPoolResource{}

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

func TestAccContainerRegistryAgentPool_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_agent_pool", "test")
	r := ContainerRegistryAgentPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryAgentPool_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_agent_pool", "test")
	r := ContainerRegistryAgentPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryAgentPool_subnet(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_agent_pool", "test")
	r := ContainerRegistryAgentPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subnet(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccContainerRegistryAgentPool_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_container_registry_agent_pool", "test")
	r := ContainerRegistryAgentPoolResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (r ContainerRegistryAgentPoolResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	client := clients.Containers.ContainerRegistryAgentPoolsClient

	id, err := parse.AgentPoolID(state.ID)
	if err != nil {
		return nil, err
	}

	if resp, err := client.Get(ctx, id.ResourceGroup, id.RegistryName, id.Name); err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return utils.Bool(true), nil
}

func (r ContainerRegistryAgentPoolResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_agent_pool" "test" {
  name                  = "acctestACRAP%d"
  container_registry_id = azurerm_container_registry.test.id
  tier                  = "S1"
}
`, template, data.RandomIntOfLength(8))
}

func (r ContainerRegistryAgentPoolResource) requiresImport(data acceptance.TestData) string {
	template := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_agent_pool" "import" {
  name                  = azurerm_container_registry_agent_pool.test.name
  container_registry_id = azurerm_container_registry_agent_pool.test.container_registry_id
  tier                  = azurerm_container_registry_agent_pool.test.tier
}
`, template)
}

func (r ContainerRegistryAgentPoolResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry_agent_pool" "test" {
  name                  = "acctestACRAP%d"
  container_registry_id = azurerm_container_registry.test.id
  tier                  = "S1"
  agent_count           = 2
  os                    = "Linux"
}
`, template, data.RandomIntOfLength(8))
}

func (r ContainerRegistryAgentPoolResource) subnet(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_virtual_network" "test" {
  name                = "acctestACRAP%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "acctestACRAP%d"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.2.0/24"
}

resource "azurerm_container_registry_agent_pool" "test" {
  name                  = "acctestACRAP%d"
  container_registry_id = azurerm_container_registry.test.id
  tier                  = "S1"
  subnet_id             = azurerm_subnet.test.id
}
`, template, data.RandomIntOfLength(8), data.RandomIntOfLength(8), data.RandomIntOfLength(8))
}

func (r ContainerRegistryAgentPoolResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-acrap-%d"
  location = "%s"
}

resource "azurerm_container_registry" "test" {
  name                = "acctestACRAP%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
}
`, data.RandomIntOfLength(8), data.Locations.Primary, data.RandomIntOfLength(8))
}
