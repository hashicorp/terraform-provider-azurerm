package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationModule_basic(t *testing.T) {
	resourceName := "azurerm_automation_module.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationModuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationModule_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationModuleExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Module link is not returned by api in Get operation
				ImportStateVerifyIgnore: []string{"module_link"},
			},
		},
	})
}

func testCheckAzureRMAutomationModuleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationModuleClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_module" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Module: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Module still exists:\n%#v", resp)

	}

	return nil
}

func testCheckAzureRMAutomationModuleExists(name string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Module: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationModuleClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Module '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationModuleClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationModule_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_module" "test" {
  name                    = "xActiveDirectory"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"

  module_link = {
    uri = "https://devopsgallerystorage.blob.core.windows.net/packages/xactivedirectory.2.19.0.nupkg"
  }
}
`, rInt, location, rInt)
}
