package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccAzureRMAzureFirewall_basic(t *testing.T) {
	resourceName := "azurerm_azure_firewall.test"
	ri := acctest.RandInt()
	config := testAccAzureRMAzureFirewall_basic(ri, testLocation())
}

func testAccAzureRMAzureFirewall_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name                      = "AzureFirewallSubnet"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.1.0/24"
}
resource "azurerm_public_ip" "test" {
  name                         = "acctestpip%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  public_ip_address_allocation = "Static"
  sku                          = "Standard"
}
resource "azurerm_azure_firewall" "test" {
  name = "acctestfirewall%d"
  location                     = "${azurerm_resource_group.test.location}"
  resource_group_name          = "${azurerm_resource_group.test.name}"
  ip_configuration {
    name                          = "configuration"
    subnet_id                     = "${azurerm_subnet.test.id}"
    internal_public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}
`, rInt, location, rInt, rInt, rInt)
}
