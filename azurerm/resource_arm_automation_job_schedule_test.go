package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationJobSchedule_basic(t *testing.T) {
	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobSchedule_basic(ri, acceptance.Location()),
				Check:  checkAccAzureRMAutomationJobSchedule_basic(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationJobSchedule_complete(t *testing.T) {
	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobSchedule_complete(ri, acceptance.Location()),
				Check:  checkAccAzureRMAutomationJobSchedule_complete(resourceName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationJobSchedule_update(t *testing.T) {
	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobSchedule_basic(ri, acceptance.Location()),
				Check:  checkAccAzureRMAutomationJobSchedule_basic(resourceName),
			},
			{
				Config: testAccAzureRMAutomationJobSchedule_complete(ri, acceptance.Location()),
				Check:  checkAccAzureRMAutomationJobSchedule_complete(resourceName),
			},
			{
				Config: testAccAzureRMAutomationJobSchedule_basic(ri, acceptance.Location()),
				Check:  checkAccAzureRMAutomationJobSchedule_basic(resourceName),
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
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_automation_job_schedule.test"
	ri := tf.AccRandTimeInt()
	location := acceptance.Location()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMAutomationJobScheduleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationJobSchedule_basic(ri, location),
				Check:  checkAccAzureRMAutomationJobSchedule_basic(resourceName),
			},
			{
				Config:      testAccAzureRMAutomationJobSchedule_requiresImport(ri, location),
				ExpectError: acceptance.RequiresImportError("azurerm_automation_job_schedule"),
			},
		},
	})
}

func testCheckAzureRMAutomationJobScheduleDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.JobScheduleClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_job_schedule" {
			continue
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		jobScheduleID := id.Path["jobSchedules"]
		jobScheduleUUID := uuid.FromStringOrNil(jobScheduleID)
		accName := rs.Primary.Attributes["account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Job Schedule: '%s'", jobScheduleUUID)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, jobScheduleUUID)

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
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Automation.JobScheduleClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		id, err := azure.ParseAzureResourceID(rs.Primary.Attributes["id"])
		if err != nil {
			return err
		}
		jobScheduleID := id.Path["jobSchedules"]
		jobScheduleUUID := uuid.FromStringOrNil(jobScheduleID)
		accName := rs.Primary.Attributes["automation_account_name"]

		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Automation Job Schedule: '%s'", jobScheduleUUID)
		}

		resp, err := conn.Get(ctx, resourceGroup, accName, jobScheduleUUID)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Automation Job Schedule '%s' (Account %q / Resource Group %q) does not exist", jobScheduleUUID, accName, resourceGroup)
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
  name                = "Output-HelloWorld"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  account_name = "${azurerm_automation_account.test.name}"
  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  publish_content_link {
    uri = "https://raw.githubusercontent.com/Azure/azure-quickstart-templates/c4935ffb69246a6058eb24f54640f53f69d3ac9f/101-automation-runbook-getvms/Runbooks/Get-AzureVMTutorial.ps1"
  }

  content = <<EOF
  param(
    [string]$Output = "World",

    [string]$Case = "Original",

    [int]$KeepCount = 10,

    [uri]$WebhookUri = "https://example.com/hook",

    [uri]$URL = "https://Example.com"
  )
  "Hello, " + $Output + "!"
EOF
}

resource "azurerm_automation_schedule" "test" {
  name                    = "acctestAS-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  frequency               = "OneTime"
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationJobSchedule_basic(rInt int, location string) string {
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

func checkAccAzureRMAutomationJobSchedule_basic(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationJobScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "job_schedule_id"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "schedule_name"),
		resource.TestCheckResourceAttrSet(resourceName, "runbook_name"),
	)
}

func testAccAzureRMAutomationJobSchedule_complete(rInt int, location string) string {
	template := testAccAzureRMAutomationJobSchedulePrerequisites(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "test" {
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  schedule_name           = "${azurerm_automation_schedule.test.name}"
  runbook_name            = "${azurerm_automation_runbook.test.name}"

  parameters = {
    output                = "Earth"
    case                  = "MATTERS"
    keepcount             = 20
    webhookuri            = "http://www.example.com/hook"
    url                   = "https://www.Example.com"
  }
}
`, template)
}

func checkAccAzureRMAutomationJobSchedule_complete(resourceName string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		testCheckAzureRMAutomationJobScheduleExists(resourceName),
		resource.TestCheckResourceAttrSet(resourceName, "job_schedule_id"),
		resource.TestCheckResourceAttrSet(resourceName, "resource_group_name"),
		resource.TestCheckResourceAttrSet(resourceName, "automation_account_name"),
		resource.TestCheckResourceAttrSet(resourceName, "schedule_name"),
		resource.TestCheckResourceAttrSet(resourceName, "runbook_name"),
		resource.TestCheckResourceAttr(resourceName, "parameters.%", "5"),
		resource.TestCheckResourceAttr(resourceName, "parameters.output", "Earth"),
		resource.TestCheckResourceAttr(resourceName, "parameters.case", "MATTERS"),
		resource.TestCheckResourceAttr(resourceName, "parameters.keepcount", "20"),
		resource.TestCheckResourceAttr(resourceName, "parameters.webhookuri", "http://www.example.com/hook"),
		resource.TestCheckResourceAttr(resourceName, "parameters.url", "https://www.Example.com"),
	)
}

func testAccAzureRMAutomationJobSchedule_requiresImport(rInt int, location string) string {
	template := testAccAzureRMAutomationJobSchedule_basic(rInt, location)
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
