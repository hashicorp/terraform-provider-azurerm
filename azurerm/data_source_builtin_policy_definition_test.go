package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceAzureRMBuiltInPolicyDefinition_AllowedResourceTypes(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_policy_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinition("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMBuiltInPolicyDefinition_AtManagementGroup_AllowedResourceTypes(t *testing.T) {
	dataSourceName := "data.azurerm_builtin_policy_definition.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBuiltInPolicyDefinitionAtManagementGroup("Allowed resource types"),
				Check: resource.ComposeTestCheckFunc(
					testAzureRMClientConfigAttr(dataSourceName, "id", "/providers/Microsoft.Authorization/policyDefinitions/a08ec900-254a-4555-9bf5-e42af04b5c5c"),
				),
			},
		},
	})
}

func testAccDataSourceBuiltInPolicyDefinition(name string) string {
	return fmt.Sprintf(`
data "azurerm_builtin_policy_definition" "test" {
  name = "%s"
}
`, name)
}

func testAccDataSourceBuiltInPolicyDefinitionAtManagementGroup(name string) string {
	return fmt.Sprintf(`

data "azurerm_client_config" "current" {}

data "azurerm_builtin_policy_definition" "test" {
  name = "%s"
  management_group_id = "${data.azurerm_client_config.current.tenant_id}"
}
`, name)
}
