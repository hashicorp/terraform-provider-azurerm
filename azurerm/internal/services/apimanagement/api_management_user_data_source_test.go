package apimanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type ApiManagementUserDataSource struct {
}

func TestAccDataSourceApiManagementUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_user", "test")
	r := ApiManagementUserDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("user_id").HasValue("test-user"),
				check.That(data.ResourceName).Key("first_name").HasValue("Acceptance"),
				check.That(data.ResourceName).Key("last_name").HasValue("Test"),
				check.That(data.ResourceName).Key("email").HasValue(fmt.Sprintf("azure-acctest%d@example.com", data.RandomInteger)),
				check.That(data.ResourceName).Key("state").HasValue("active"),
				check.That(data.ResourceName).Key("note").HasValue("Used for testing in dimension C-137."),
			),
		},
	})
}

func (ApiManagementUserDataSource) basic(data acceptance.TestData) string {
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

resource "azurerm_api_management_user" "test" {
  user_id             = "test-user"
  api_management_name = azurerm_api_management.test.name
  resource_group_name = azurerm_resource_group.test.name
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  note                = "Used for testing in dimension C-137."
}

data "azurerm_api_management_user" "test" {
  user_id             = azurerm_api_management_user.test.user_id
  api_management_name = azurerm_api_management_user.test.api_management_name
  resource_group_name = azurerm_api_management_user.test.resource_group_name
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}
