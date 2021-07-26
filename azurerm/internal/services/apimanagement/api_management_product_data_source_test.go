package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementProductDataSource struct {
}

func TestAccDataSourceApiManagementProduct_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_product", "test")
	r := ApiManagementProductDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("product_id").HasValue("test-product"),
				check.That(data.ResourceName).Key("display_name").HasValue("Test Product"),
				check.That(data.ResourceName).Key("subscription_required").HasValue("true"),
				check.That(data.ResourceName).Key("approval_required").HasValue("true"),
				check.That(data.ResourceName).Key("published").HasValue("true"),
				check.That(data.ResourceName).Key("description").HasValue("This is an example description"),
				check.That(data.ResourceName).Key("terms").HasValue("These are some example terms and conditions"),
			),
		},
	})
}

func (ApiManagementProductDataSource) basic(data acceptance.TestData) string {
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
