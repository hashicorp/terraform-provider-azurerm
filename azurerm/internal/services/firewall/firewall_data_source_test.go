package firewall_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type FirewallDataSource struct {
}

func TestAccFirewallDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")
	r := FirewallDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
			),
		},
	})
}

func TestAccFirewallDataSource_enableDNS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")
	r := FirewallDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.enableDNS(data, "1.1.1.1", "8.8.8.8"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("dns_servers.#").HasValue("2"),
				check.That(data.ResourceName).Key("dns_servers.0").HasValue("1.1.1.1"),
				check.That(data.ResourceName).Key("dns_servers.1").HasValue("8.8.8.8"),
			),
		},
	})
}

func TestAccFirewallDataSource_withManagementIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")
	r := FirewallDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withManagementIp(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("ip_configuration.0.name").HasValue("configuration"),
				check.That(data.ResourceName).Key("ip_configuration.0.private_ip_address").Exists(),
				check.That(data.ResourceName).Key("management_ip_configuration.0.name").HasValue("management_configuration"),
				check.That(data.ResourceName).Key("management_ip_configuration.0.public_ip_address_id").Exists(),
			),
		},
	})
}

func TestAccFirewallDataSource_withFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")
	r := FirewallDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.withFirewallPolicy(data, "policy1"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("firewall_policy_id").Exists(),
			),
		},
	})
}

func TestAccFirewallDataSource_inVirtualhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")
	r := FirewallDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.inVirtualHub(data, 2),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("virtual_hub.0.virtual_hub_id").Exists(),
				check.That(data.ResourceName).Key("virtual_hub.0.public_ip_count").HasValue("2"),
				check.That(data.ResourceName).Key("virtual_hub.0.public_ip_addresses.#").HasValue("2"),
			),
		},
	})
}

func (FirewallDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctestvirtnet%d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "AzureFirewallSubnet"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefix       = "10.0.1.0/24"
}

resource "azurerm_public_ip" "test" {
  name                = "acctestpip%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_firewall" "test" {
  name                = "acctestfirewall%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  ip_configuration {
    name                 = "configuration"
    subnet_id            = azurerm_subnet.test.id
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func (FirewallDataSource) enableDNS(data acceptance.TestData, dnsServers ...string) string {
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, FirewallResource{}.enableDNS(data, dnsServers...))
}

func (FirewallDataSource) withManagementIp(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, FirewallResource{}.withManagementIp(data))
}

func (FirewallDataSource) withFirewallPolicy(data acceptance.TestData, policyName string) string {
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, FirewallResource{}.withFirewallPolicy(data, policyName))
}

func (FirewallDataSource) inVirtualHub(data acceptance.TestData, pipCount int) string {
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, FirewallResource{}.inVirtualHub(data, pipCount))
}
