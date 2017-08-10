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

func TestAccAzureRMAutomationSchedule_oneTime(t *testing.T) {
	ri := acctest.RandInt()
	config := testAccAzureRMAutomationSchedule_oneTime(ri, testLocation())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationScheduleExistsAndFrequencyType("azurerm_automation_schedule.test", automation.OneTime),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testCheckAzureRMAutomationScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_schedule" {
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
			return fmt.Errorf("Automation Schedule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMAutomationScheduleExistsAndFrequencyType(name string, freq automation.ScheduleFrequency) resource.TestCheckFunc {

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
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}

		conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient

		resp, err := conn.Get(resourceGroup, accName, name)

		if err != nil {
			return fmt.Errorf("Bad: Get on automationScheduleClient: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Automation Schedule '%s' (resource group: '%s') does not exist", name, resourceGroup)
		}


		if resp.Frequency != freq {
			return fmt.Errorf("Current frequency %s is not consistent with checked value %s", resp.Frequency, freq) 
		}
		return nil
	}
}

func testAccAzureRMAutomationSchedule_oneTime(rInt int, location string) string {
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

resource "azurerm_automation_schedule" "test" {
  name	 	      = "OneTimer-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  frequency	      = "OneTime"
  first_run {
        "hour" = 20
        "minute" = 5
        "second" = 0
  }
  description	      = "This is a test runbook for terraform acceptance test"
}
`, rInt, location, rInt, rInt)
}
