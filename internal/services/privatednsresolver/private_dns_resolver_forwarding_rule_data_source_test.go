package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverForwardingRuleDataSource struct{}

func TestAccPrivateDNSResolverForwardingRuleDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_forwarding_rule", "test")
	d := PrivateDNSResolverForwardingRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dns_forwarding_ruleset_id").Exists(),
				check.That(data.ResourceName).Key("domain_name").HasValue("onprem.local."),
				check.That(data.ResourceName).Key("target_dns_servers.0.ip_address").HasValue("10.10.0.1"),
				check.That(data.ResourceName).Key("target_dns_servers.0.port").HasValue("53"),
			),
		},
	})
}

func TestAccPrivateDNSResolverForwardingRuleDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver_forwarding_rule", "test")
	d := PrivateDNSResolverForwardingRuleDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("dns_forwarding_ruleset_id").Exists(),
				check.That(data.ResourceName).Key("domain_name").HasValue("onprem.local."),
				check.That(data.ResourceName).Key("enabled").HasValue("true"),
				check.That(data.ResourceName).Key("target_dns_servers.0.ip_address").HasValue("10.10.0.1"),
				check.That(data.ResourceName).Key("target_dns_servers.0.port").HasValue("53"),
				check.That(data.ResourceName).Key("metadata.key").HasValue("value"),
			),
		},
	})
}

func (d PrivateDNSResolverForwardingRuleDataSource) template(data acceptance.TestData) string {
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

resource "azurerm_private_dns_resolver_dns_forwarding_ruleset" "test" {
  name                                       = "acctest-drdfr-%[2]d"
  resource_group_name                        = azurerm_resource_group.test.name
  location                                   = azurerm_resource_group.test.location
  private_dns_resolver_outbound_endpoint_ids = [azurerm_private_dns_resolver_outbound_endpoint.test.id]
}
`, data.Locations.Primary, data.RandomInteger)
}

func (d PrivateDNSResolverForwardingRuleDataSource) basic(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = "acctest-drfr-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem.local."
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
}

data "azurerm_private_dns_resolver_forwarding_rule" "test" {
	name         		  		= azurerm_private_dns_resolver_forwarding_rule.test.name
	dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
}
`, template, data.RandomInteger)
}

func (d PrivateDNSResolverForwardingRuleDataSource) complete(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_private_dns_resolver_forwarding_rule" "test" {
  name                      = "acctest-drfr-%d"
  dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
  domain_name               = "onprem.local."
  enabled                   = true
  target_dns_servers {
    ip_address = "10.10.0.1"
    port       = 53
  }
  metadata = {
    key = "value"
  }
}

data "azurerm_private_dns_resolver_forwarding_rule" "test" {
	name         		  		= azurerm_private_dns_resolver_forwarding_rule.test.name
	dns_forwarding_ruleset_id = azurerm_private_dns_resolver_dns_forwarding_ruleset.test.id
}
`, template, data.RandomInteger)
}
