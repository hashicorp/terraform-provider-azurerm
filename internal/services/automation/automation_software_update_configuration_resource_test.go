// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automation_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/softwareupdateconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SoftwareUpdateConfigurationResource struct {
	startTime  string
	expireTime string
}

func newSoftwareUpdateConfigurationResource() SoftwareUpdateConfigurationResource {
	// The start time of the schedule must be at least 5 minutes after the time you create the schedule,
	// so we cannot hardcode the time string.
	// we use timezone as UTC so the time string should be in UTC format
	ref := time.Now().Round(time.Minute)
	ins := SoftwareUpdateConfigurationResource{
		startTime:  ref.Add(time.Hour * 1).In(time.UTC).Format(time.RFC3339),
		expireTime: ref.Add(time.Hour * 2).In(time.UTC).Format(time.RFC3339),
	}
	return ins
}

func (a SoftwareUpdateConfigurationResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := softwareupdateconfiguration.ParseSoftwareUpdateConfigurationID(state.ID)
	if err != nil {
		return nil, err
	}
	resp, err := client.Automation.SoftwareUpdateConfigClient.GetByName(ctx, *id, softwareupdateconfiguration.DefaultGetByNameOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func TestAccSoftwareUpdateConfiguration_linuxBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_linuxComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_linuxUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.linuxComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_CompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_withTask(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.withTask(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_defaultTimeZone(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.defaultTimeZone(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_update(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.linuxBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_windowsBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_windowsComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func TestAccSoftwareUpdateConfiguration_windowsUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, automation.SoftwareUpdateConfigurationResource{}.ResourceType(), "test")
	r := newSoftwareUpdateConfigurationResource()
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.windowsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.windowsComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
		{
			Config: r.windowsBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		// scheduleInfo.advancedSchedule always returns null - https://github.com/Azure/azure-rest-api-specs/issues/24436
		data.ImportStep("schedule.0.advanced", "schedule.0.monthly_occurrence"),
	})
}

func (a SoftwareUpdateConfigurationResource) defaultTimeZone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  linux {
    classification_included = "Security"
    excluded_packages       = ["apt"]
    included_packages       = ["vim"]
    reboot                  = "IfRequired"
  }

  duration            = "PT1H1M1S"
  virtual_machine_ids = []

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
    }

    non_azure_query {
      function_alias = "savedSearch1"
      workspace_id   = azurerm_log_analytics_workspace.test.id
    }
  }

  schedule {
    description         = "foo-schedule"
    start_time          = "%[3]s"
    is_enabled          = true
    interval            = 1
    frequency           = "Hour"
    advanced_week_days  = ["Monday", "Tuesday"]
    advanced_month_days = [1, 10, 15]
    monthly_occurrence {
      occurrence = 1
      day        = "Tuesday"
    }
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) linuxBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  linux {
    classifications_included = ["Security"]
  }

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
    }
  }

  schedule {
    frequency = "OneTime"
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger)
}

func (a SoftwareUpdateConfigurationResource) linuxComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  linux {
    classifications_included = ["Critical", "Security"]
    excluded_packages        = ["apt"]
    included_packages        = ["vim"]
    reboot                   = "RebootOnly"
  }

  duration            = "PT1H1M1S"
  virtual_machine_ids = []

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
    }

    non_azure_query {
      function_alias = "savedSearch1"
      workspace_id   = azurerm_log_analytics_workspace.test.id
    }
  }

  schedule {
    description         = "foo-schedule"
    start_time          = "%[3]s"
    expiry_time         = "%[4]s"
    is_enabled          = true
    interval            = 1
    frequency           = "Hour"
    time_zone           = "Etc/UTC"
    advanced_week_days  = ["Monday", "Tuesday"]
    advanced_month_days = [1, 10, 15]
    monthly_occurrence {
      occurrence = 1
      day        = "Tuesday"
    }
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) withTask(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_automation_runbook" "test" {
  name                    = "Get-AzureVMTutorial"
  location                = azurerm_resource_group.test.location
  resource_group_name     = azurerm_resource_group.test.name
  automation_account_name = azurerm_automation_account.test.name

  log_verbose  = "true"
  log_progress = "true"
  description  = "This is a test runbook for terraform acceptance test"
  runbook_type = "PowerShell"

  content = <<CONTENT
# Some test content
# for Terraform acceptance test
CONTENT
  tags = {
    ENV = "runbook_test"
  }
}

