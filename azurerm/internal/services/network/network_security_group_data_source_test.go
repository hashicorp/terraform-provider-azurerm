package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type NetworkSecurityGroupDataSource struct {
}

func TestAccDataSourceNetworkSecurityGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")
	r := NetworkSecurityGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceNetworkSecurityGroup_rules(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")
	r := NetworkSecurityGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withRules(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("1"),
				check.That(data.ResourceName).Key("security_rule.0.name").HasValue("test123"),
				check.That(data.ResourceName).Key("security_rule.0.priority").HasValue("100"),
				check.That(data.ResourceName).Key("security_rule.0.direction").HasValue("Inbound"),
				check.That(data.ResourceName).Key("security_rule.0.access").HasValue("Allow"),
				check.That(data.ResourceName).Key("security_rule.0.protocol").HasValue("Tcp"),
				check.That(data.ResourceName).Key("security_rule.0.source_port_range").HasValue("*"),
				check.That(data.ResourceName).Key("security_rule.0.destination_port_range").HasValue("*"),
				check.That(data.ResourceName).Key("security_rule.0.source_address_prefix").HasValue("*"),
				check.That(data.ResourceName).Key("security_rule.0.destination_address_prefix").HasValue("*"),
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourceNetworkSecurityGroup_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_network_security_group", "test")
	r := NetworkSecurityGroupDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.tags(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
				check.That(data.ResourceName).Key("security_rule.#").HasValue("0"),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.environment").HasValue("staging"),
			),
		},
	})
}

func (NetworkSecurityGroupDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_network_security_group" "test" {
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (NetworkSecurityGroupDataSource) withRules(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
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

data "azurerm_network_security_group" "test" {
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (NetworkSecurityGroupDataSource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_network_security_group" "test" {
  name                = "acctestnsg-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    environment = "staging"
  }
}

data "azurerm_network_security_group" "test" {
  name                = azurerm_network_security_group.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
