package monitor_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/acceptance/check"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type MonitorScheduledQueryRulesResource struct {
}

func TestAccMonitorScheduledQueryRules_AlertingActionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesResource{}
	ts := time.Now().Format(time.RFC3339)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AlertingActionConfigBasic(data, ts),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRules_AlertingActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesResource{}
	ts := time.Now().Format(time.RFC3339)

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AlertingActionConfigBasic(data, ts),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.AlertingActionConfigUpdate(data, ts),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRules_AlertingActionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AlertingActionConfigComplete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRules_AlertingActionCrossResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.AlertingActionCrossResourceConfig(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MonitorScheduledQueryRulesResource) AlertingActionConfigBasic(data acceptance.TestData, ts string) string {
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

func (MonitorScheduledQueryRulesResource) AlertingActionConfigUpdate(data acceptance.TestData, ts string) string {
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

func (MonitorScheduledQueryRulesResource) AlertingActionConfigComplete(data acceptance.TestData) string {
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

func (MonitorScheduledQueryRulesResource) AlertingActionCrossResourceConfig(data acceptance.TestData) string {
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

func (t MonitorScheduledQueryRulesResource) Exists(ctx context.Context, clients *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := azure.ParseAzureResourceID(state.ID)
	if err != nil {
		return nil, err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["scheduledqueryrules"]

	resp, err := clients.Monitor.ScheduledQueryRulesClient.Get(ctx, resourceGroup, name)
	if err != nil {
		return nil, fmt.Errorf("reading Scheduled Query Rules (%s): %+v", id, err)
	}

	return utils.Bool(resp.ID != nil), nil
}
