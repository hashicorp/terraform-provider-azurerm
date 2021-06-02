package privatedns_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type PrivateDnsZoneDatasource struct {
}

func TestAccDataSourcePrivateDNSZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccDataSourcePrivateDNSZone_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.tags(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
				check.That(data.ResourceName).Key("tags.hello").HasValue("world"),
			),
		},
	})
}

func TestAccDataSourcePrivateDNSZone_withoutResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_private_dns_zone", "test")
	r := PrivateDnsZoneDatasource{}
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.onlyNamePrep(data, resourceGroupName),
		},
		{
			Config: r.onlyName(data, resourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
			),
		},
	})
}

func (PrivateDnsZoneDatasource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PrivateDnsZoneDatasource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_private_dns_zone" "test" {
  name                = azurerm_private_dns_zone.test.name
  resource_group_name = azurerm_private_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (PrivateDnsZoneDatasource) onlyNamePrep(data acceptance.TestData, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_private_dns_zone" "test" {
  name                = "acctestzone%d.internal"
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, data.RandomInteger)
}

func (r PrivateDnsZoneDatasource) onlyName(data acceptance.TestData, resourceGroupName string) string {
	return fmt.Sprintf(`
%s

data "azurerm_private_dns_zone" "test" {
  name = azurerm_private_dns_zone.test.name
}
`, r.onlyNamePrep(data, resourceGroupName))
}
