package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMCdnProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCdnProfile_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists("data.azurerm_cdn_profile.test"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMCdnProfile_withTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_cdn_profile", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMCdnProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMCdnProfile_withTags(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMCdnProfileExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.environment", "Production"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.cost_center", "MSFT"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMCdnProfile_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"
}

data "azurerm_cdn_profile" "test" {
  name                = azurerm_cdn_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}

func testAccDataSourceAzureRMCdnProfile_withTags(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_cdn_profile" "test" {
  name                = "acctestcdnprof%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard_Verizon"

  tags = {
    environment = "Production"
    cost_center = "MSFT"
  }
}

data "azurerm_cdn_profile" "test" {
  name                = azurerm_cdn_profile.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
