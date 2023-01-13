package privatednsresolver_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type PrivateDNSResolverDnsResolverDataSource struct{}

func TestAccPrivateDNSResolverDnsResolverDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver", "test")
	d := PrivateDNSResolverDnsResolverDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("virtual_network_id").Exists(),
			),
		},
	})
}

func TestAccPrivateDNSResolverDnsResolverDataSource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_resolver", "test")
	d := PrivateDNSResolverDnsResolverDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: d.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("resource_group_name").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.Locations.Primary)),
				check.That(data.ResourceName).Key("virtual_network_id").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("value"),
			),
		},
	})
}

func (d PrivateDNSResolverDnsResolverDataSource) template(data acceptance.TestData) string {
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
`, data.Locations.Primary, data.RandomInteger)
}

func (d PrivateDNSResolverDnsResolverDataSource) basic(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
}

data "azurerm_private_dns_resolver" "test" {
  name                = azurerm_private_dns_resolver.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}

func (d PrivateDNSResolverDnsResolverDataSource) complete(data acceptance.TestData) string {
	template := d.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_private_dns_resolver" "test" {
  name                = "acctest-dr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_network_id  = azurerm_virtual_network.test.id
  tags = {
    key = "value"
  }
}

data "azurerm_private_dns_resolver" "test" {
  name                = azurerm_private_dns_resolver.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, template, data.RandomInteger)
}
