package tests

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TestAccAzureRMMonitorScheduledQueryRules_AlertingActionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	ts := time.Now().Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRules_AlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigBasic(data, ts),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRules_AlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorScheduledQueryRules_AlertingActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	ts := time.Now().Format(time.RFC3339)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRules_AlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigBasic(data, ts),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRules_AlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigUpdate(data, ts),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRules_AlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorScheduledQueryRules_AlertingActionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRules_AlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigComplete(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRules_AlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func TestAccAzureRMMonitorScheduledQueryRules_AlertingActionCrossResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { acceptance.PreCheck(t) },
		Providers:    acceptance.SupportedProviders,
		CheckDestroy: testCheckAzureRMMonitorScheduledQueryRules_AlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMMonitorScheduledQueryRules_AlertingActionCrossResourceConfig(data),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMMonitorScheduledQueryRules_AlertExists(data.ResourceName),
				),
			},
			data.ImportStep(),
		},
	})
}

func testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigBasic(data acceptance.TestData, ts string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  data_source_id = azurerm_application_insights.test.id
  query          = <<-QUERY
	let d=datatable(TimeGenerated: datetime, usage_percent: double) [  '%s', 25.4, '%s', 75.4 ];
	d | summarize AggregatedValue=avg(usage_percent) by bin(TimeGenerated, 1h)
QUERY


  frequency   = 60
  time_window = 60

  action {
    action_group = [azurerm_monitor_action_group.test.id]
  }

  trigger {
    operator  = "GreaterThan"
    threshold = 5000
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, ts, ts)
}

func testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigUpdate(data acceptance.TestData, ts string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  data_source_id = azurerm_application_insights.test.id
  query          = <<-QUERY
	let d=datatable(TimeGenerated: datetime, usage_percent: double) [  '%s', 25.4, '%s', 75.4 ];
	d | summarize AggregatedValue=avg(usage_percent) by bin(TimeGenerated, 1h)
QUERY


  enabled     = false
  description = "test description"

  frequency   = 30
  time_window = 30

  action {
    action_group = [azurerm_monitor_action_group.test.id]
  }

  trigger {
    operator  = "GreaterThan"
    threshold = 1000
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, ts, ts)
}

func testAccAzureRMMonitorScheduledQueryRules_AlertingActionConfigComplete(data acceptance.TestData) string {
	ts := time.Now().Format(time.RFC3339)

	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "test alerting action"
  enabled             = true

  data_source_id = azurerm_application_insights.test.id
  query          = "let d=datatable(TimeGenerated: datetime, usage_percent: double) [  '%s', 25.4, '%s', 75.4 ]; d | summarize AggregatedValue=avg(usage_percent) by bin(TimeGenerated, 1h)"

  frequency   = 60
  time_window = 60

  severity   = 3
  throttling = 5
  action {
    action_group           = [azurerm_monitor_action_group.test.id]
    email_subject          = "Custom alert email subject"
    custom_webhook_payload = "{}"
  }

  trigger {
    operator  = "GreaterThan"
    threshold = 5000
    metric_trigger {
      operator            = "GreaterThan"
      threshold           = 1
      metric_trigger_type = "Total"
      metric_column       = "TimeGenerated"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, ts, ts)
}

func testAccAzureRMMonitorScheduledQueryRules_AlertingActionCrossResourceConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  application_type    = "web"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestWorkspace-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  description         = "test alerting action cross-resource"
  enabled             = true

  authorized_resource_ids = ["${azurerm_application_insights.test.id}", "${azurerm_log_analytics_workspace.test.id}"]
  data_source_id          = "${azurerm_application_insights.test.id}"
  query = format(<<-QUERY
	let a=workspace('%%s').Perf
		| where Computer='dependency' and TimeGenerated > ago(1h)
		| where ObjectName == 'Processor' and CounterName == '%%%% Processor Time'
		| summarize cpu=avg(CounterValue) by bin(TimeGenerated, 1m)
		| extend ts=tostring(TimeGenerated); let b=requests
		| where resultCode == '200' and timestamp > ago(1h)
		| summarize reqs=count() by bin(timestamp, 1m)
		| extend ts = tostring(timestamp); a
		| join b on $left.ts == $right.ts
		| where cpu > 50 and reqs > 5
QUERY
  , azurerm_log_analytics_workspace.test.id)

  frequency   = 60
  time_window = 60

  severity = 3
  action {
    action_group  = ["${azurerm_monitor_action_group.test.id}"]
    email_subject = "Custom alert email subject"
  }

  trigger {
    operator  = "GreaterThan"
    threshold = 5000
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger, data.RandomInteger, data.RandomInteger)
}

func testCheckAzureRMMonitorScheduledQueryRules_AlertDestroy(s *terraform.State) error {
	client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_monitor_scheduled_query_rules_alert" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return nil
		}

		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("Scheduled Query Rule still exists:\n%#v", resp)
		}
	}

	return nil
}

func testCheckAzureRMMonitorScheduledQueryRules_AlertExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Ensure we have enough information in state to look up in API
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup, hasResourceGroup := rs.Primary.Attributes["resource_group_name"]
		if !hasResourceGroup {
			return fmt.Errorf("Bad: no resource group found in state for Scheduled Query Rule Instance: %s", name)
		}

		client := acceptance.AzureProvider.Meta().(*clients.Client).Monitor.ScheduledQueryRulesClient
		ctx := acceptance.AzureProvider.Meta().(*clients.Client).StopContext

		resp, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			return fmt.Errorf("Bad: Get on monitorScheduledQueryRulesClient: %+v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Scheduled Query Rule Instance %q (resource group: %q) does not exist", name, resourceGroup)
		}

		return nil
	}
}
