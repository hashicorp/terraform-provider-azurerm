package dns_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDNSZone_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDNSZone_tags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_tags(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDNSZone_withoutResourceGroupName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_dns_zone", "test")
	// resource group of DNS zone is always small case
	resourceGroupName := fmt.Sprintf("acctestrg-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_onlyName(data, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
				),
			},
		},
	})
}

func testAccDataSourceDNSZone_basic(data acceptance.TestData) string {
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

func testAccDataSourceDNSZone_tags(data acceptance.TestData) string {
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

func testAccDataSourceDNSZone_onlyName(data acceptance.TestData, resourceGroupName string) string {
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
