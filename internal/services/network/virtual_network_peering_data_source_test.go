package network_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type VirtualNetworkPeeringDataSource struct{}

func TestAccDataSourceVirtualNetworkPeering_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_virtual_network_peering", "test-1")
	secondResourceName := "data.azurerm_virtual_network_peering.test-2"

	r := VirtualNetworkPeeringDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(data.ResourceName).Key("allow_forwarded_traffic").HasValue("true"),
				check.That(data.ResourceName).Key("allow_gateway_transit").HasValue("false"),
				check.That(data.ResourceName).Key("only_ipv6_peering_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("peer_complete_virtual_networks_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("use_remote_gateways").HasValue("false"),
				check.That(secondResourceName).Key("allow_virtual_network_access").HasValue("true"),
				check.That(secondResourceName).Key("allow_forwarded_traffic").HasValue("true"),
				check.That(secondResourceName).Key("allow_gateway_transit").HasValue("false"),
				check.That(secondResourceName).Key("only_ipv6_peering_enabled").HasValue("false"),
				check.That(secondResourceName).Key("peer_complete_virtual_networks_enabled").HasValue("true"),
				check.That(secondResourceName).Key("use_remote_gateways").HasValue("false"),
			),
		},
	})
}

func (r VirtualNetworkPeeringDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRg-%[1]d"
  location = "%[2]s"
}

resource "azurerm_virtual_network" "test-1" {
  name                = "acctestvnet1-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.1.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network" "test-2" {
  name                = "acctestvnet2-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  address_space       = ["10.0.2.0/24"]
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_network_peering" "test-1" {
  name                         = "peer1to2"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test-1.name
  remote_virtual_network_id    = azurerm_virtual_network.test-2.id
  allow_virtual_network_access = true
}

resource "azurerm_virtual_network_peering" "test-2" {
  name                         = "peer2to1"
  resource_group_name          = azurerm_resource_group.test.name
  virtual_network_name         = azurerm_virtual_network.test-2.name
  remote_virtual_network_id    = azurerm_virtual_network.test-1.id
  allow_virtual_network_access = true
}

data "azurerm_virtual_network_peering" "test-1" {
  name                 = "peer1to2"
  virtual_network_name = azurerm_virtual_network.test-1.name
  resource_group_name  = azurerm_resource_group.test.name
}

data "azurerm_virtual_network_peering" "test-2" {
  name                 = "peer2to1"
  virtual_network_name = azurerm_virtual_network.test-2.name
  resource_group_name  = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary)
}
