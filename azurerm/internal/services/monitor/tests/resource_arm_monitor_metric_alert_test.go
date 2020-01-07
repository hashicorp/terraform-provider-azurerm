package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
)

func TestAccAzureRMMonitorMetricAlert_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_requiresImport(t *testing.T) {
	if !features.ShouldResourcesBeImported() {
		t.Skip("Skipping since resources aren't required to be imported")
		return
	}

	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
				),
			},
			{
				Config:      testAccAzureRMMonitorMetricAlert_requiresImport(data),
				ExpectError: acceptance.RequiresImportError("azurerm_monitor_metric_alert"),
			},
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_mitigate", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "severity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is a complete metric alert resource."),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "PT30M"),
					resource.TestCheckResourceAttr(data.ResourceName, "window_size", "PT12H"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_name", "Transactions"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.aggregation", "Maximum"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operator", "Equals"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.threshold", "99"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.name", "GeoType"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.operator", "Include"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.values.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.name", "ApiName"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.operator", "Include"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.values.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.aggregation", "Total"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.operator", "GreaterThanOrEqual"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.threshold", "66.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.dimension.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "2"),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorMetricAlert_basicAndCompleteUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_metric_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorMetricAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_mitigate", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "severity", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "PT1M"),
					resource.TestCheckResourceAttr(data.ResourceName, "window_size", "PT5M"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
			{
				Config: testAccAzureRMMonitorMetricAlert_complete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_mitigate", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "severity", "4"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", "This is a complete metric alert resource."),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "PT30M"),
					resource.TestCheckResourceAttr(data.ResourceName, "window_size", "PT12H"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_name", "Transactions"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.aggregation", "Maximum"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operator", "Equals"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.threshold", "99"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.#", "2"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.name", "GeoType"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.operator", "Include"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.0.values.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.name", "ApiName"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.operator", "Include"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.values.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.1.values.0", "*"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.aggregation", "Total"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.operator", "GreaterThanOrEqual"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.threshold", "66.6"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.1.dimension.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "2"),
				),
			},
			{
				Config: testAccAzureRMMonitorMetricAlert_basic(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorMetricAlertExists(data.ResourceName),
					resource.TestCheckResourceAttr(data.ResourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(data.ResourceName, "auto_mitigate", "false"),
					resource.TestCheckResourceAttr(data.ResourceName, "severity", "3"),
					resource.TestCheckResourceAttr(data.ResourceName, "description", ""),
					resource.TestCheckResourceAttr(data.ResourceName, "frequency", "PT1M"),
					resource.TestCheckResourceAttr(data.ResourceName, "window_size", "PT5M"),
					resource.TestCheckResourceAttr(data.ResourceName, "scopes.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.#", "1"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_namespace", "Microsoft.Storage/storageAccounts"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.metric_name", "UsedCapacity"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.aggregation", "Average"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.operator", "GreaterThan"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.threshold", "55.5"),
					resource.TestCheckResourceAttr(data.ResourceName, "criteria.0.dimension.#", "0"),
					resource.TestCheckResourceAttr(data.ResourceName, "action.#", "0"),
				),
			},
		},
	})
}

func testAccAzureRMMonitorMetricAlert_basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_storage_account.test.id}"]

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString, data.RandomInteger)
}

func testAccAzureRMMonitorMetricAlert_requiresImport(data acceptance.TestData) string {
	template := testAccAzureRMMonitorMetricAlert_basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_metric_alert" "import" {
  name                = "${azurerm_monitor_metric_alert.test.name}"
  resource_group_name = "${azurerm_monitor_metric_alert.test.resource_group_name}"
  scopes              = "${azurerm_monitor_metric_alert.test.scopes}"

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Average"
    operator         = "GreaterThan"
    threshold        = 55.5
  }
}
`, template)
}

func testAccAzureRMMonitorMetricAlert_complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa1%[3]s"
  resource_group_name      = "${azurerm_resource_group.test.name}"
  location                 = "${azurerm_resource_group.test.location}"
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_monitor_action_group" "test1" {
  name                = "acctestActionGroup1-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag1"
}

resource "azurerm_monitor_action_group" "test2" {
  name                = "acctestActionGroup2-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag2"
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestMetricAlert-%[1]d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  scopes              = ["${azurerm_storage_account.test.id}"]
  enabled             = true
  auto_mitigate       = true
  severity            = 4
  description         = "This is a complete metric alert resource."
  frequency           = "PT30M"
  window_size         = "PT12H"

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "Transactions"
    aggregation      = "Maximum"
    operator         = "Equals"
    threshold        = 99

    dimension {
      name     = "GeoType"
      operator = "Include"
      values   = ["*"]
    }

    dimension {
      name     = "ApiName"
      operator = "Include"
      values   = ["*"]
    }
  }

  criteria {
    metric_namespace = "Microsoft.Storage/storageAccounts"
    metric_name      = "UsedCapacity"
    aggregation      = "Total"
    operator         = "GreaterThanOrEqual"
    threshold        = 66.6
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test1.id}"
  }

  action {
    action_group_id = "${azurerm_monitor_action_group.test2.id}"
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomString)
}

func testCheckAzureRMMonitorMetricAlertDestroy(s *terraform.State) error {
	conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.MetricAlertsClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_metric_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Metric alert still exists:\n%#v", resp)
		}
	}
	return nil
}

func testCheckAzureRMMonitorMetricAlertExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.MetricAlertsClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Metric Alert Instance: %s", name)
		}

		resp, err := conn.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorMetricAlertsClient: %+v", err)
		}
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Metric Alert Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
