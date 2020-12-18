package monitor_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMonitorDiagnosticSetting_eventhub(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_eventhub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub_name"),
					resource.TestCheckResourceAttrSet(data.ResourceName, "eventhub_authorization_rule_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_eventhub(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorDiagnosticSetting_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_diagnostic_setting"),
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "log_analytics_workspace_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_storageAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_storageAccount(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
					resource.TestCheckResourceAttrSet(data.ResourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_activityLog(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_diagnostic_setting", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_activityLog(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testCheckAzureRMMonitorDiagnosticSettingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.DiagnosticSettingsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		actualResourceId := rs.Primary.Attributes["target_resource_id"]
		targetResourceId := strings.TrimPrefix(actualResourceId, "/")

		resp, err := client.Get(ctx, targetResourceId, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Monitor Diagnostic Setting %q does not exist for Resource ID %s", name, targetResourceId)
			}

			return fmt.Errorf("Bad: Get on monitorDiagnosticSettingsClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMMonitorDiagnosticSettingDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.DiagnosticSettingsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_diagnostic_setting" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		actualResourceId := rs.Primary.Attributes["target_resource_id"]
		targetResourceId := strings.TrimPrefix(actualResourceId, "/")

		resp, err := client.Get(ctx, targetResourceId, name)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil
			}

			return err
		}
	}

	return nil
}

func testAccAzureRMMonitorDiagnosticSetting_eventhub(data acceptance.TestData) string {
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

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func testAccAzureRMMonitorDiagnosticSetting_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorDiagnosticSetting_eventhub(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_diagnostic_setting" "import" {
  name                           = azurerm_monitor_diagnostic_setting.test.name
  target_resource_id             = azurerm_monitor_diagnostic_setting.test.target_resource_id
  eventhub_authorization_rule_id = azurerm_monitor_diagnostic_setting.test.eventhub_authorization_rule_id
  eventhub_name                  = azurerm_monitor_diagnostic_setting.test.eventhub_name

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
`, template)
}

func testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(data acceptance.TestData) string {
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

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(data acceptance.TestData) string {
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

  log {
    category = "ActivityRuns"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "PipelineRuns"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "TriggerRuns"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISIntegrationRuntimeLogs"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISPackageEventMessageContext"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISPackageEventMessages"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISPackageExecutableStatistics"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISPackageExecutionComponentPhases"
    retention_policy {
      enabled = false
    }
  }

  log {
    category = "SSISPackageExecutionDataStatistics"
    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    retention_policy {
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func testAccAzureRMMonitorDiagnosticSetting_storageAccount(data acceptance.TestData) string {
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

  log {
    category = "AuditEvent"
    enabled  = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"

    retention_policy {
      enabled = false
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}

func testAccAzureRMMonitorDiagnosticSetting_activityLog(data acceptance.TestData) string {
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

  log {
    category = "Administrative"
    enabled  = true
  }

  log {
    category = "Alert"
    enabled  = true
  }

  log {
    category = "Autoscale"
    enabled  = true
  }

  log {
    category = "Policy"
    enabled  = true
  }

  log {
    category = "Recommendation"
    enabled  = true
  }

  log {
    category = "ResourceHealth"
    enabled  = true
  }

  log {
    category = "Security"
    enabled  = true
  }

  log {
    category = "ServiceHealth"
    enabled  = true
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomIntOfLength(17))
}
