package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMDNSZone_basic(t *testing.T) {
	dataSourceName := "data.azurerm_dns_zone.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()
	resourceGroupName := fmt.Sprintf("acctestRG-%d", rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

  tags = {
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
