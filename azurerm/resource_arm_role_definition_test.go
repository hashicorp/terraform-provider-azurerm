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

func TestAccAzureRMRoleDefinition_basic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRoleDefinition_basic(uuid.New().String(), ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists("azurerm_role_definition.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_complete(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMRoleDefinition_complete(uuid.New().String(), ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists("azurerm_role_definition.test"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_update(t *testing.T) {
	resourceName := "azurerm_role_definition.test"
	id := uuid.New().String()
	ri := acctest.RandInt()

	config := testAccAzureRMRoleDefinition_basic(id, ri)
	updatedConfig := testAccAzureRMRoleDefinition_updated(id, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.not_actions.#", "0"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.not_actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/read"),
				),
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_emptyName(t *testing.T) {
	resourceName := "azurerm_role_definition.test"

	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_emptyId(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
				),
			},
		},
	})
}

func testCheckAzureRMRoleDefinitionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %q", name)
		}

		scope := rs.Primary.Attributes["scope"]
		roleDefinitionId := rs.Primary.Attributes["role_definition_id"]

		client := testAccProvider.Meta().(*ArmClient).roleDefinitionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleDefinitionId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Definition %q (Scope: %q) does not exist", name, scope)
			}
			return fmt.Errorf("Bad: Get on roleDefinitionsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMRoleDefinitionDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_role_definition" {
			continue
		}

		scope := rs.Primary.Attributes["scope"]
		roleDefinitionId := rs.Primary.Attributes["role_definition_id"]

		client := testAccProvider.Meta().(*ArmClient).roleDefinitionsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
		resp, err := client.Get(ctx, scope, roleDefinitionId)

		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return err
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMRoleDefinition_basic(id string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, id, rInt)
}

func testAccAzureRMRoleDefinition_complete(id string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Acceptance Test Role Definition"

  permissions {
    actions     = ["*"]
    not_actions = ["Microsoft.Authorization/*/read"]
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, id, rInt)
}

func testAccAzureRMRoleDefinition_updated(id string, rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d-updated"
  scope              = "${data.azurerm_subscription.primary.id}"
  description        = "Acceptance Test Role Definition"

  permissions {
    actions     = ["*"]
    not_actions = ["Microsoft.Authorization/*/read"]
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, id, rInt)
}

func testAccAzureRMRoleDefinition_emptyId(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  name               = "acctestrd-%d"
  scope              = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, rInt)
}
