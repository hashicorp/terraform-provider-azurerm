package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMAutomationIntVariable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_variable_int.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationIntVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "1234"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationIntVariable_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationIntVariable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_int" "test" {
  name                    = "${azurerm_automation_variable_int.test.name}"
  resource_group_name     = "${azurerm_automation_variable_int.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_variable_int.test.automation_account_name}"
}
`, config)
}
