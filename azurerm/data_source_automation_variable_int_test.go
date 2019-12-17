package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAzureRMAutomationVariableInt_basic(t *testing.T) {
	dataSourceName := "data.azurerm_automation_variable_int.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationVariableInt_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "value", "1234"),
				),
			},
		},
	})
}

func testAccDataSourceAutomationVariableInt_basic(rInt int, location string) string {
	config := testAccAzureRMAutomationVariableInt_basic(rInt, location)
	return fmt.Sprintf(`
%s

data "azurerm_automation_variable_int" "test" {
  name                    = "${azurerm_automation_variable_int.test.name}"
  resource_group_name     = "${azurerm_automation_variable_int.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_variable_int.test.automation_account_name}"
}
`, config)
}
