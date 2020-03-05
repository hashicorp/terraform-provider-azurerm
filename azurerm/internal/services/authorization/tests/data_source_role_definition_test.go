package tests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMRoleDefinition_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleDefinition_basic(id, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_basicByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")

	id := uuid.New().String()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRoleDefinition_byName(id, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(data.ResourceName, "name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_contributor(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "5"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.3", "Microsoft.Blueprint/blueprintAssignments/write"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.4", "Microsoft.Blueprint/blueprintAssignments/delete"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_owner(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Owner"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_reader(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Reader"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.actions.0", "*/read"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMRoleDefinition_builtIn_virtualMachineContributor(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_role_definition", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMRoleDefinition_builtIn("Virtual Machine Contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/9980e02c-c2be-4d73-94e8-173b1dc7cf3c"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "description"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "type"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMRoleDefinition_builtIn(name string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_role_definition" "test" {
  name = "%s"
}
`, name)
}

func testAccDataSourceRoleDefinition_basic(id string, data acceptance.TestData) string {
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

func testAccDataSourceRoleDefinition_byName(id string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_role_definition" "byName" {
  name  = azurerm_role_definition.test.name
  scope = data.azurerm_subscription.primary.id
}
`, testAccDataSourceRoleDefinition_basic(id, data))
}
