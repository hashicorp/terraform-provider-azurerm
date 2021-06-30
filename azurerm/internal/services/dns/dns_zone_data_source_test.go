package dns_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type AzureRMDNSZoneDataSource struct {
}

func TestAccAzureRMDNSZoneDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")
	r := AzureRMDNSZoneDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("tags.%").HasValue("0"),
			),
		},
	})
}

func TestAccAzureRMDNSZoneDataSource_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")
	r := AzureRMDNSZoneDataSource{}

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

func TestAccAzureRMDNSZoneDataSource_withoutResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")
	r := AzureRMDNSZoneDataSource{}
	// resource group of DNS zone is always small case
	resourceGroupName := fmt.Sprintf("acctestrg-%d", data.RandomInteger)

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.onlyName(data, resourceGroupName),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_group_name").HasValue(resourceGroupName),
			),
		},
	})
}

func (AzureRMDNSZoneDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_dns_zone" "test" {
  name                = azurerm_dns_zone.test.name
  resource_group_name = azurerm_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AzureRMDNSZoneDataSource) tags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name

  tags = {
    hello = "world"
  }
}

data "azurerm_dns_zone" "test" {
  name                = azurerm_dns_zone.test.name
  resource_group_name = azurerm_dns_zone.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func (AzureRMDNSZoneDataSource) onlyName(data acceptance.TestData, resourceGroupName string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_dns_zone" "test" {
  name                = "acctestzone%d.com"
  resource_group_name = azurerm_resource_group.test.name
}

data "azurerm_dns_zone" "test" {
  name = azurerm_dns_zone.test.name
}
`, resourceGroupName, data.Locations.Primary, data.RandomInteger)
}
