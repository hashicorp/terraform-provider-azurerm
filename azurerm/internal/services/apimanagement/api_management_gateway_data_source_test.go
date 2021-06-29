package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementGatewayDataSource struct {
}

func TestAccDataSourceApiManagementGateway_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_gateway", "test")
	r := ApiManagementGatewayDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").HasValue("old world"),
				check.That(data.ResourceName).Key("description").HasValue("this is a test gateway"),
			),
		},
	})
}

func (ApiManagementGatewayDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name                = "acctestAM-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  publisher_name      = "pub1"
  publisher_email     = "pub1@email.com"

  sku_name = "Developer_1"
}

resource "azurerm_api_management_gateway" "test" {
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  gateway_id          = "TestGateway"
  location            = "old world updated"
  description         = "this is a test gateway updated"
}

data "azurerm_api_management_gateway" "test" {
  gateway_id          = azurerm_api_management_gateway.test.gateway_id
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
