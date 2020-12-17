package firewall_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMFirewall_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewall_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "ip_configuration.0.name", "configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_configuration.0.private_ip_address"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFirewall_enableDNS(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewall_enableDNS(data, "1.1.1.1", "8.8.8.8"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "dns_servers.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "dns_servers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr(data.ResourceName, "dns_servers.1", "8.8.8.8"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFirewall_withManagementIp(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewall_withManagementIp(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "management_ip_configuration.0.name", "management_configuration"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_ip_configuration.0.private_ip_address"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_ip_configuration.0.subnet_id"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "management_ip_configuration.0.public_ip_address_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFirewall_withFirewallPolicy(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewall_withFirewallPolicy(data, "policy1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "firewall_policy_id"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMFirewall_inVirtualhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_firewall", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMFirewallDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFirewall_inVirtualHub(data, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "virtual_hub.0.virtual_hub_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_hub.0.public_ip_count", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "virtual_hub.0.public_ip_addresses.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourceFirewall_basic(data acceptance.TestData) string {
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

func testAccDataSourceFirewall_enableDNS(data acceptance.TestData, dnsServers ...string) string {
	template := testAccAzureRMFirewall_enableDNS(data, dnsServers...)
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceFirewall_withManagementIp(data acceptance.TestData) string {
	template := testAccAzureRMFirewall_withManagementIp(data)
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceFirewall_withFirewallPolicy(data acceptance.TestData, policyName string) string {
	template := testAccAzureRMFirewall_withFirewallPolicy(data, policyName)
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}

func testAccDataSourceFirewall_inVirtualHub(data acceptance.TestData, pipCount int) string {
	template := testAccAzureRMFirewall_inVirtualHub(data, pipCount)
	return fmt.Sprintf(`
%s

data "azurerm_firewall" "test" {
  name                = azurerm_firewall.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template)
}
