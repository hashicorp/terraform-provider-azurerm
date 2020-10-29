package tests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
)

func TestAccDataSourceAutomationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_automation_account", "test")
	resourceGroupName := fmt.Sprintf("acctestRG-%d", data.RandomInteger)

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { acceptance.PreCheck(t) },
		Providers: acceptance.SupportedProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAutomationAccount_complete(resourceGroupName, data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(data.ResourceName, "resource_group_name", resourceGroupName),
					resource.TestMatchResourceAttr(data.ResourceName, "id",
						regexp.MustCompile(`^/subscriptions/[^/]+/resourceGroups/[^/]+/providers/Microsoft\.Automation/automationAccounts/[^/]+$`)),
				),
			},
		},
	})
}

func testAccDataSourceAutomationAccount_complete(resourceGroupName string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "%s"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestautomationAccount-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

data "azurerm_automation_account" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = azurerm_automation_account.test.name
}
`, resourceGroupName, data.Locations.Primary, data.RandomInteger)
}
