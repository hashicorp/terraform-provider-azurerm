package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type DNSResolverInboundEndpointDataSource struct{}

func TestAccDNSResolverInboundEndpointDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_inbound_endpoint", "test")
	d := DNSResolverInboundEndpointDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("private_dns_resolver_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("ip_configurations.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("ip_configurations.0.private_ip_address").Exists(),
			),
		},
	})
}

func TestAccDNSResolverInboundEndpointDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_inbound_endpoint", "test")
	d := DNSResolverInboundEndpointDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("private_dns_resolver_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("ip_configurations.0.private_ip_allocation_method").HasValue("Dynamic"),
				check.That(data.ResourceName).Key("ip_configurations.0.subnet_id").Exists(),
				check.That(data.ResourceName).Key("ip_configurations.0.private_ip_address").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("value"),
			),
		},
	})
}

func (d DNSResolverInboundEndpointDataSource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}


resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%[2]d"
  location = "%[1]s"
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest-rg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "test" {
  name                 = "inbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.0/28"]

  delegation {
    name = "Microsoft.Network.dnsResolvers"
    service_delegation {
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
      name    = "Microsoft.Network/dnsResolvers"
    }
  }
}

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (d DNSResolverInboundEndpointDataSource) basic(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  ip_configurations {
    subnet_id = azurerm_subnet.test.id
  }
}

data "azurerm_private_dns_resolver_inbound_endpoint" "test" {
	name         		    = azurerm_private_dns_resolver_inbound_endpoint.test.name
	private_dns_resolver_id = azurerm_private_dns_resolver.test.id
}
`, template, data.RandomInteger)
}

func (d DNSResolverInboundEndpointDataSource) complete(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_inbound_endpoint" "test" {
  name                    = "acctest-drie-%d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  ip_configurations {
    private_ip_allocation_method = "Dynamic"
    subnet_id                    = azurerm_subnet.test.id
  }
  tags = {
    key = "value"
  }
}

data "azurerm_private_dns_resolver_inbound_endpoint" "test" {
	name         		    = azurerm_private_dns_resolver_inbound_endpoint.test.name
	private_dns_resolver_id = azurerm_private_dns_resolver.test.id
}
`, template, data.RandomInteger)
}
