package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementGroupDataSource struct {
}

func TestAccDataSourceApiManagementGroup_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_group", "test")
	r := ApiManagementGroupDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("display_name").HasValue("Test Group"),
				check.That(data.ResourceName).Key("description").HasValue(""),
				check.That(data.ResourceName).Key("external_id").HasValue(""),
				check.That(data.ResourceName).Key("type").HasValue("custom"),
			),
		},
	})
}

func (ApiManagementGroupDataSource) basic(data acceptance.TestData) string {
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
  sku_name            = "Developer_1"
}

resource "azurerm_api_management_group" "test" {
  name                = "acctestAMGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  api_management_name = azurerm_api_management.test.name
  display_name        = "Test Group"
}

data "azurerm_api_management_group" "test" {
  name                = azurerm_api_management_group.test.name
  api_management_name = azurerm_api_management_group.test.api_management_name
  resource_group_name = azurerm_api_management_group.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
