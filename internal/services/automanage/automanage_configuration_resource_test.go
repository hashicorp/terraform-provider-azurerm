// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AutoManageConfigurationProfileResource struct{}

func TestAccAutoManageConfigurationProfile_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
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

func TestAccAutoManageConfigurationProfile_antimalware(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.antimalware(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("0"),
				check.That(data.ResourceName).Key("automation_account_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_azureSecurityBaseline(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.azureSecurityBaseline(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_security_baseline.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_security_baseline.0.assignment_type").HasValue("ApplyAndAutoCorrect"),
			),
		},
		data.ImportStep(),
		{
			Config: r.azureSecurityBaselineUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_security_baseline.#").HasValue("1"),
				check.That(data.ResourceName).Key("azure_security_baseline.0.assignment_type").HasValue("DeployAndAutoCorrect"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("azure_security_baseline.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccAutoManageConfigurationProfile_backup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.backup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.policy_name").Exists(),
				check.That(data.ResourceName).Key("backup.0.time_zone").HasValue("UTC"),
				check.That(data.ResourceName).Key("backup.0.instant_rp_retention_range_in_days").HasValue("2"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_days.#").HasValue("2"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_days.0").HasValue("Monday"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_days.1").HasValue("Tuesday"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_times.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_times.0").HasValue("12:00"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_policy_type").HasValue("SimpleSchedulePolicy"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.retention_policy_type").HasValue("LongTermRetentionPolicy"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_times.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_times.0").HasValue("12:00"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.0.duration_type").HasValue("Days"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.0.retention_times.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.0.retention_times.0").HasValue("14:00"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.0.retention_duration.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.0.retention_duration.0.count").HasValue("4"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.0.retention_duration.0.duration_type").HasValue("Weeks"),
			),
		},
		data.ImportStep(),
		{
			Config: r.backupUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.policy_name").Exists(),
				check.That(data.ResourceName).Key("backup.0.time_zone").HasValue("UTC"),
				check.That(data.ResourceName).Key("backup.0.instant_rp_retention_range_in_days").HasValue("5"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_frequency").HasValue("Daily"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_days.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_days.0").HasValue("Monday"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_times.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.schedule_policy.0.schedule_run_times.0").HasValue("12:00"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.retention_policy_type").HasValue("LongTermRetentionPolicy"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_times.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_times.0").HasValue("12:00"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.#").HasValue("1"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.0.count").HasValue("7"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.daily_schedule.0.retention_duration.0.duration_type").HasValue("Days"),
				check.That(data.ResourceName).Key("backup.0.retention_policy.0.weekly_schedule.#").HasValue("0"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("backup.#").HasValue("0"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_logAnalytics(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalytics(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_analytics_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_analytics_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.0.extensions").HasValue("exe;dll"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.0.processes").HasValue("svchost.exe;notepad.exe"),
				check.That(data.ResourceName).Key("antimalware.0.real_time_protection_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_day").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_type").HasValue("Quick"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_time_in_minutes").HasValue("1339"),
				check.That(data.ResourceName).Key("automation_account_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("boot_diagnostics_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("defender_for_cloud_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("guest_configuration_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("status_change_alert_enabled").HasValue("true"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccAutoManageConfigurationProfile_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_automanage_configuration", "test")
	r := AutoManageConfigurationProfileResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("antimalware.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.#").HasValue("1"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.0.extensions").HasValue("exe"),
				check.That(data.ResourceName).Key("antimalware.0.exclusions.0.processes").HasValue("svchost.exe"),
				check.That(data.ResourceName).Key("antimalware.0.real_time_protection_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_enabled").HasValue("true"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_day").HasValue("2"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_type").HasValue("Full"),
				check.That(data.ResourceName).Key("antimalware.0.scheduled_scan_time_in_minutes").HasValue("1338"),
				check.That(data.ResourceName).Key("automation_account_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("boot_diagnostics_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("defender_for_cloud_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("guest_configuration_enabled").HasValue("false"),
				check.That(data.ResourceName).Key("status_change_alert_enabled").HasValue("false"),
			),
		},
		data.ImportStep(),
	})
}

func (r AutoManageConfigurationProfileResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	client := clients.Automanage.ConfigurationProfilesClient

	id, err := configurationprofiles.ParseConfigurationProfileID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	return pointer.To(resp.Model != nil), nil
}

func (r AutoManageConfigurationProfileResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctest-rg-%d"
  location = "%s"
}
`, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) antimalware(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  antimalware {
    exclusions {
      extensions = "exe;dll"
    }
    real_time_protection_enabled = true
  }
  automation_account_enabled = true
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) logAnalytics(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
				%s

resource "azurerm_automanage_configuration" "test" {
  name                  = "acctest-amcp-%d"
  resource_group_name   = azurerm_resource_group.test.name
  location              = "%s"
  log_analytics_enabled = true
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) requiresImport(data acceptance.TestData) string {
	config := r.antimalware(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "import" {
  name                = azurerm_automanage_configuration.test.name
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
}
`, config, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  antimalware {
    exclusions {
      extensions = "exe;dll"
      paths      = "C:\\Windows\\Temp;D:\\Temp"
      processes  = "svchost.exe;notepad.exe"
    }
    real_time_protection_enabled   = true
    scheduled_scan_enabled         = true
    scheduled_scan_type            = "Quick"
    scheduled_scan_day             = 1
    scheduled_scan_time_in_minutes = 1339
  }
  automation_account_enabled  = true
  boot_diagnostics_enabled    = true
  defender_for_cloud_enabled  = true
  guest_configuration_enabled = true
  status_change_alert_enabled = true
  tags = {
    "env" = "test"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) update(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = "%s"
  antimalware {
    exclusions {
      extensions = "exe"
      processes  = "svchost.exe"
    }
    real_time_protection_enabled   = false
    scheduled_scan_enabled         = true
    scheduled_scan_type            = "Full"
    scheduled_scan_day             = 2
    scheduled_scan_time_in_minutes = 1338
  }
  tags = {
    "env2" = "test2"
  }
}
`, template, data.RandomInteger, data.Locations.Primary)
}

func (r AutoManageConfigurationProfileResource) azureSecurityBaseline(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  azure_security_baseline {
    assignment_type = "ApplyAndAutoCorrect"
  }
}
`, template, data.RandomInteger)
}

func (r AutoManageConfigurationProfileResource) azureSecurityBaselineUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  azure_security_baseline {
    assignment_type = "DeployAndAutoCorrect"
  }
}
`, template, data.RandomInteger)
}

func (r AutoManageConfigurationProfileResource) backup(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  backup {
    policy_name                        = "acctest-backup-policy-%d"
    time_zone                          = "UTC"
    instant_rp_retention_range_in_days = 2

    schedule_policy {
      schedule_run_frequency = "Daily"
      schedule_run_days      = ["Monday", "Tuesday"]
      schedule_run_times     = ["12:00"]
      schedule_policy_type   = "SimpleSchedulePolicy"
    }

    retention_policy {
      retention_policy_type = "LongTermRetentionPolicy"

      daily_schedule {
        retention_times = ["12:00"]
        retention_duration {
          count         = 7
          duration_type = "Days"
        }
      }

      weekly_schedule {
        retention_times = ["14:00"]
        retention_duration {
          count         = 4
          duration_type = "Weeks"
        }
      }
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}

func (r AutoManageConfigurationProfileResource) backupUpdate(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
			%s

resource "azurerm_automanage_configuration" "test" {
  name                = "acctest-amcp-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  backup {
    policy_name                        = "acctest-backup-policy-%d"
    time_zone                          = "UTC"
    instant_rp_retention_range_in_days = 5

    schedule_policy {
      schedule_run_frequency = "Daily"
      schedule_run_days      = ["Monday"]
      schedule_run_times     = ["12:00"]
    }

    retention_policy {
      retention_policy_type = "LongTermRetentionPolicy"

      daily_schedule {
        retention_times = ["12:00"]
        retention_duration {
          count         = 7
          duration_type = "Days"
        }
      }
    }
  }
}
`, template, data.RandomInteger, data.RandomInteger)
}
