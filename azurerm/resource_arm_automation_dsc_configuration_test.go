package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationDscConfiguration_basic(t *testing.T) {
	resourceName := "azurerm_automation_dsc_configuration.test"
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationDscConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationDscConfiguration_basic(ri, testLocation()),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationDscConfigurationExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "location"),
					resource.TestCheckResourceAttr(resourceName, "description", "test"),
					resource.TestCheckResourceAttrSet(resourceName, "log_verbose"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttr(resourceName, "content_embedded", "configuration acctest {}"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMAutomationDscConfigurationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationDscConfigurationClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_dsc_configuration" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Dsc Configuration: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Dsc Configuration still exists:\n%#v", resp)

	}

	return nil
}

func testCheckAzureRMAutomationDscConfigurationExists(name string) resource.TestCheckFunc {

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
			return fmt.Errorf("Bad: no resource group found in state for Automation Dsc Configuration: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationDscConfigurationClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, resourceGroup, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Dsc Configuration '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationDscConfigurationClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		return nil
	}
}

func testAccAzureRMAutomationDscConfiguration_basic(rInt int, location string) string {
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

resource "azurerm_automation_dsc_configuration" "test" {
  name                    = "acctest"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  location                = "${azurerm_resource_group.test.location}"
  content_embedded        = "configuration acctest {}"
  description             = "test"
}
`, rInt, location, rInt)
}
