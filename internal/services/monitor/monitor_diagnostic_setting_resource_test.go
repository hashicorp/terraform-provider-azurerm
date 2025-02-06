// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorDiagnosticSettingResource struct{}

func TestAccMonitorDiagnosticSetting_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_name").Exists(),
				check.That(data.ResourceName).Key("eventhub_authorization_rule_id").Exists(),
				check.That(data.ResourceName).Key("metric.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_CategoryGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.categoryGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("eventhub_name").Exists(),
				check.That(data.ResourceName).Key("eventhub_authorization_rule_id").Exists(),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("1"),
				check.That(data.ResourceName).Key("metric.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.eventhub(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_diagnostic_setting"),
		},
	})
}

func TestAccMonitorDiagnosticSetting_logAnalyticsWorkspace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalyticsWorkspace(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_analytics_workspace_id").Exists(),
				check.That(data.ResourceName).Key("metric.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalyticsWorkspaceDedicated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_partnerSolution(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.partnerSolution(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("partner_solution_id").Exists(),
				check.That(data.ResourceName).Key("metric.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.storageAccount(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("storage_account_id").Exists(),
				check.That(data.ResourceName).Key("metric.#").HasValue("1"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_activityLog(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.activityLog(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_logAnalyticsDestinationType(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.logAnalyticsDestinationTypeOmit(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_analytics_destination_type").HasValue("AzureDiagnostics"),
			),
		},
		data.ImportStep(),
		{
			Config: r.logAnalyticsDestinationTypeUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("log_analytics_destination_type").HasValue("Dedicated"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_enabledLogsMix(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enabledLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsCategoryGroupUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsCategoryGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_enabledLogsCategoryGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enabledLogsCategoryGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsCategoryGroupUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsCategoryGroup(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorDiagnosticSetting_enabledLogs(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	r := MonitorDiagnosticSettingResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.enabledLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogsUpdated(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("1"),
			),
		},
		data.ImportStep(),
		{
			Config: r.enabledLogs(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("enabled_log.#").HasValue("2"),
			),
		},
	})
}

func (t MonitorDiagnosticSettingResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := monitor.ParseMonitorDiagnosticId(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.DiagnosticSettingsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading diagnostic setting (%s): %+v", id, err)
	}

	return utils.Bool(resp.Model != nil && resp.Model.Id != nil), nil
}

func (MonitorDiagnosticSettingResource) eventhub(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = false
      days    = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) categoryGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  enabled_log {
    category_group = "Audit"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (r MonitorDiagnosticSettingResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_diagnostic_setting" "import" {
  name                           = azurerm_monitor_diagnostic_setting.test.name
  target_resource_id             = azurerm_monitor_diagnostic_setting.test.target_resource_id
  eventhub_authorization_rule_id = azurerm_monitor_diagnostic_setting.test.eventhub_authorization_rule_id
  eventhub_name                  = azurerm_monitor_diagnostic_setting.test.eventhub_name

  metric {
    category = "AllMetrics"

    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, r.eventhub(data))
}

func (MonitorDiagnosticSettingResource) logAnalyticsWorkspace(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[1]d"
  target_resource_id         = azurerm_key_vault.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  metric {
    category = "AllMetrics"

    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) logAnalyticsWorkspaceDedicated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_data_factory" "test" {
  name                = "acctest-DF-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[1]d"
  target_resource_id         = azurerm_data_factory.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  log_analytics_destination_type = "Dedicated"

  enabled_log {
    category = "ActivityRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "PipelineRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "TriggerRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISIntegrationRuntimeLogs"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessageContext"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessages"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutableStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionComponentPhases"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionDataStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxActivityRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxPipelineRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MonitorDiagnosticSettingResource) partnerSolution(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_elastic_cloud_elasticsearch" "test" {
  name                        = "acctest-elastic%[3]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  sku_name                    = "ess-consumption-2024_Monthly"
  elastic_cloud_email_address = "user@example.com"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                = "acctest-DS-%[1]d"
  target_resource_id  = azurerm_key_vault.test.id
  partner_solution_id = azurerm_elastic_cloud_elasticsearch.test.id

  metric {
    category = "AllMetrics"

    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) storageAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_replication_type = "LRS"
  account_tier             = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  target_resource_id = azurerm_key_vault.test.id
  storage_account_id = azurerm_storage_account.test.id

  metric {
    category = "AllMetrics"

    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) activityLog(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}


data "azurerm_subscription" "current" {
}


resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctest%[3]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_replication_type = "LRS"
  account_tier             = "Standard"
}


resource "azurerm_monitor_diagnostic_setting" "test" {
  name               = "acctest-DS-%[1]d"
  target_resource_id = data.azurerm_subscription.current.id
  storage_account_id = azurerm_storage_account.test.id

  enabled_log {
    category = "Administrative"
  }

  enabled_log {
    category = "Alert"
  }

  enabled_log {
    category = "Autoscale"
  }

  enabled_log {
    category = "Policy"
  }

  enabled_log {
    category = "Recommendation"
  }

  enabled_log {
    category = "ResourceHealth"
  }

  enabled_log {
    category = "Security"
  }

  enabled_log {
    category = "ServiceHealth"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) enabledLogs(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  enabled_log {
    category = "AuditEvent"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "AzurePolicyEvaluationDetails"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = false
      days    = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) enabledLogsUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  enabled_log {
    category = "AuditEvent"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = false
      days    = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) enabledLogsCategoryGroup(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  enabled_log {
    category_group = "allLogs"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category_group = "audit"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = false
      days    = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) enabledLogsCategoryGroupUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctest-EHN-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctest-EH-%[1]d"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = azurerm_eventhub_namespace.test.name
  resource_group_name = azurerm_resource_group.test.name
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctest%[3]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  tenant_id           = data.azurerm_client_config.current.tenant_id
  sku_name            = "standard"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctest-DS-%[1]d"
  target_resource_id             = azurerm_key_vault.test.id
  eventhub_authorization_rule_id = azurerm_eventhub_namespace_authorization_rule.test.id
  eventhub_name                  = azurerm_eventhub.test.name

  enabled_log {
    category_group = "allLogs"

    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = true

    retention_policy {
      enabled = false
      days    = 7
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func (MonitorDiagnosticSettingResource) logAnalyticsDestinationTypeOmit(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_data_factory" "test" {
  name                = "acctest-DF-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[1]d"
  target_resource_id         = azurerm_data_factory.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  enabled_log {
    category = "ActivityRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "PipelineRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "TriggerRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISIntegrationRuntimeLogs"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessageContext"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessages"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutableStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionComponentPhases"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionDataStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxActivityRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxPipelineRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = false
    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (MonitorDiagnosticSettingResource) logAnalyticsDestinationTypeUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-LAW-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_data_factory" "test" {
  name                = "acctest-DF-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctest-DS-%[1]d"
  target_resource_id         = azurerm_data_factory.test.id
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id

  log_analytics_destination_type = "Dedicated"

  enabled_log {
    category = "ActivityRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "PipelineRuns"
    retention_policy {
      enabled = false
      days    = 0
    }
  }

  enabled_log {
    category = "TriggerRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISIntegrationRuntimeLogs"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessageContext"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageEventMessages"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutableStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionComponentPhases"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SSISPackageExecutionDataStatistics"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxActivityRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  enabled_log {
    category = "SandboxPipelineRuns"
    retention_policy {
      days    = 0
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled  = false
    retention_policy {
      days    = 0
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}
