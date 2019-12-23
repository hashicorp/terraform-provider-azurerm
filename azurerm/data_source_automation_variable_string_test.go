package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAutomationVariableString_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_variable_string.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationVariableString_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
		},
	})
}

func testAccDataSourceAutomationVariableString_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationVariableString_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_string" "test" {
  name                    = "${azurerm_automation_variable_string.test.name}"
  resource_group_name     = "${azurerm_automation_variable_string.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_variable_string.test.automation_account_name}"
}
`, config)
}
