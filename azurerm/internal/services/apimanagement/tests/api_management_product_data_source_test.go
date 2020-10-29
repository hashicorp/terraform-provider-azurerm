package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementProduct_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_product", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementProduct_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "product_id", "test-product"),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", "Test Product"),
					resource.TestCheckResourceAttr(data.ResourceName, "subscription_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "approval_required", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "published", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is an example description"),
					resource.TestCheckResourceAttr(data.ResourceName, "terms", "These are some example terms and conditions"),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementProduct_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"
  sku_name            = "Developer_1"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_api_management_product" "test" {
  product_id            = "test-product"
  api_management_name   = azurerm_api_management.test.name
  resource_group_name   = azurerm_resource_group.test.name
  display_name          = "Test Product"
  subscription_required = true
  approval_required     = true
  subscriptions_limit   = 2
  published             = true
  description           = "This is an example description"
  terms                 = "These are some example terms and conditions"
}

data "azurerm_api_management_product" "test" {
  product_id          = azurerm_api_management_product.test.product_id
  api_management_name = azurerm_api_management_product.test.api_management_name
  resource_group_name = azurerm_api_management_product.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
