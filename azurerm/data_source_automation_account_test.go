package azurerm

import (
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceAutomationAccountRegistrationInformation(t *testing.T) {
	dataSourceName := "data.azurerm_automation_account_registration_info.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAutomationAccountRegistrationInformation_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "primary_key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "secondary_key"),
					resource.TestCheckResourceAttr(dataSourceName, "endpoint"),
				),
			},
		},
	})
}

func testAccResourceAutomationAccountRegistrationInformation_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
resource "azurerm_automation_account" "test" {
  name                = "automationAccount1"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name = "Basic"
}
data "azurerm_automation_account_registration_info" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_automation_account.test.name}"
}
`, rInt, location)
}
