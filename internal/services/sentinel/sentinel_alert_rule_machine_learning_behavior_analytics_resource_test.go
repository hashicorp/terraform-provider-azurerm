// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelAlertRuleMLBehaviorAnalyticsResource struct{}

func TestAccSentinelAlertRuleMLBehaviorAnalytics_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", "test")
	r := SentinelAlertRuleMLBehaviorAnalyticsResource{}

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

func TestAccSentinelAlertRuleMLBehaviorAnalytics_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", "test")
	r := SentinelAlertRuleMLBehaviorAnalyticsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleMLBehaviorAnalytics_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", "test")
	r := SentinelAlertRuleMLBehaviorAnalyticsResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleMLBehaviorAnalytics_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", "test")
	r := SentinelAlertRuleMLBehaviorAnalyticsResource{}

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

func (r SentinelAlertRuleMLBehaviorAnalyticsResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	alertRuleClient := client.Sentinel.AlertRulesClient
	id, err := alertrules.ParseAlertRuleID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := alertRuleClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving Sentinel Alert Rule MLBehaviorAnalytics (%q): %+v", state.String(), err)
	}

	if model := resp.Model; model != nil {
		rule, ok := model.(alertrules.MLBehaviorAnalyticsAlertRule)
		if !ok {
			return nil, fmt.Errorf("the Alert Rule %q is not a Fusion Alert Rule", id)
		}
		return utils.Bool(rule.Id != nil), nil
	}

	return utils.Bool(false), nil
}

func (r SentinelAlertRuleMLBehaviorAnalyticsResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "(Preview) Anomalous SSH Login Detection"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
}

resource "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics" "test" {
  name                       = "acctest-SentinelAlertRule-MLBehaviorAnalytics-%d"
  log_analytics_workspace_id = azurerm_log_analytics_workspace.test.id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name

}
`, r.template(data), data.RandomInteger)
}

func (r SentinelAlertRuleMLBehaviorAnalyticsResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_template" "test" {
  display_name               = "(Preview) Anomalous SSH Login Detection"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
}

resource "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics" "test" {
  name                       = "acctest-SentinelAlertRule-MLBehaviorAnalytics-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test.name
  enabled                    = false
}

data "azurerm_sentinel_alert_rule_template" "test2" {
  display_name               = "(Preview) Anomalous RDP Login Detections"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
}

resource "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics" "test2" {
  name                       = "acctest-SentinelAlertRule-MLBehaviorAnalytics-2-%d"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  alert_rule_template_guid   = data.azurerm_sentinel_alert_rule_template.test2.name
  enabled                    = false
}
`, r.template(data), data.RandomInteger, data.RandomInteger)
}

func (r SentinelAlertRuleMLBehaviorAnalyticsResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_machine_learning_behavior_analytics" "import" {
  name                       = azurerm_sentinel_alert_rule_machine_learning_behavior_analytics.test.name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_machine_learning_behavior_analytics.test.log_analytics_workspace_id
  alert_rule_template_guid   = azurerm_sentinel_alert_rule_machine_learning_behavior_analytics.test.alert_rule_template_guid
}
`, r.basic(data))
}

func (r SentinelAlertRuleMLBehaviorAnalyticsResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-sentinel-%d"
  location = "%s"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctestLAW-%d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
}

resource "azurerm_sentinel_log_analytics_workspace_onboarding" "test" {
  workspace_id = azurerm_log_analytics_workspace.test.id
}
`, data.RandomInteger, data.Locations.Primary, data.RandomInteger)
}
