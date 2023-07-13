// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type MonitorScheduledQueryRulesLogResource struct{}

func TestAccMonitorScheduledQueryRules_LogToMetricActionBasic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_log", "test")
	r := MonitorScheduledQueryRulesLogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.LogToMetricActionConfigBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRules_LogToMetricActionUpdate(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_log", "test")
	r := MonitorScheduledQueryRulesLogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.LogToMetricActionConfigBasic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.LogToMetricActionConfigUpdate(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorScheduledQueryRules_LogToMetricActionComplete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_scheduled_query_rules_log", "test")
	r := MonitorScheduledQueryRulesLogResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.LogToMetricActionConfigComplete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (MonitorScheduledQueryRulesLogResource) LogToMetricActionConfigBasic(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestWorkspace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_scheduled_query_rules_log" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  data_source_id = azurerm_log_analytics_workspace.test.id

  criteria {
    metric_name = "Average_%% Idle Time"
    dimension {
      name     = "InstanceName"
      operator = "Include"
      values   = ["1"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (MonitorScheduledQueryRulesLogResource) LogToMetricActionConfigUpdate(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestWorkspace-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_scheduled_query_rules_log" "test" {
  name                = "acctestsqr-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "test log to metric action"
  enabled             = true

  data_source_id = azurerm_log_analytics_workspace.test.id

  criteria {
    metric_name = "Average_%% Idle Time"
    dimension {
      name     = "InstanceName"
      operator = "Include"
      values   = ["2"]
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger, data.RandomInteger)
}

func (MonitorScheduledQueryRulesLogResource) LogToMetricActionConfigComplete(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-%[1]d"
  location = "%[2]s"
}

resource "azurerm_application_insights" "test" {
  name                = "acctestAppInsights-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestWorkspace-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_monitor_action_group" "test" {
  name                = "acctestActionGroup-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "acctestag"
}

resource "azurerm_monitor_scheduled_query_rules_log" "test" {
  name                = "acctestsqr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  description         = "test log to metric action"
  enabled             = true

  data_source_id          = azurerm_log_analytics_workspace.test.id
  authorized_resource_ids = [azurerm_application_insights.test.id, azurerm_log_analytics_workspace.test.id]

  criteria {
    metric_name = "Average_%% Idle Time"
    dimension {
      name     = "Computer"
      operator = "Include"
      values   = ["*"]
    }
  }
  tags = {
    ENV = "test"
  }
}

resource "azurerm_monitor_metric_alert" "test" {
  name                = "acctestmal-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  scopes              = [azurerm_log_analytics_workspace.test.id]
  description         = "Action will be triggered when Average %% Idle Time is less than 10."

  criteria {
    metric_namespace = "Microsoft.OperationalInsights/workspaces"
    metric_name      = azurerm_monitor_scheduled_query_rules_log.test.criteria[0].metric_name
    aggregation      = "Average"
    operator         = "LessThan"
    threshold        = 10
  }

  action {
    action_group_id = azurerm_monitor_action_group.test.id
  }
}
`, data.RandomInteger, data.Locations.Primary)
}

func (t MonitorScheduledQueryRulesLogResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := clients.Monitor.ScheduledQueryRulesClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("reading (%s): %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}
