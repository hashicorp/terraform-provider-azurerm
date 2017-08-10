package azurerm

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAzureRMAutomationCredential_testCredential(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationCredential_testCredential(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationCredentialExistsAndUserName("azurerm_automation_credential.test", "test_user"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAutomationCredentialDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationCredentialClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_credential" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(resourceGroup, accName, name)

		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Automation Credential still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMAutomationCredentialExistsAndUserName(name string, username string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Credential: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationCredentialClient

		resp, err := conn.Get(resourceGroup, accName, name)

		if err != nil {
			return fmt.Errorf("Bad: Get on automationCredentialClient: %s\nName: %s, Account name: %s, Resource group: %s OBJECT: %+v", err, name, accName, resourceGroup, rs.Primary)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Automation Credential '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}

		if *resp.UserName != username {
			return fmt.Errorf("Current username %s is not consistant with the checked value %s", resp.UserName, username)
		}

		return nil
	}
}

func testAccAzureRMAutomationCredential_testCredential(rInt int, location string) string {
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

resource "azurerm_automation_credential" "test" {
  name     	      = "acctest-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  user_name           = "test_user"
  password            = "test_pwd"
  description         = "This is a test credential for terraform acceptance test"
}
`, rInt, location, rInt, rInt)
}
