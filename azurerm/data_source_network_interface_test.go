package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceArmVirtualNetworkInterface_basic(t *testing.T) {
	dataSourceName := "data.azurerm_network_interface.test"
	ri := acctest.RandInt()

	name := fmt.Sprintf("acctest-nic-%d", ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmVirtualNetworkInterface_basic(ri, testLocation()),
			},
			{
				Config: testAccDataSourceArmVirtualNetworkInterface_withDataSource(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "id"),
					resource.TestCheckResourceAttr(dataSourceName, "private_ip_address", "10.0.1.4"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_security_group_id"),
				),
			},
		},
	})
}

func testAccDataSourceArmVirtualNetworkInterface_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_security_group" "test" { 
  name                 = "acctest-nsg-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_network_interface" "test" {
  name                      = "acctest-nic-%d"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }

  tags {
    environment = "staging"
  }
}

`, rInt, location, rInt, rInt, rInt)
}

func testAccDataSourceArmVirtualNetworkInterface_withDataSource(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest-%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-vn-%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_network_security_group" "test" { 
  name                 = "acctest-nsg-%d"
  location             = "${azurerm_resource_group.test.location}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "subnet1"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_network_interface" "test" {
  name                      = "acctest-nic-%d"
  location                  = "${azurerm_resource_group.test.location}"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  network_security_group_id = "${azurerm_network_security_group.test.id}"

  ip_configuration {
    name                          = "testconfiguration1"
    subnet_id                     = "${azurerm_subnet.test.id}"
    private_ip_address_allocation = "dynamic"
  }

  tags {
    environment = "staging"
  }
}

data "azurerm_network_interface" "test" {
  name = "acctest-nic-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
}


`, rInt, location, rInt, rInt, rInt, rInt)
}
