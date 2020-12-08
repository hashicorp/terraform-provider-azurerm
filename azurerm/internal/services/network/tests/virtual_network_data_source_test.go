package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmVirtualNetwork_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network", "test")

	name := fmt.Sprintf("acctestvnet-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmVirtualNetwork_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "location", azure.NormalizeLocation(data.Locations.Primary)),
					resource.TestCheckResourceAttr(data.ResourceName, "dns_servers.0", "10.0.0.4"),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(data.ResourceName, "subnets.0", "subnet1"),
				),
			},
		},
	})
}

func TestAccDataSourceArmVirtualNetwork_peering(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network", "test")

	virtualNetworkName := fmt.Sprintf("acctestvnet-1-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmVirtualNetwork_peering(data),
			},
			{
				Config: testAccDataSourceArmVirtualNetwork_peeringWithDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", virtualNetworkName),
					resource.TestCheckResourceAttr(data.ResourceName, "address_space.0", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(data.ResourceName, "vnet_peerings.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceArmVirtualNetwork_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  dns_servers         = ["10.0.0.4"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

data "azurerm_virtual_network" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_virtual_network.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceArmVirtualNetwork_peering(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet-1-%d"
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet-2-%d"
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network_peering" "test1" {
  name                      = "peer-1to2"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.test1.name
  remote_virtual_network_id = azurerm_virtual_network.test2.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccDataSourceArmVirtualNetwork_peeringWithDataSource(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet-1-%d"
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet-2-%d"
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_virtual_network_peering" "test1" {
  name                      = "peer-1to2"
  resource_group_name       = azurerm_resource_group.test.name
  virtual_network_name      = azurerm_virtual_network.test1.name
  remote_virtual_network_id = azurerm_virtual_network.test2.id
}

data "azurerm_virtual_network" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_virtual_network.test1.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
