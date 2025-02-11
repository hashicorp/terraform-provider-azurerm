// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/runbook"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2023-11-01/schedule"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type AutomationJobScheduleResource struct{}

func TestAccAutomationJobSchedule_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_updateRunbook(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data, "Update Runbook auto update"),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutomationJobSchedule_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automation_job_schedule", "test")
	r := AutomationJobScheduleResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_automation_job_schedule"),
		},
	})
}

func (t AutomationJobScheduleResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := commonids.ParseCompositeResourceID(state.ID, &schedule.ScheduleId{}, &runbook.RunbookId{})
	if err != nil {
		return nil, err
	}

	resp, err := automation.GetJobScheduleFromTFID(ctx, clients.Automation.JobSchedule, id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %v", *id, err)
	}

	return pointer.To(resp != nil), nil
}

func (AutomationJobScheduleResource) template(data acceptance.TestData, runbookDesc ...string) string {
	var description string
	if len(runbookDesc) > 0 {
		description = runbookDesc[0]
	}

	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-auto-%[1]d"
  location = "%[2]s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAA-%[1]d"
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
  description             = "This is a test runbook for terraform acceptance test.%[3]s"
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
  name                    = "acctestAS-%[1]d"
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  frequency               = "OneTime"
}
`, data.RandomInteger, data.Locations.Primary, description)
}

func (AutomationJobScheduleResource) basic(data acceptance.TestData, runbookDesc ...string) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_job_schedule" "test" {
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name
  schedule_name           = azurerm_automation_schedule.test.name
  runbook_name            = azurerm_automation_runbook.test.name
}
`, AutomationJobScheduleResource{}.template(data, runbookDesc...))
}

func (AutomationJobScheduleResource) complete(data acceptance.TestData, runbookDesc ...string) string {
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
`, AutomationJobScheduleResource{}.template(data, runbookDesc...))
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
