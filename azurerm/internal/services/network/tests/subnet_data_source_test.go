package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceSubnet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubnet_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "route_table_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMSubnet_basic_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubnet_basic_addressPrefixes(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefixes.#"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "route_table_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceSubnet_networkSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				// since the network security group association is a separate resource this forces it
				Config: testAccDataSourceSubnet_networkSecurityGroupDependencies(data),
			},
			{
				Config: testAccDataSourceSubnet_networkSecurityGroup(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "network_security_group_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "route_table_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceSubnet_routeTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				// since the route table association is a separate resource this forces it
				Config: testAccDataSourceSubnet_routeTableDependencies(data),
			},
			{
				Config: testAccDataSourceSubnet_routeTable(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefix"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "route_table_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_security_group_id", ""),
				),
			},
		},
	})
}

func TestAccDataSourceSubnet_serviceEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSubnet_serviceEndpoint(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "resource_group_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_network_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "address_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "network_security_group_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "route_table_id", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "service_endpoints.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_endpoints.0", "Microsoft.Sql"),
					resource.TestCheckResourceAttr(data.ResourceName, "service_endpoints.1", "Microsoft.Storage"),
				),
			},
		},
	})
}

func testAccDataSourceSubnet_basic(data acceptance.TestData) string {
	template := testAccDataSourceSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, template)
}

func testAccDataSourceSubnet_basic_addressPrefixes(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}
resource "azurerm_resource_group" "test" {
  name     = "acctest%d-rg"
  location = "%s"
}
resource "azurerm_virtual_network" "test" {
  name                = "acctest%d-vn"
  address_space       = ["10.0.0.0/16", "ace:cab:deca::/48"]
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}
resource "azurerm_subnet" "test" {
  name                 = "acctest%d-private"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
  address_prefixes     = ["10.0.0.0/24", "ace:cab:deca:deed::/64"]
}
data "azurerm_subnet" "test" {
  name                 = "${azurerm_subnet.test.name}"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  virtual_network_name = "${azurerm_virtual_network.test.name}"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func testAccDataSourceSubnet_networkSecurityGroupDependencies(data acceptance.TestData) string {
	template := testAccDataSourceSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

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

resource "azurerm_subnet_network_security_group_association" "test" {
  subnet_id                 = azurerm_subnet.test.id
  network_security_group_id = azurerm_network_security_group.test.id
}
`, template, data.RandomInteger)
}

func testAccDataSourceSubnet_networkSecurityGroup(data acceptance.TestData) string {
	template := testAccDataSourceSubnet_networkSecurityGroupDependencies(data)
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, template)
}

func testAccDataSourceSubnet_routeTableDependencies(data acceptance.TestData) string {
	template := testAccDataSourceSubnet_template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_subnet" "test" {
  name                 = "internal"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.0.0/24"
}

resource "azurerm_route_table" "test" {
  name                = "acctestrt-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  route {
    name                   = "first"
    address_prefix         = "10.100.0.0/14"
    next_hop_type          = "VirtualAppliance"
    next_hop_in_ip_address = "10.10.1.1"
  }
}

resource "azurerm_subnet_route_table_association" "test" {
  subnet_id      = azurerm_subnet.test.id
  route_table_id = azurerm_route_table.test.id
}
`, template, data.RandomInteger)
}

func testAccDataSourceSubnet_routeTable(data acceptance.TestData) string {
	template := testAccDataSourceSubnet_routeTableDependencies(data)
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, template)
}

func testAccDataSourceSubnet_serviceEndpoint(data acceptance.TestData) string {
	template := testAccAzureRMSubnet_serviceEndpointsUpdated(data)
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, template)
}

func testAccDataSourceSubnet_template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest%d-vn"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
