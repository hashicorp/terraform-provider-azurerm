package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMBuiltInRoleDefinition_contributor(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Contributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/b24988ac-6180-42a0-ab88-20f7382dd24c"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "5"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.0", "Microsoft.Authorization/*/Delete"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.1", "Microsoft.Authorization/*/Write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.2", "Microsoft.Authorization/elevateAccess/Action"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.3", "Microsoft.Blueprint/blueprintAssignments/write"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.4", "Microsoft.Blueprint/blueprintAssignments/delete"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_owner(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Owner"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/8e3af657-a8ff-443c-a75c-2fe8c4bcb635"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_reader(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Reader"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/acdd72a7-3385-48ef-bd42-f606fba81ae7"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "*/read"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_virtualMachineContributor(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("VirtualMachineContributor"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/roleDefinitions/9980e02c-c2be-4d73-94e8-173b1dc7cf3c"),
					resource.TestCheckResourceAttrSet(dataSourceName, "description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "type"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.#", "38"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.0", "Microsoft.Authorization/*/read"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.actions.15", "Microsoft.Network/networkSecurityGroups/join/action"),
					resource.TestCheckResourceAttr(dataSourceName, "permissions.0.not_actions.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceBuiltInRoleDefinition(name string) string {
	return fmt.Sprintf(`
data "azurerm_builtin_role_definition" "test" {
  name = "%s"
}
`, name)
}
