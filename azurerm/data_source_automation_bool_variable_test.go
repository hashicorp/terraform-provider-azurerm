package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMAutomationBoolVariable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_bool_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationBoolVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "false"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationBoolVariable_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationBoolVariable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_bool_variable" "test" {
  name                    = "${azurerm_automation_bool_variable.test.name}"
  resource_group_name     = "${azurerm_automation_bool_variable.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_bool_variable.test.automation_account_name}"
}
`, config)
}
