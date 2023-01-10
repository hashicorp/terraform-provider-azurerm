package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverDnsForwardingRulesetDataSource struct{}

func TestAccPrivateDNSResolverDnsForwardingRulesetDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	d := PrivateDNSResolverDnsForwardingRulesetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("private_dns_resolver_outbound_endpoint_ids.0").Exists(),
			),
		},
	})
}

func TestAccPrivateDNSResolverDnsForwardingRulesetDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_dns_forwarding_ruleset", "test")
	d := PrivateDNSResolverDnsForwardingRulesetDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("private_dns_resolver_outbound_endpoint_ids.0").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("value"),
			),
		},
	})
}

func (d PrivateDNSResolverDnsForwardingRulesetDataSource) template(data acceptance.TestData) string {
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
  name                 = "outbounddns"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.0.64/28"]

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

resource "azurerm_private_dns_resolver_outbound_endpoint" "test" {
  name                    = "acctest-droe-%[2]d"
  private_dns_resolver_id = azurerm_private_dns_resolver.test.id
  location                = azurerm_private_dns_resolver.test.location
  subnet_id               = azurerm_subnet.test.id
}
`, data.Locations.Primary, data.RandomInteger)
}

func (d PrivateDNSResolverDnsForwardingRulesetDataSource) basic(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
}

data "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
	name         		  = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.name
	resource_group_name = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}

func (d PrivateDNSResolverDnsForwardingRulesetDataSource) complete(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
  tags = {
    key = "value"
  }
}

data "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
	name         		  = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.name
	resource_group_name = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}
