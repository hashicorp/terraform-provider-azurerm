package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/Azure/azure-sdk-for-go/arm/automation"
)

func TestAccAzureRMAutomationAccount_skuBasic(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationAccount_skuBasic(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExistsAndSku("azurerm_automation_account.test", automation.Basic),
				),
			},
		},
	})
}

func TestAccAzureRMAutomationAccount_skuFree(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationAccount_skuFree(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationAccountExistsAndSku("azurerm_automation_account.test", automation.Free),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationAccountDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationAccountClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_account" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Automation Account still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMAutomationAccountExistsAndSku(name string, sku automation.SkuNameEnum) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Account: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationAccountClient

		resp, err := conn.Get(resourceGroup, name)

		if err != nil {
			return fmt.Errorf("Bad: Get on automationClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Automation Account '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		if resp.Sku.Name != sku {
			return fmt.Errorf("Actual sku %s is not consistent with the checked value %s", resp.Sku.Name, sku)
		}

		return nil
	}
}

func testAccAzureRMAutomationAccount_skuBasic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name = "acctestRG-%d"
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
`, rInt, location, rInt)
}

func testAccAzureRMAutomationAccount_skuFree(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
 name = "acctestRG-%d"
 location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctest-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku {
        name = "Free"
  }
}
`, rInt, location, rInt)
}
