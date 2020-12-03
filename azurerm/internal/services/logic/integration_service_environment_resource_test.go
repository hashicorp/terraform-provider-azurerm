package logic_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type IntegrationServiceEnvironmentResource struct {
}

func TestAccIntegrationServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	r := IntegrationServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIntegrationServiceEnvironment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	r := IntegrationServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("development"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIntegrationServiceEnvironment_developer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	r := IntegrationServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.developer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Developer_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("development"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIntegrationServiceEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	r := IntegrationServiceEnvironmentResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.skuName(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_1"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.ImportStep(),
	})
}

func TestAccIntegrationServiceEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	r := IntegrationServiceEnvironmentResource{}
	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("location").HasValue(data.Locations.Primary),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
				check.That(data.ResourceName).Key("sku_name").HasValue("Premium_0"),
				check.That(data.ResourceName).Key("access_endpoint_type").HasValue("Internal"),
				check.That(data.ResourceName).Key("virtual_network_subnet_ids.#").HasValue("4"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("development"),
				check.That(data.ResourceName).Key("connector_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("connector_outbound_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_endpoint_ip_addresses.#").Exists(),
				check.That(data.ResourceName).Key("workflow_outbound_ip_addresses.#").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func (IntegrationServiceEnvironmentResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.IntegrationServiceEnvironmentID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Logic.IntegrationServiceEnvironmentClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return nil, fmt.Errorf("retrieving Integration Service Environment %s (resource group: %s): %v", id.Name, id.ResourceGroup, err)
	}

	return utils.Bool(resp.Properties != nil), nil
}

func (IntegrationServiceEnvironmentResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-logic-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vnet-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.0.0/22"]
}

resource "azurerm_subnet" "isesubnet1" {
  name                 = "isesubnet1"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/26"]

  delegation {
    name = "integrationServiceEnvironments"
    service_delegation {
      name    = "Microsoft.Logic/integrationServiceEnvironments"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "isesubnet2" {
  name                 = "isesubnet2"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.64/26"]
}

resource "azurerm_subnet" "isesubnet3" {
  name                 = "isesubnet3"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.128/26"]
}

resource "azurerm_subnet" "isesubnet4" {
  name                 = "isesubnet4"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.192/26"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (r IntegrationServiceEnvironmentResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_integration_service_environment" "test" {
  name                 = "acctestRG-logic-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku_name             = "Premium_0"
  access_endpoint_type = "Internal"
  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r IntegrationServiceEnvironmentResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_integration_service_environment" "test" {
  name                 = "acctestRG-logic-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku_name             = "Premium_0"
  access_endpoint_type = "Internal"
  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
  tags = {
    environment = "development"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r IntegrationServiceEnvironmentResource) developer(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_integration_service_environment" "test" {
  name                 = "acctestRG-logic-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku_name             = "Developer_0"
  access_endpoint_type = "Internal"
  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
  tags = {
    environment = "development"
  }
}
`, r.template(data), data.RandomInteger)
}

func (r IntegrationServiceEnvironmentResource) skuName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_integration_service_environment" "test" {
  name                 = "acctestRG-logic-%d"
  location             = azurerm_resource_group.test.location
  resource_group_name  = azurerm_resource_group.test.name
  sku_name             = "Premium_1"
  access_endpoint_type = "Internal"
  virtual_network_subnet_ids = [
    azurerm_subnet.isesubnet1.id,
    azurerm_subnet.isesubnet2.id,
    azurerm_subnet.isesubnet3.id,
    azurerm_subnet.isesubnet4.id
  ]
}
  `, r.template(data), data.RandomInteger)
}

func (r IntegrationServiceEnvironmentResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_integration_service_environment" "import" {
  name                       = azurerm_integration_service_environment.test.name
  location                   = azurerm_resource_group.test.location
  resource_group_name        = azurerm_resource_group.test.name
  sku_name                   = azurerm_integration_service_environment.test.sku_name
  access_endpoint_type       = azurerm_integration_service_environment.test.access_endpoint_type
  virtual_network_subnet_ids = azurerm_integration_service_environment.test.virtual_network_subnet_ids
  tags                       = azurerm_integration_service_environment.test.tags
}
`, r.basic(data))
}
