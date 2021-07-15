package network_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type LocalNetworkGatewayDataSource struct {
}

func TestAccDataSourceLocalNetworkGateway_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_local_network_gateway", "test")
	r := LocalNetworkGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("gateway_address").HasValue("127.0.0.1"),
				check.That(data.ResourceName).Key("address_space.0").HasValue("127.0.1.0/24"),
				check.That(data.ResourceName).Key("address_space.1").HasValue("127.0.0.0/24"),
				check.That(data.ResourceName).Key("bgp_settings.#").HasValue("1"),
				check.That(data.ResourceName).Key("bgp_settings.0.asn").HasValue("2468"),
				check.That(data.ResourceName).Key("bgp_settings.0.bgp_peering_address").HasValue("10.104.1.1"),
				check.That(data.ResourceName).Key("bgp_settings.0.peer_weight").HasValue("15"),
			),
		},
	})
}

func (LocalNetworkGatewayDataSource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-lngw-%d"
  location = "%s"
}

resource "azurerm_local_network_gateway" "test" {
  name                = "acctestlng-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  gateway_address     = "127.0.0.1"
  address_space       = ["127.0.1.0/24", "127.0.0.0/24"]

  bgp_settings {
    asn                 = 2468
    bgp_peering_address = "10.104.1.1"
    peer_weight         = 15
  }

  tags = {
    environment = "acctest"
  }
}

data "azurerm_local_network_gateway" "test" {
  name                = azurerm_local_network_gateway.test.name
  resource_group_name = azurerm_local_network_gateway.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
