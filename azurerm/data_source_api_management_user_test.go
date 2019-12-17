package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMApiManagementUser_basic(t *testing.T) {
	dataSourceName := "data.azurerm_api_management_user.test"
	rInt := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceApiManagementUser_basic(rInt, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "user_id", "test-user"),
					resource.TestCheckResourceAttr(dataSourceName, "first_name", "Acceptance"),
					resource.TestCheckResourceAttr(dataSourceName, "last_name", "Test"),
					resource.TestCheckResourceAttr(dataSourceName, "email", fmt.Sprintf("azure-acctest%d@example.com", rInt)),
					resource.TestCheckResourceAttr(dataSourceName, "state", "active"),
					resource.TestCheckResourceAttr(dataSourceName, "note", "Used for testing in dimension C-137."),
				),
			},
		},
	})
}

func testAccDataSourceApiManagementUser_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "amtestRG-%d"
  location = "%s"
}

resource "azurerm_api_management" "test" {
  name            = "acctestAM-%d"
  publisher_name  = "pub1"
  publisher_email = "pub1@email.com"

  sku {
    name     = "Developer"
    capacity = 1
  }

  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_api_management_user" "test" {
  user_id             = "test-user"
  api_management_name = "${azurerm_api_management.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  first_name          = "Acceptance"
  last_name           = "Test"
  email               = "azure-acctest%d@example.com"
  state               = "active"
  note                = "Used for testing in dimension C-137."
}

data "azurerm_api_management_user" "test" {
  user_id             = "${azurerm_api_management_user.test.user_id}"
  api_management_name = "${azurerm_api_management_user.test.api_management_name}"
  resource_group_name = "${azurerm_api_management_user.test.resource_group_name}"
}
`, rInt, location, rInt, rInt)
}
