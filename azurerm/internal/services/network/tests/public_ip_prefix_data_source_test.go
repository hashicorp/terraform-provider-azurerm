package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPublicIPPrefix_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_public_ip_prefix", "test")
	name := fmt.Sprintf("acctestpublicipprefix-%d", data.RandomInteger)
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPublicIPPrefixBasic(name, resourceGroupName, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", name),
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestCheckResourceAttr(data.ResourceName, "location", data.Locations.Primary),
					resource.TestCheckResourceAttr(data.ResourceName, "sku", "Standard"),
					resource.TestCheckResourceAttr(data.ResourceName, "prefix_length", "31"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "ip_prefix"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "tags.env", "test"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPublicIPPrefixBasic(name string, resourceGroupName string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"

  tags = {
    env = "test"
  }
}

resource "azurerm_public_ip_prefix" "test" {
  name                = "%s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Standard"
  prefix_length       = 31

  tags = {
    env = "test"
  }
}

data "azurerm_public_ip_prefix" "test" {
  name                = azurerm_public_ip_prefix.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, resourceGroupName, data.Locations.Primary, name)
}
