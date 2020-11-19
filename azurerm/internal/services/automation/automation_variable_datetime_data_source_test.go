package automation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAutomationVariableDateTime_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_variable_datetime", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationVariableDateTime_basic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationVariableDateTime_basic(data acceptance.TestData) string {
	config := testAccAzureRMAutomationVariableDateTime_basic(data)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_datetime" "test" {
  name                    = azurerm_automation_variable_datetime.test.name
  resource_group_name     = azurerm_automation_variable_datetime.test.resource_group_name
  automation_account_name = azurerm_automation_variable_datetime.test.automation_account_name
}
`, config)
}
