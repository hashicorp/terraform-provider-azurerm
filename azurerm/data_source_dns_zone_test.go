package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMDNSZone_basic(t *testing.T) {
	dataSourceName := "data.azurerm_dns_zone.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDNSZone_tags(t *testing.T) {
	dataSourceName := "data.azurerm_dns_zone.test"
	rInt := acctest.RandInt()
	location := testLocation()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_tags(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "tags.hello", "world"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMDNSZone_withoutResourceGroupName(t *testing.T) {
	dataSourceName := "data.azurerm_dns_zone.test"
	rInt := acctest.RandInt()
	location := testLocation()
	resourceGroupName := fmt.Sprintf("acctestRG-%d", rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMDnsZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZone_onlyName(rInt, location, resourceGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
				),
			},
		},
	})
}

func testAccDataSourceDNSZone_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_dns_zone" "test" {
	name                = "acctestzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

data "azurerm_dns_zone" "test" {
	name                = "${azurerm_dns_zone.test.name}"
	resource_group_name = "${azurerm_dns_zone.test.resource_group_name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceDNSZone_tags(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "acctestRG-%d"
	location = "%s"
}

resource "azurerm_dns_zone" "test" {
	name                = "acctestzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
	tags {
		hello = "world"
	}
}

data "azurerm_dns_zone" "test" {
	name                = "${azurerm_dns_zone.test.name}"
	resource_group_name = "${azurerm_dns_zone.test.resource_group_name}"
}
`, rInt, location, rInt)
}

func testAccDataSourceDNSZone_onlyName(rInt int, location, resourceGroupName string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
	name     = "%s"
	location = "%s"
}

resource "azurerm_dns_zone" "test" {
	name                = "acctestzone%d.com"
	resource_group_name = "${azurerm_resource_group.test.name}"
}

data "azurerm_dns_zone" "test" {
	name = "${azurerm_dns_zone.test.name}"
}
`, resourceGroupName, location, rInt)
}
