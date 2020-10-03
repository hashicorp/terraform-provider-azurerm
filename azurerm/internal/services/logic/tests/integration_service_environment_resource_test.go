package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/logic/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMIntegrationServiceEnvironment_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIntegrationServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIntegrationServiceEnvironment_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIntegrationServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "development"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIntegrationServiceEnvironment_developer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIntegrationServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_developer(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Developer_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "development"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIntegrationServiceEnvironment_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIntegrationServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_skuName(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_1"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMIntegrationServiceEnvironment_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_integration_service_environment", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMIntegrationServiceEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMIntegrationServiceEnvironment_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMIntegrationServiceEnvironmentExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", fmt.Sprintf("acctestRG-logic-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "sku_name", "Premium_0"),
					resource.TestCheckResourceAttr(data.ResourceName, "access_endpoint_type", "Internal"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_network_subnet_ids.#", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "development"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "connector_outbound_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_endpoint_ip_addresses.#"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "workflow_outbound_ip_addresses.#"),
				),
			},
			data.RequiresImportErrorStep(testAccAzureRMIntegrationServiceEnvironment_requiresImport),
		},
	})
}

func testCheckAzureRMIntegrationServiceEnvironmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.IntegrationServiceEnvironmentClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Integration Service Environment not found: %s", resourceName)
		}
		id, err := parse.IntegrationServiceEnvironmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Integration Service Environment %q (resource group %q) does not exist", id.Name, id.ResourceGroup)
			}
			return fmt.Errorf("Bad: Get on IntegrationServiceEnvironmentClient: %+v", err)
		}
		return nil
	}
}

func testCheckAzureRMIntegrationServiceEnvironmentDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Logic.IntegrationServiceEnvironmentClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_integration_service_environment" {
			continue
		}
		id, err := parse.IntegrationServiceEnvironmentID(rs.Primary.ID)
		if err != nil {
			return err
		}
		if resp, err := client.Get(ctx, id.ResourceGroup, id.Name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on IntegrationServiceEnvironmentClient: %+v", err)
			}
		}
		return nil
	}
	return nil
}

func testAccAzureRMIntegrationServiceEnvironment_template(data acceptance.TestData) string {
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

func testAccAzureRMIntegrationServiceEnvironment_basic(data acceptance.TestData) string {
	template := testAccAzureRMIntegrationServiceEnvironment_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMIntegrationServiceEnvironment_complete(data acceptance.TestData) string {
	template := testAccAzureRMIntegrationServiceEnvironment_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMIntegrationServiceEnvironment_developer(data acceptance.TestData) string {
	template := testAccAzureRMIntegrationServiceEnvironment_template(data)
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
`, template, data.RandomInteger)
}

func testAccAzureRMIntegrationServiceEnvironment_skuName(data acceptance.TestData) string {
	template := testAccAzureRMIntegrationServiceEnvironment_template(data)
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
  `, template, data.RandomInteger)
}

func testAccAzureRMIntegrationServiceEnvironment_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMIntegrationServiceEnvironment_basic(data)
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
`, template)
}
