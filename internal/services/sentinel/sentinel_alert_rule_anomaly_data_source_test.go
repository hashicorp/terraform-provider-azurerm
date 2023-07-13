// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
)

type SentinelAlertRuleAnomalyDataSource struct{}

func TestAccSentinelAlertRuleAnomalyDataSource_basicWithThreshold(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_anomaly", "test")
	r := SentinelAlertRuleAnomalyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic_withThreshold(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("anomaly_version").Exists(),
				check.That(data.ResourceName).Key("anomaly_settings_version").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
				check.That(data.ResourceName).Key("frequency").Exists(),
				check.That(data.ResourceName).Key("required_data_connector.#").HasValue("1"),
				check.That(data.ResourceName).Key("mode").Exists(),
				check.That(data.ResourceName).Key("settings_definition_id").Exists(),
				check.That(data.ResourceName).Key("tactics.#").HasValue("1"),
				check.That(data.ResourceName).Key("techniques.#").HasValue("1"),
				check.That(data.ResourceName).Key("threshold_observation.#").HasValue("2"),
			),
		},
	})
}

func TestAccSentinelAlertRuleAnomalyDataSource_basicWithSingleSelect(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_anomaly", "test")
	r := SentinelAlertRuleAnomalyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic_withSingleSelect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("anomaly_version").Exists(),
				check.That(data.ResourceName).Key("anomaly_settings_version").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
				check.That(data.ResourceName).Key("frequency").Exists(),
				check.That(data.ResourceName).Key("required_data_connector.#").HasValue("1"),
				check.That(data.ResourceName).Key("mode").Exists(),
				check.That(data.ResourceName).Key("settings_definition_id").Exists(),
				check.That(data.ResourceName).Key("tactics.#").HasValue("1"),
				check.That(data.ResourceName).Key("techniques.#").HasValue("2"),
				check.That(data.ResourceName).Key("single_select_observation.#").HasValue("2"),
			),
		},
	})
}

func TestAccSentinelAlertRuleAnomalyDataSource_basicWithMultiSelect(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_anomaly", "test")
	r := SentinelAlertRuleAnomalyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic_withMultiSelect(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("anomaly_version").Exists(),
				check.That(data.ResourceName).Key("anomaly_settings_version").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
				check.That(data.ResourceName).Key("frequency").Exists(),
				check.That(data.ResourceName).Key("required_data_connector.#").HasValue("1"),
				check.That(data.ResourceName).Key("mode").Exists(),
				check.That(data.ResourceName).Key("settings_definition_id").Exists(),
				check.That(data.ResourceName).Key("tactics.#").HasValue("1"),
				check.That(data.ResourceName).Key("techniques.#").HasValue("1"),
				check.That(data.ResourceName).Key("multi_select_observation.#").HasValue("1"),
			),
		},
	})
}

func TestAccSentinelAlertRuleAnomalyDataSource_basicWithPrioritized(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azurerm_sentinel_alert_rule_anomaly", "test")
	r := SentinelAlertRuleAnomalyDataSource{}

	data.DataSourceTest(t, []acceptance.TestStep{
		{
			Config: r.basic_withPrioritizeExclude(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
				check.That(data.ResourceName).Key("display_name").Exists(),
				check.That(data.ResourceName).Key("anomaly_version").Exists(),
				check.That(data.ResourceName).Key("anomaly_settings_version").Exists(),
				check.That(data.ResourceName).Key("description").Exists(),
				check.That(data.ResourceName).Key("enabled").Exists(),
				check.That(data.ResourceName).Key("frequency").Exists(),
				check.That(data.ResourceName).Key("required_data_connector.#").HasValue("1"),
				check.That(data.ResourceName).Key("mode").Exists(),
				check.That(data.ResourceName).Key("settings_definition_id").Exists(),
				check.That(data.ResourceName).Key("tactics.#").HasValue("2"),
				check.That(data.ResourceName).Key("techniques.#").HasValue("2"),
				check.That(data.ResourceName).Key("prioritized_exclude_observation.#").HasValue("2"),
			),
		},
	})
}

func (SentinelAlertRuleAnomalyDataSource) basic_withThreshold(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Potential data staging"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDataSource) basic_withSingleSelect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Suspicious geography change in Palo Alto GlobalProtect account logins"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDataSource) basic_withMultiSelect(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Attempted user account bruteforce per logon type"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDataSource) basic_withPrioritizeExclude(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Anomalous web request activity"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}
