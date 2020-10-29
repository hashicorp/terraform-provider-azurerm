package tests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAutomationVariableString_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_string", "test")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationVariableString_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
		},
	})
}

func testAccDataSourceAutomationVariableString_basic(data acceptance.TestData) string {
	config := testAccAzureRMAutomationVariableString_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_string" "test" {
  name                    = azurerm_automation_variable_string.test.name
  resource_group_name     = azurerm_automation_variable_string.test.resource_group_name
  automation_account_name = azurerm_automation_variable_string.test.automation_account_name
}
`, config)
}
