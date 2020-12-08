package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	uuid "github.com/satori/go.uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type AutomationJobScheduleResource struct {
}

func TestAccAutomationJobSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_automation_job_schedule"),
		},
	})
}

func (t AutomationJobScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	jobScheduleID := id.Path["jobSchedules"]
	jobScheduleUUID := uuid.FromStringOrNil(jobScheduleID)
	resourceGroup := id.ResourceGroup
	accountName := id.Path["automationAccounts"]

	resp, err := clients.Automation.JobScheduleClient.Get(ctx, resourceGroup, accountName, jobScheduleUUID)
	if err != nil {
		return nil, fmt.Errorf("retrieving Automation Job Schedule '%s' (Account %q / Resource Group %q) does not exist", jobScheduleUUID, accountName, resourceGroup)
	}

	return utils.Bool(resp.JobScheduleProperties != nil), nil
}

func (AutomationJobScheduleResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_automation_runbook" "test" {
  name                    = "Output-HelloWorld"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  log_verbose             = "true"
  log_progress            = "true"
  description             = "This is a test runbook for terraform acceptance test"
  runbook_type            = "PowerShell"

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
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (AutomationJobScheduleResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  schedule_name           = azurerm_automation_schedule.test.name
  runbook_name            = azurerm_automation_runbook.test.name
}
`, AutomationJobScheduleResource{}.template(data))
}

func (AutomationJobScheduleResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  schedule_name           = azurerm_automation_schedule.test.name
  runbook_name            = azurerm_automation_runbook.test.name

  parameters = {
    output     = "Earth"
    case       = "MATTERS"
    keepcount  = 20
    webhookuri = "http://www.example.com/hook"
    url        = "https://www.Example.com"
  }
}
`, AutomationJobScheduleResource{}.template(data))
}

func (AutomationJobScheduleResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "import" {
  resource_group_name     = azurerm_automation_job_schedule.test.resource_group_name
  automation_account_name = azurerm_automation_job_schedule.test.automation_account_name
  schedule_name           = azurerm_automation_job_schedule.test.schedule_name
  runbook_name            = azurerm_automation_job_schedule.test.runbook_name
  job_schedule_id         = azurerm_automation_job_schedule.test.job_schedule_id
}
`, AutomationJobScheduleResource{}.basic(data))
}
