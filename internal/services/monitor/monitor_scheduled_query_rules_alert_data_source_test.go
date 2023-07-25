// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type MonitorScheduledQueryRulesDataSource struct{}

func TestAccDataSourceMonitorScheduledQueryRules_AlertingAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.AlertingActionConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func TestAccDataSourceMonitorScheduledQueryRules_AlertingActionCrossResource(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_monitor_scheduled_query_rules_alert", "test")
	r := MonitorScheduledQueryRulesDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.AlertingActionCrossResourceConfig(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
			),
		},
	})
}

func (MonitorScheduledQueryRulesDataSource) AlertingActionConfig(data acceptance.TestData) string {
	ts := time.Now().Format(time.RFC3339)

	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, MonitorScheduledQueryRulesResource{}.AlertingActionConfigBasic(data, ts))
}

func (MonitorScheduledQueryRulesDataSource) AlertingActionCrossResourceConfig(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_monitor_scheduled_query_rules_alert" "test" {
  name                = basename(azurerm_monitor_scheduled_query_rules_alert.test.id)
  resource_group_name = "${azurerm_resource_group.test.name}"
}
`, MonitorScheduledQueryRulesResource{}.AlertingActionCrossResourceConfig(data))
}
