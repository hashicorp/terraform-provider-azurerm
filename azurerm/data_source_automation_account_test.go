package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAutomationAccount(t *testing.T) {
	dataSourceName := "data.azurerm_automation_account.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()
	resourceGroupName := fmt.Sprintf("acctestRG-%d", ri)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationAccount_complete(resourceGroupName, location, ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "resource_group_name", resourceGroupName),
				),
			},
		},
	})
}

func testAccDataSourceAutomationAccount_complete(resourceGroupName string, location string, ri int) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}
resource "azurerm_automation_account" "test" {
  name                = "acctestautomationAccount-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku_name = "Basic"
}
data "azurerm_automation_account" "test" {
  resource_group_name = "${azurerm_resource_group.test.name}"
  name                = "${azurerm_automation_account.test.name}"
}
`, resourceGroupName, location, ri)
}