%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  linux {
    classification_included = "Security"
    excluded_packages       = ["apt"]
    included_packages       = ["vim"]
    reboot                  = "IfRequired"
  }

  duration            = "PT1H1M1S"
  virtual_machine_ids = []

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
    }

    non_azure_query {
      function_alias = "savedSearch1"
      workspace_id   = azurerm_log_analytics_workspace.test.id
    }
  }

  schedule {
    description         = "foo-schedule"
    start_time          = "%[3]s"
    is_enabled          = true
    interval            = 1
    frequency           = "Hour"
    time_zone           = "Etc/UTC"
    advanced_week_days  = ["Monday", "Tuesday"]
    advanced_month_days = [1, 10, 15]
    monthly_occurrence {
      occurrence = 1
      day        = "Tuesday"
    }
  }

  pre_task {
    source = azurerm_automation_runbook.test.name
    parameters = {
      COMPUTERNAME = "Foo"
    }
  }

  post_task {
    source = azurerm_automation_runbook.test.name
    parameters = {
      COMPUTERNAME = "Foo"
    }
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`

%s

data "azurerm_client_config" "current" {}

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  linux {
    classification_included = "Security"
    excluded_packages       = ["apt"]
    included_packages       = ["vim"]
    reboot                  = "Always"
  }

  duration            = "PT2H2M2S"
  virtual_machine_ids = []

  target {
    azure_query {
      scope     = ["/subscriptions/${data.azurerm_client_config.current.subscription_id}"]
      locations = [azurerm_resource_group.test.location]
      tags {
        tag    = "foo"
        values = ["barbar2"]
      }
      tag_filter = "Any"
    }

    non_azure_query {
      function_alias = "savedSearch2"
      workspace_id   = azurerm_log_analytics_workspace.test.id
    }
  }

  schedule {
    description        = "foobar-schedule"
    start_time         = "%[3]s"
    expiry_time        = "%[4]s"
    is_enabled         = true
    interval           = 2
    frequency          = "Hour"
    time_zone          = "Etc/UTC"
    advanced_week_days = ["Monday", "Tuesday"]
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) windowsBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  windows {
    classifications_included = ["Security"]
  }

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
    }
  }

  schedule {
    frequency = "OneTime"
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) windowsComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`


%s

resource "azurerm_automation_software_update_configuration" "test" {
  automation_account_id = azurerm_automation_account.test.id
  name                  = "acctest-suc-%[2]d"

  windows {
    classifications_included = ["Critical", "Security"]
    reboot                   = "RebootOnly"
  }

  duration            = "PT1H1M1S"
  virtual_machine_ids = []

  target {
    azure_query {
      scope     = [azurerm_resource_group.test.id]
      locations = [azurerm_resource_group.test.location]
      tags {
        tag    = "foo"
        values = ["barbar2"]
      }
      tag_filter = "Any"
    }

    non_azure_query {
      function_alias = "savedSearch1"
      workspace_id   = azurerm_log_analytics_workspace.test.id
    }
  }

  schedule {
    description         = "foo-schedule"
    start_time          = "%[3]s"
    expiry_time         = "%[4]s"
    is_enabled          = true
    interval            = 1
    frequency           = "Hour"
    time_zone           = "Etc/UTC"
    advanced_week_days  = ["Monday", "Tuesday"]
    advanced_month_days = [1, 10, 15]
    monthly_occurrence {
      occurrence = 1
      day        = "Tuesday"
    }
  }

  depends_on = [azurerm_log_analytics_linked_service.test]
}
`, a.template(data), data.RandomInteger, a.startTime, a.expireTime)
}

func (a SoftwareUpdateConfigurationResource) template(data acceptance.TestData) string {
	// software update needs a log analytic location map correct, a location like `East US` will cause
	// error like `chosen Azure Automation does not have a Log Analytics workspace linked for operation to succeed.`
	// so use a mapping of eastus locations
	// see more https://learn.microsoft.com/en-us/azure/automation/how-to/region-mappings
	loc := data.Locations.Primary
	logLoc := strings.ToLower(loc)
	switch logLoc {
	case "eastus":
		logLoc = "eastus2"
	case "eastus2":
		logLoc = "eastus"
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
  name                = "acctest-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%[1]d"
  location            = "%[3]s"
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_linked_service" "test" {
  resource_group_name = azurerm_resource_group.test.name
  workspace_id        = azurerm_log_analytics_workspace.test.id
  read_access_id      = azurerm_automation_account.test.id
}
`, data.RandomInteger, loc, logLoc)
}
