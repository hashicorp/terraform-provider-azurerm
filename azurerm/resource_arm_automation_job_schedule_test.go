package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationJobScheduleCreate(t *testing.T) {
	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobScheduleCreate(ri, testLocation()),
				Check:  checkAccAzureRMAutomationJobScheduleCreate(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationJobSchedule_requiresImport(t *testing.T) {
	if !requireResourcesToBeImported {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobScheduleCreate(ri, location),
				Check:  checkAccAzureRMAutomationJobScheduleCreate(resourceName),
			},
			{
				Config:      testAccAzureRMAutomationJobSchedule_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_automation_job_schedule"),
			},
		},
	})
}

func testCheckAzureRMAutomationJobScheduleDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*ArmClient).automationJobScheduleClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_job_schedule" {
			continue
		}

		id, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		name := id.Path["jobSchedules"]
		nameUUID := uuid.FromStringOrNil(name)
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Job Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, nameUUID)

		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}

		return fmt.Errorf("Automation Job Schedule still exists:\n%#v", resp)
	}

	return nil
}

func testCheckAzureRMAutomationJobScheduleExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*ArmClient).automationJobScheduleClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := parseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		name := id.Path["jobSchedules"]
		nameUUID := uuid.FromStringOrNil(name)
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Job Schedule: '%s'", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, nameUUID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Job Schedule '%s' (Account %q / Resource Group %q) does not exist", name, accName, resourceGroup)
			}

			return fmt.Errorf("Bad: Get on automationJobScheduleClient: %+v", err)
		}

		return nil
	}
}

func testAccAzureRMAutomationJobSchedulePrerequisites(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku {
    name = "Basic"
  }
}

resource "azurerm_automation_runbook" "test" {
  name                = "Get-AzureVMTutorial"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  account_name = "${azurerm_automation_account.test.name}"
  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShellWorkflow"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/master/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }
}

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationJobScheduleCreate(rInt int, location string) string {
	template := testAccAzureRMAutomationJobSchedulePrerequisites(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  schedule_name           = "${azurerm_automation_schedule.test.name}"
  runbook_name            = "${azurerm_automation_runbook.test.name}"
}
`, template)
}

func checkAccAzureRMAutomationJobScheduleCreate(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationJobScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "schedule_name"),
		resource.TestCheckResourceAttrSet(resourceName, "runbook_name"),
	)
}

func testAccAzureRMAutomationJobSchedule_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationJobScheduleCreate(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "import" {
  resource_group_name     = "${azurerm_automation_job_schedule.test.resource_group_name}"
  automation_account_name = "${azurerm_automation_job_schedule.test.automation_account_name}"
  schedule_name           = "${azurerm_automation_job_schedule.test.schedule_name}"
  runbook_name            = "${azurerm_automation_job_schedule.test.runbook_name}"
}
`, template)
}
