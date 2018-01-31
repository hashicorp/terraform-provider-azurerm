package azurerm

import (
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/automation/mgmt/2015-10-31/automation"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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
			},
		},
	})
}

func testCheckAzureRMAutomationScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationScheduleClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_schedule" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Schedule: '%s'", name)
		}
		conn.ResourceGroupName = resourceGroup

		resp, err := conn.Get(ctx, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Schedule still exists:\n%#v", resp)
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
		conn.ResourceGroupName = resourceGroup
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		resp, err := conn.Get(ctx, accName, name)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Schedule '%s' (resource group: '%s') does not exist", name, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationScheduleClient: %+v", err)
		}

		if resp.Frequency != freq {
			return fmt.Errorf("Current frequency %s is not consistent with checked value %s", resp.Frequency, freq)
		}
		return nil
	}
}

func testAccAzureRMAutomationSchedule_oneTime(rInt int, location string) string {
	startTime := time.Now().UTC().Add(time.Duration(7) * time.Minute)
	startTime = startTime.Add(time.Duration(-1*startTime.Second()) * time.Second)

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

resource "azurerm_automation_schedule" "test" {
  name	 	          = "OneTimer-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  account_name        = "${azurerm_automation_account.test.name}"
  frequency	          = "OneTime"
  timezone	          = "Central Europe Standard Time"
  start_time	      = "%s"
  description	      = "This is a test runbook for terraform acceptance test"
}
`, rInt, location, rInt, rInt, startTime.Format(time.RFC3339))
}
