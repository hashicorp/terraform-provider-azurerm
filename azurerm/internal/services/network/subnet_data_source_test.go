package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type SubnetDataSource struct {
}

func TestAccDataSourceSubnet_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")
	r := SubnetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_name").Exists(),
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").HasValue(""),
				check.That(data.ResourceName).Key("route_table_id").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceAzureRMSubnet_basic_addressPrefixes(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")
	r := SubnetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic_addressPrefixes(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_name").Exists(),
				check.That(data.ResourceName).Key("address_prefixes.#").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").HasValue(""),
				check.That(data.ResourceName).Key("route_table_id").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceSubnet_networkSecurityGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")
	r := SubnetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// since the network security group association is a separate resource this forces it
			Config: r.networkSecurityGroupDependencies(data),
		},
		{
			Config: r.networkSecurityGroup(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_name").Exists(),
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").Exists(),
				check.That(data.ResourceName).Key("route_table_id").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceSubnet_routeTable(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")
	r := SubnetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			// since the route table association is a separate resource this forces it
			Config: r.routeTableDependencies(data),
		},
		{
			Config: r.routeTable(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_name").Exists(),
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("route_table_id").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").HasValue(""),
			),
		},
	})
}

func TestAccDataSourceSubnet_serviceEndpoints(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_subnet", "test")
	r := SubnetDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.serviceEndpoint(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("virtual_network_name").Exists(),
				check.That(data.ResourceName).Key("address_prefix").Exists(),
				check.That(data.ResourceName).Key("network_security_group_id").HasValue(""),
				check.That(data.ResourceName).Key("route_table_id").HasValue(""),
				check.That(data.ResourceName).Key("service_endpoints.#").HasValue("2"),
				check.That(data.ResourceName).Key("service_endpoints.0").HasValue("Microsoft.Sql"),
				check.That(data.ResourceName).Key("service_endpoints.1").HasValue("Microsoft.Storage"),
			),
		},
	})
}

func (r SubnetDataSource) basic(data acceptance.TestData) string {
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
`, r.template(data))
}

func (SubnetDataSource) basic_addressPrefixes(data acceptance.TestData) string {
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

func (r SubnetDataSource) networkSecurityGroupDependencies(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r SubnetDataSource) networkSecurityGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, r.networkSecurityGroupDependencies(data))
}

func (r SubnetDataSource) routeTableDependencies(data acceptance.TestData) string {
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
`, r.template(data), data.RandomInteger)
}

func (r SubnetDataSource) routeTable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, r.routeTableDependencies(data))
}

func (SubnetDataSource) serviceEndpoint(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_subnet" "test" {
  name                 = azurerm_subnet.test.name
  virtual_network_name = azurerm_subnet.test.virtual_network_name
  resource_group_name  = azurerm_subnet.test.resource_group_name
}
`, SubnetResource{}.serviceEndpointsUpdated(data))
}

func (SubnetDataSource) template(data acceptance.TestData) string {
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
