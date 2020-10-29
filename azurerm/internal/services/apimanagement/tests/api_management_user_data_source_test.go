package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementUser_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_api_management_user", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementUser_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "user_id", "test-user"),
					resource.TestCheckResourceAttr(data.ResourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(data.ResourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(data.ResourceName, "email", fmt.Sprintf("azure-acctest%d@example.com", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "state", "active"),
					resource.TestCheckResourceAttr(data.ResourceName, "note", "Used for testing in dimension C-137."),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementUser_basic(data acceptance.TestData) string {
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
