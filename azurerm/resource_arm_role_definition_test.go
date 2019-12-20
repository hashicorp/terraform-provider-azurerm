package azurerm

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMRoleDefinition_basic(t *testing.T) {
	resourceName := "azurerm_role_definition.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(uuid.New().String(), ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role_definition_id", "scope"},
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_role_definition.test"
	id := uuid.New().String()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(id, ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMRoleDefinition_requiresImport(id, ri),
				ExpectError: acceptance.RequiresImportError("azurerm_role_definition"),
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_complete(t *testing.T) {
	resourceName := "azurerm_role_definition.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_complete(uuid.New().String(), ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role_definition_id", "scope"},
			},
		},
	})
}

func TestAccAzureRMRoleDefinition_update(t *testing.T) {
	resourceName := "azurerm_role_definition.test"
	id := uuid.New().String()
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_basic(id, ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.not_actions.#", "0"),
				),
			},
			{
				Config: testAccAzureRMRoleDefinition_updated(id, ri),
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

func TestAccAzureRMRoleDefinition_updateEmptyId(t *testing.T) {
	resourceName := "azurerm_role_definition.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMRoleDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMRoleDefinition_emptyId(ri),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMRoleDefinitionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(resourceName, "permissions.0.not_actions.#", "0"),
				),
			},
			{
				Config: testAccAzureRMRoleDefinition_updateEmptyId(ri),
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
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
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

func testCheckAzureRMRoleDefinitionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		scope := rs.Primary.Attributes["scope"]
		roleDefinitionId := rs.Primary.Attributes["role_definition_id"]

		client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
		resp, err := client.Get(ctx, scope, roleDefinitionId)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Role Definition %q (Scope: %q) does not exist", roleDefinitionId, scope)
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

		client := acceptance.AzureProvider.Meta().(*clients.Client).Authorization.RoleDefinitionsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
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

func testAccAzureRMRoleDefinition_requiresImport(id string, rInt int) string {
	return fmt.Sprintf(`
%s

resource "azurerm_role_definition" "import" {
  role_definition_id = "${azurerm_role_definition.test.role_definition_id}"
  name               = "${azurerm_role_definition.test.name}"
  scope              = "${azurerm_role_definition.test.scope}"

  permissions {
    actions     = ["*"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, testAccAzureRMRoleDefinition_basic(id, rInt))
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
    actions          = ["*"]
    data_actions     = ["Microsoft.Storage/storageAccounts/blobServices/containers/blobs/read"]
    not_actions      = ["Microsoft.Authorization/*/read"]
    not_data_actions = []
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
  name  = "acctestrd-%d"
  scope = "${data.azurerm_subscription.primary.id}"

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

func testAccAzureRMRoleDefinition_updateEmptyId(rInt int) string {
	return fmt.Sprintf(`
data "azurerm_subscription" "primary" {}

resource "azurerm_role_definition" "test" {
  name  = "acctestrd-%d"
  scope = "${data.azurerm_subscription.primary.id}"

  permissions {
    actions     = ["*"]
    not_actions = ["Microsoft.Authorization/*/read"]
  }

  assignable_scopes = [
    "${data.azurerm_subscription.primary.id}",
  ]
}
`, rInt)
}
