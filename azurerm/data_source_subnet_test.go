package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMSubnet_basic(t *testing.T) {
	resourceName := "data.azurerm_subnet.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubnet_basic(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(resourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(resourceName, "address_prefix"),
					resource.TestCheckResourceAttr(resourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttr(resourceName, "route_table_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubnet_networkSecurityGroup(t *testing.T) {
	dataSourceName := "data.azurerm_subnet.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubnet_networkSecurityGroup(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(dataSourceName, "network_security_group_id"),
					resource.TestCheckResourceAttr(dataSourceName, "route_table_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubnet_routeTable(t *testing.T) {
	dataSourceName := "data.azurerm_subnet.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubnet_routeTable(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_prefix"),
					resource.TestCheckResourceAttr(dataSourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttrSet(dataSourceName, "route_table_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubnet_serviceEndpoints(t *testing.T) {
	dataSourceName := "data.azurerm_subnet.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMSubnet_serviceEndpoints(ri, acceptance.Location()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "address_prefix"),
					resource.TestCheckResourceAttr(dataSourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttr(dataSourceName, "route_table_id", ""),
					resource.TestCheckResourceAttr(dataSourceName, "service_endpoints.#", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "service_endpoints.0", "Microsoft.Sql"),
					resource.TestCheckResourceAttr(dataSourceName, "service_endpoints.1", "Microsoft.Storage"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMSubnet_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%d-vn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctest%d-private"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.0.0/24"
}

data "azurerm_subnet" "test" {
  name                 = "${azurerm_subnet.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}
`, rInt, location, rInt, rInt)
}

func testAccDataSourceAzureRMSubnet_networkSecurityGroup(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  security_rule {
    name                       = "test123"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "*"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%d-vn"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                      = "acctest%d-private"
  resource_group_name       = "${azurerm_resource_group.test.name}"
  virtual_network_name      = "${azurerm_virtual_network.test.name}"
  address_prefix            = "10.0.0.0/24"
  network_security_group_id = "${azurerm_network_security_group.test.id}"
}

data "azurerm_subnet" "test" {
  name                 = "${azurerm_subnet.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccDataSourceAzureRMSubnet_routeTable(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_route_table" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  route {
    name                   = "acctest-%d"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_subnet" "test" {
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  route_table_id       = "${azurerm_route_table.test.id}"
}

data "azurerm_subnet" "test" {
  name                 = "${azurerm_subnet.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccDataSourceAzureRMSubnet_serviceEndpoints(rInt int, location string) string {
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
  name                 = "acctestsubnet%d"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefix       = "10.0.2.0/24"
  service_endpoints    = ["Microsoft.Sql", "Microsoft.Storage"]
}

data "azurerm_subnet" "test" {
  name                 = "${azurerm_subnet.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}
`, rInt, location, rInt, rInt)
}
