package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMAutomationDatetimeVariable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_variable_datetime.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationDatetimeVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "2019-04-24T21:40:54.074Z"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationDatetimeVariable_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationDatetimeVariable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_datetime" "test" {
  name                    = "${azurerm_automation_variable_datetime.test.name}"
  resource_group_name     = "${azurerm_automation_variable_datetime.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_variable_datetime.test.automation_account_name}"
}
`, config)
}
