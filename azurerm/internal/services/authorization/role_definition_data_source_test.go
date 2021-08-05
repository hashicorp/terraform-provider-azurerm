package authorization_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
)

type RoleDefinitionDataSource struct{}

func TestAccRoleDefinitionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	id := uuid.New().String()

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.basic(id, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.0").HasValue("*"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("3"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.0").HasValue("Microsoft.Authorization/*/Delete"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.1").HasValue("Microsoft.Authorization/*/Write"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.2").HasValue("Microsoft.Authorization/elevateAccess/Action"),
			),
		},
	})
}

func TestAccRoleDefinitionDataSource_basicByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	id := uuid.New().String()

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.byName(id, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.0").HasValue("*"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("3"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.0").HasValue("Microsoft.Authorization/*/Delete"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.1").HasValue("Microsoft.Authorization/*/Write"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.2").HasValue("Microsoft.Authorization/elevateAccess/Action"),
			),
		},
	})
}

func TestAccRoleDefinitionDataSource_builtIn_contributor(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.builtIn("Contributor"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.0").HasValue("*"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("6"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.0").HasValue("Microsoft.Authorization/*/Delete"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.1").HasValue("Microsoft.Authorization/*/Write"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.2").HasValue("Microsoft.Authorization/elevateAccess/Action"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.3").HasValue("Microsoft.Blueprint/blueprintAssignments/write"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.4").HasValue("Microsoft.Blueprint/blueprintAssignments/delete"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.5").HasValue("Microsoft.Compute/galleries/share/action"),
			),
		},
	})
}

func TestAccRoleDefinitionDataSource_builtIn_owner(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.builtIn("Owner"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.0").HasValue("*"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("0"),
			),
		},
	})
}

func TestAccRoleDefinitionDataSource_builtIn_reader(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.builtIn("Reader"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.actions.0").HasValue("*/read"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("0"),
			),
		},
	})
}

func TestAccRoleDefinitionDataSource_builtIn_virtualMachineContributor(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: RoleDefinitionDataSource{}.builtIn("Virtual Machine Contributor"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").HasValue("/providers/Microsoft.Authorization/roleDefinitions/9980e02c-c2be-4d73-94e8-173b1dc7cf3c"),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("type").Exists(),
				check.That(data.ResourceName).Key("permissions.#").HasValue("1"),
				check.That(data.ResourceName).Key("permissions.0.not_actions.#").HasValue("0"),
			),
		},
	})
}

func (d RoleDefinitionDataSource) builtIn(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "test" {
  name = "%s"
}
`, name)
}

func (d RoleDefinitionDataSource) basic(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

resource "azurerm_role_definition" "test" {
  role_definition_id = "%s"
  name               = "acctestrd-%d"
  scope              = data.azurerm_subscription.primary.id
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
    data.azurerm_subscription.primary.id,
  ]
}

data "azurerm_role_definition" "test" {
  role_definition_id = azurerm_role_definition.test.role_definition_id
  scope              = data.azurerm_subscription.primary.id
}
`, id, data.RandomInteger)
}

func (d RoleDefinitionDataSource) byName(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_role_definition" "byName" {
  name  = azurerm_role_definition.test.name
  scope = data.azurerm_subscription.primary.id
}
`, d.basic(id, data))
}
