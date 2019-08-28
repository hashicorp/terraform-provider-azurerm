package azurerm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMMonitorDiagnosticSetting_eventhub(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostic_setting.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_eventhub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub_name"),
					resource.TestCheckResourceAttrSet(resourceName, "eventhub_authorization_rule_id"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	resourceName := "azurerm_monitor_diagnostic_setting.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_eventhub(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(resourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorDiagnosticSetting_requiresImport(ri, location),
				ExpectError: testRequiresImportError("azurerm_monitor_diagnostic_setting"),
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostic_setting.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "log_analytics_workspace_id"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostic_setting.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "log_analytics_workspace_id"),
					resource.TestCheckResourceAttr(resourceName, "log_analytics_destination_type", "Dedicated"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "log.3188484811.category", "ActivityRuns"),
					resource.TestCheckResourceAttr(resourceName, "log.595859111.category", "PipelineRuns"),
					resource.TestCheckResourceAttr(resourceName, "log.2542277390.category", "TriggerRuns"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.4109484471.category", "AllMetrics"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMMonitorDiagnosticSetting_storageAccount(t *testing.T) {
	resourceName := "azurerm_monitor_diagnostic_setting.test"
	ri := acctest.RandIntRange(10000, 99999)
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMMonitorDiagnosticSettingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorDiagnosticSetting_storageAccount(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorDiagnosticSettingExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "storage_account_id"),
					resource.TestCheckResourceAttr(resourceName, "log.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "log.782743152.category", "AuditEvent"),
					resource.TestCheckResourceAttr(resourceName, "metric.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "metric.1439188313.category", "AllMetrics"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckAzureRMMonitorDiagnosticSettingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %q", resourceName)
		}

		client := testAccProvider.Meta().(*ArmClient).monitor.DiagnosticSettingsClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext
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
	client := testAccProvider.Meta().(*ArmClient).monitor.DiagnosticSettingsClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

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

func testAccAzureRMMonitorDiagnosticSetting_eventhub(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest%d"
  location = "%s"
}

resource "azurerm_eventhub_namespace" "test" {
  name                = "acctesteventhubnamespace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "Basic"
}

resource "azurerm_eventhub" "test" {
  name                = "acctesteventhub-%d"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  partition_count     = 2
  message_retention   = 1
}

resource "azurerm_eventhub_namespace_authorization_rule" "test" {
  name                = "example"
  namespace_name      = "${azurerm_eventhub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  listen              = true
  send                = true
  manage              = true
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                           = "acctestds%d"
  target_resource_id             = "${azurerm_key_vault.test.id}"
  eventhub_authorization_rule_id = "${azurerm_eventhub_namespace_authorization_rule.test.id}"
  eventhub_name                  = "${azurerm_eventhub.test.name}"

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
`, rInt, location, rInt, rInt, rInt, rInt)
}

func testAccAzureRMMonitorDiagnosticSetting_requiresImport(rInt int, location string) string {
	template := testAccAzureRMMonitorDiagnosticSetting_eventhub(rInt, location)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_diagnostic_setting" "import" {
  name                           = "${azurerm_monitor_diagnostic_setting.test.name}"
  target_resource_id             = "${azurerm_monitor_diagnostic_setting.test.target_resource_id}"
  eventhub_authorization_rule_id = "${azurerm_monitor_diagnostic_setting.test.eventhub_authorization_rule_id}"
  eventhub_name                  = "${azurerm_monitor_diagnostic_setting.test.eventhub_name}"

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

func testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspace(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctestds%d"
  target_resource_id         = "${azurerm_key_vault.test.id}"
  log_analytics_workspace_id = "${azurerm_log_analytics_workspace.test.id}"

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
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMonitorDiagnosticSetting_logAnalyticsWorkspaceDedicated(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestlaw%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_data_factory" "test" {
  name                = "acctestdf%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name                       = "acctestds%d"
  target_resource_id         = "${azurerm_data_factory.test.id}"
  log_analytics_workspace_id = "${azurerm_log_analytics_workspace.test.id}"

  log_analytics_destination_type = "Dedicated"

  log {
    category = "ActivityRuns"

    retention_policy {
      enabled = false
    }
  }

  log {
    category = "PipelineRuns"
    enabled = false

    retention_policy {
      enabled = false
    }
  }

  log {
    category = "TriggerRuns"
    enabled = false

    retention_policy {
      enabled = false
    }
  }

  metric {
    category = "AllMetrics"
    enabled = false

    retention_policy {
      enabled = false
    }
  }
}
`, rInt, location, rInt, rInt, rInt)
}

func testAccAzureRMMonitorDiagnosticSetting_storageAccount(rInt int, location string) string {
	return fmt.Sprintf(`
data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "test" {
  name     = "acctest%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestlogs%d"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_replication_type = "LRS"
  account_tier             = "Standard"
}

resource "azurerm_key_vault" "test" {
  name                = "acctestkv%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  sku {
    name = "standard"
  }
}

resource "azurerm_monitor_diagnostic_setting" "test" {
  name               = "acctestds%d"
  target_resource_id = "${azurerm_key_vault.test.id}"
  storage_account_id = "${azurerm_storage_account.test.id}"

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
`, rInt, location, rInt, rInt, rInt)
}
