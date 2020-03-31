package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMPolicySetDefinition_customByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPolicySetDefinition_customByName(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestpolset-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("acctestpolset-display-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_type", "Custom"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "parameters"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_definitions"),
				),
			},
		},
	})
}

func TestAccDataSourceAzureRMPolicySetDefinition_customByDisplayName(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_policy_set_definition", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMPolicySetDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAzureRMPolicySetDefinition_customByDisplayName(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "name", fmt.Sprintf("acctestpolset-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "display_name", fmt.Sprintf("acctestpolset-display-%d", data.RandomInteger)),
					resource.TestCheckResourceAttr(data.ResourceName, "policy_type", "Custom"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "parameters"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "policy_definitions"),
				),
			},
		},
	})
}

func testAccDataSourceAzureRMPolicySetDefinition_customByName(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  name = azurerm_policy_set_definition.test.name
}
`, template)
}

func testAccDataSourceAzureRMPolicySetDefinition_customByDisplayName(data acceptance.TestData) string {
	template := testAzureRMPolicySetDefinition_custom(data)
	return fmt.Sprintf(`
%s

data "azurerm_policy_set_definition" "test" {
  display_name = azurerm_policy_set_definition.test.display_name
}
`, template)
}
