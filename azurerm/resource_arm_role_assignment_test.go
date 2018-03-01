package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoleAssignment_emptyName(t *testing.T) {
	resourceName := "azurerm_role_assignment.test"
	config := testAccAzureRMRoleAssignment_emptyName()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleAssignment_roleName(t *testing.T) {
	id := uuid.New().String()
	resourceName := "azurerm_role_assignment.test"
	config := testAccAzureRMRoleAssignment_roleName(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "role_definition_id"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleAssignment_builtin(t *testing.T) {
	id := uuid.New().String()
	config := testAccAzureRMRoleAssignment_builtin(id)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleAssignment_custom(t *testing.T) {
	roleDefinitionId := uuid.New().String()
	roleAssignmentId := uuid.New().String()
	rInt := acctest.RandInt()
	config := testAccAzureRMRoleAssignment_custom(roleDefinitionId, roleAssignmentId, rInt)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleAssignmentExists("azurerm_role_assignment.test"),
				),
			},
		},
	})
}

func testCheckAzureRMRoleAssignmentExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).roleAssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleAssignmentName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Assignment %q (Scope: %q) does not exist", name, scope)
			}
			return fmt.Errorf("Bad: Get on roleDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRoleAssignmentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_role_assignment" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		roleAssignmentName := rs.Primary.Attributes["name"]

		client := testAccProvider.Meta().(*ArmClient).roleAssignmentsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleAssignmentName)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Role Definition still exists:\n%#v", resp)
	}

	return nil
}

func testAccAzureRMRoleAssignment_emptyName() string {
	return `
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_builtin_role_definition" "test" {
  name = "Reader"
}

resource "azurerm_role_assignment" "test" {
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_builtin_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`
}

func testAccAzureRMRoleAssignment_roleName(id string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_assignment" "test" {
  name                 = "%s"
  scope                = "${data.azurerm_subscription.primary.id}"
  role_definition_name = "Reader"
  principal_id         = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, id)
}

func testAccAzureRMRoleAssignment_builtin(id string) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

data "azurerm_builtin_role_definition" "test" {
  name = "Reader"
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${data.azurerm_subscription.primary.id}${data.azurerm_builtin_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, id)
}

func testAccAzureRMRoleAssignment_custom(roleDefinitionId string, roleAssignmentId string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

data "azurerm_client_config" "test" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Created by the Role Assignment Acceptance Test"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}

resource "azurerm_role_assignment" "test" {
  name               = "%s"
  scope              = "${data.azurerm_subscription.primary.id}"
  role_definition_id = "${azurerm_role_definition.test.id}"
  principal_id       = "${data.azurerm_client_config.test.service_principal_object_id}"
}
`, roleDefinitionId, rInt, roleAssignmentId)
}
