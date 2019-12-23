package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceArmVirtualNetworkInterface_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_interface", "test")

	name := fmt.Sprintf("acctest-nic-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArmVirtualNetworkInterface_basic(data),
			},
			{
				Config: testAccDataSourceArmVirtualNetworkInterface_withDataSource(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttrSet(data.ResourceName, "id"),
					resource.TestCheckResourceAttr(data.ResourceName, "private_ip_address", "10.0.1.4"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
				),
			},
		},
	})
}

func testAccDataSourceArmVirtualNetworkInterface_basic(data acceptance.TestData) string {
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
  name                = "acctest-nsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testAccDataSourceArmVirtualNetworkInterface_withDataSource(data acceptance.TestData) string {
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
  name                = "acctest-nsg-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
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
    private_ip_address_allocation = "Dynamic"
  }

  tags = {
    environment = "staging"
  }
}

data "azurerm_network_interface" "test" {
  name                = "acctest-nic-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}
