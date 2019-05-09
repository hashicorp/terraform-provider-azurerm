package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
)

func TestAccDataSourceAzureRMAutomationStringVariable_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_string_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationStringVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "Hello, Terraform Basic Test."),
				),
			},
		},
	})
}

func testAccDataSourceAutomationStringVariable_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationStringVariable_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_string_variable" "test" {
  name                    = "${azurerm_automation_string_variable.test.name}"
  resource_group_name     = "${azurerm_automation_string_variable.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_string_variable.test.automation_account_name}"
}
`, config)
}
