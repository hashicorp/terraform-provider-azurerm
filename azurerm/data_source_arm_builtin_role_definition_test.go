package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMBuiltInRoleDefinition_contributor(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Contributor"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "b24988ac-6180-42a0-ab88-20f7382dd24c"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_owner(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Owner"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "8e3af657-a8ff-443c-a75c-2fe8c4bcb635"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_reader(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("Reader"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "acdd72a7-3385-48ef-bd42-f606fba81ae7"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInRoleDefinition_virtualMachineContributor(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_role_definition.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInRoleDefinition("VirtualMachineContributor"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "d73bb868-a0df-4d4d-bd69-98a00b01fccb"),
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
