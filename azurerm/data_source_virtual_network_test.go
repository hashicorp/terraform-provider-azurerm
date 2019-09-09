package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceArmVirtualNetwork_basic(t *testing.T) {
	dataSourceName := "data.azurerm_virtual_network.test"
	ri := tf.AccRandTimeInt()

	name := fmt.Sprintf("acctestvnet-%d", ri)
	location := testLocation()
	config := testAccDataSourceArmVirtualNetwork_basic(ri, testLocation())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "location", azure.NormalizeLocation(location)),
					resource.TestCheckResourceAttr(dataSourceName, "dns_servers.0", "10.0.0.4"),
					resource.TestCheckResourceAttr(dataSourceName, "address_spaces.0", "10.0.0.0/16"),
					resource.TestCheckResourceAttr(dataSourceName, "subnets.0", "subnet1"),
				),
			},
		},
	})
}

func TestAccDataSourceArmVirtualNetwork_peering(t *testing.T) {
	dataSourceName := "data.azurerm_virtual_network.test"
	ri := tf.AccRandTimeInt()

	virtualNetworkName := fmt.Sprintf("acctestvnet-1-%d", ri)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmVirtualNetwork_peering(ri, location),
			},
			{
				Config: testAccDataSourceArmVirtualNetwork_peeringWithDataSource(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", virtualNetworkName),
					resource.TestCheckResourceAttr(dataSourceName, "address_spaces.0", "10.0.1.0/24"),
					resource.TestCheckResourceAttr(dataSourceName, "vnet_peerings.%", "1"),
				),
			},
		},
	})
}

func testAccDataSourceArmVirtualNetwork_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvnet-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  dns_servers         = ["10.0.0.4"]

  subnet {
    name           = "subnet1"
    address_prefix = "10.0.1.0/24"
  }
}

data "azurerm_virtual_network" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_virtual_network.test.name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceArmVirtualNetwork_peering(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet-1-%d"
  address_space       = ["10.0.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet-2-%d"
  address_space       = ["10.0.2.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network_peering" "test1" {
  name                      = "peer-1to2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test1.name}"
  remote_virtual_network_id = "${azurerm_virtual_network.test2.id}"
}
`, rInt, location, rInt, rInt)
}

func testAccDataSourceArmVirtualNetwork_peeringWithDataSource(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test1" {
  name                = "acctestvnet-1-%d"
  address_space       = ["10.0.1.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network" "test2" {
  name                = "acctestvnet-2-%d"
  address_space       = ["10.0.2.0/24"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_virtual_network_peering" "test1" {
  name                      = "peer-1to2"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test1.name}"
  remote_virtual_network_id = "${azurerm_virtual_network.test2.id}"
}

data "azurerm_virtual_network" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_virtual_network.test1.name}"
}
`, rInt, location, rInt, rInt)
}
