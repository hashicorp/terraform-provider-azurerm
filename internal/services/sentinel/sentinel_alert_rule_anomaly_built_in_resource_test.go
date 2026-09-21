// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SentinelAlertRuleAnomalyBuiltInResource struct{}

func (r SentinelAlertRuleAnomalyBuiltInResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := securitymlanalyticssettings.ParseSecurityMLAnalyticsSettingID(state.ID)
	if err != nil {
		return nil, err
	}

	client := clients.Sentinel.AnalyticsSettingsClient

	resp, err := client.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

func TestAccSentinelAlertRuleAnomalyBuiltIn_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_built_in", "test")
	r := SentinelAlertRuleAnomalyBuiltInResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").HasValue("8bada072-c58c-4df3-a17e-e02392b48240"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleAnomalyBuiltIn_basicByName(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_built_in", "test")
	r := SentinelAlertRuleAnomalyBuiltInResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicByName(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("display_name").HasValue("UEBA Anomalous Account Deletion"),
			),
		},
		data.ImportStep(),
	})
}

func TestAccSentinelAlertRuleAnomalyBuiltIn_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_built_in", "test")
	r := SentinelAlertRuleAnomalyBuiltInResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
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
	})
}

func (SentinelAlertRuleAnomalyBuiltInResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sentinel_alert_rule_anomaly_built_in" "test" {
  display_name               = "UEBA Anomalous Account Deletion"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  enabled                    = true
  mode                       = "Production"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyBuiltInResource) basicByName(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sentinel_alert_rule_anomaly_built_in" "test" {
  name                       = "8bada072-c58c-4df3-a17e-e02392b48240"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  enabled                    = true
  mode                       = "Production"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyBuiltInResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sentinel_alert_rule_anomaly_built_in" "test" {
  display_name               = "Anomalous scanning activity"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  enabled                    = true
  mode                       = "Flighting"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyBuiltInResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sentinel_alert_rule_anomaly_built_in" "test" {
  display_name               = "Anomalous scanning activity"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  enabled                    = false
  mode                       = "Production"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}
