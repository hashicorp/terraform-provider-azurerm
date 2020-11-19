package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAutomationVariableBool_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_bool", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationVariableBool_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "value", "false"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationVariableBool_basic(data acceptance.TestData) string {
	config := testAccAzureRMAutomationVariableBool_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_bool" "test" {
  name                    = azurerm_automation_variable_bool.test.name
  resource_group_name     = azurerm_automation_variable_bool.test.resource_group_name
  automation_account_name = azurerm_automation_variable_bool.test.automation_account_name
}
`, config)
}
