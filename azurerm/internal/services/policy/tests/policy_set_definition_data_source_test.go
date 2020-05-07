package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPolicySetDefinition_byName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPolicySetDefinition_byName(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestPolSet-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("acctestPolSet-display-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_type", "Custom"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "parameters"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_definitions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicySetDefinition_byDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPolicySetDefinition_byDisplayName(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestPolSet-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("acctestPolSet-display-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_type", "Custom"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "parameters"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_definitions"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPolicySetDefinition_byName(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  name = azurerm_policy_set_definition.test.name
}
`, template)
}

func testAccDataSourceAzureRMPolicySetDefinition_byDisplayName(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  display_name = azurerm_policy_set_definition.test.display_name
}
`, template)
}
