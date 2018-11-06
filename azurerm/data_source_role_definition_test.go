package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMRoleDefinition_basic(t *testing.T) {
	dataSourceName := "data.azurerm_role_definition.test"

	id := uuid.New().String()
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleDefinition(id, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
				),
			},
		},
	})
}

func testAccDataSourceRoleDefinition(id string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Created by the Data Source Role Definition Acceptance Test"

  permissions {
    actions = ["*"]

    not_actions = [
      "Microsoft.Authorization/*/Delete",
      "Microsoft.Authorization/*/Write",
      "Microsoft.Authorization/elevateAccess/Action",
    ]
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

data "azurerm_role_definition" "test" {
  role_definition_id = "${azurerm_role_definition.test.role_definition_id}"
  scope              = "${data.azurerm_subscription.primary.id}"
}
`, id, rInt)
}
