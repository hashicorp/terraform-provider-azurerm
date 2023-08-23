// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2022-10-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SentinelAlertRuleAnomalyDuplicateResource struct{}

func (r SentinelAlertRuleAnomalyDuplicateResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := parse.MLAnalyticsSettingsID(state.ID)
	if err != nil {
		return nil, err
	}

	workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName)
	client := clients.Sentinel.AnalyticsSettingsClient
	resp, err := sentinel.AlertRuleAnomalyReadWithPredicate(ctx, client.BaseClient, workspaceId, func(r *azuresdkhacks.AnomalySecurityMLAnalyticsSettings) bool {
		if r.Name != nil && strings.EqualFold(sentinel.AlertRuleAnomalyIdFromWorkspaceId(workspaceId, *r.Name), id.ID()) {
			return true
		}
		return false
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving Sentinel Alert Rule Anomaly Built In %q (Workspace %q / Resource Group %q): %+v", id.SecurityMLAnalyticsSettingName, id.WorkspaceName, id.ResourceGroup, err)
	}
	return utils.Bool(resp != nil), nil
}

func TestAccSentinelAlertRuleAnomalyDuplicate_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_duplicate", "test")
	r := SentinelAlertRuleAnomalyDuplicateResource{}

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

func TestAccSentinelAlertRuleAnomalyDuplicate_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_duplicate", "test")
	r := SentinelAlertRuleAnomalyDuplicateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: regexp.MustCompile("only one duplicate rule of the same built-in rule is allowed, there is an existing duplicate rule of .+"),
		},
	})
}

func TestAccSentinelAlertRuleAnomalyDuplicate_withCustomObservation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_duplicate", "test")
	r := SentinelAlertRuleAnomalyDuplicateResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basicWithThresholdObservation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithMultiSelectObservation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithSingleSelectObservation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.basicWithPrioritizeExcludeObservation(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (SentinelAlertRuleAnomalyDuplicateResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Potential data staging"
}

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "test" {
  display_name               = "acctest duplicate rule"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  built_in_rule_id           = data.azurerm_sentinel_alert_rule_anomaly.test.id
  enabled                    = true
  mode                       = "Flighting"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDuplicateResource) basicWithThresholdObservation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "UEBA Anomalous Sign In"
}

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "test" {
  display_name               = "acctest duplicate rule"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  built_in_rule_id           = data.azurerm_sentinel_alert_rule_anomaly.test.id
  enabled                    = true
  mode                       = "Flighting"

  threshold_observation {
    name  = "Anomaly score threshold"
    value = "0.6"
  }

}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDuplicateResource) basicWithSingleSelectObservation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Unusual web traffic detected with IP in URL path"
}

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "test" {
  display_name               = "acctest duplicate Unusual web traffic detected with IP in URL path"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  built_in_rule_id           = data.azurerm_sentinel_alert_rule_anomaly.test.id
  enabled                    = true
  mode                       = "Flighting"

  single_select_observation {
    name  = "Device vendor"
    value = "Zscaler"
  }
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (SentinelAlertRuleAnomalyDuplicateResource) basicWithMultiSelectObservation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Anomalous scanning activity"
}

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "test" {
  display_name               = "acctest duplicate Anomalous scanning activity"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  built_in_rule_id           = data.azurerm_sentinel_alert_rule_anomaly.test.id
  enabled                    = true
  mode                       = "Flighting"

  multi_select_observation {
    name   = "Device action"
    values = ["accept"]
  }
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}
func (SentinelAlertRuleAnomalyDuplicateResource) basicWithPrioritizeExcludeObservation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm_sentinel_alert_rule_anomaly" "test" {
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  display_name               = "Anomalous web request activity"
}

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "test" {
  display_name               = "acctest duplicate Anomalous web request activity"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  built_in_rule_id           = data.azurerm_sentinel_alert_rule_anomaly.test.id
  enabled                    = true
  mode                       = "Flighting"

  prioritized_exclude_observation {
    name       = "Prioritize script suffixes of the URI stems"
    prioritize = ".asp, .aspx, .armx, .asax, .ashz"
  }

  prioritized_exclude_observation {
    name    = "Exclude noisy URI stems"
    exclude = "test.com"
  }

}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}

func (r SentinelAlertRuleAnomalyDuplicateResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_sentinel_alert_rule_anomaly_duplicate" "import" {
  display_name               = azurerm_sentinel_alert_rule_anomaly_duplicate.test.display_name
  log_analytics_workspace_id = azurerm_sentinel_alert_rule_anomaly_duplicate.test.log_analytics_workspace_id
  built_in_rule_id           = azurerm_sentinel_alert_rule_anomaly_duplicate.test.built_in_rule_id
  enabled                    = azurerm_sentinel_alert_rule_anomaly_duplicate.test.enabled
  mode                       = azurerm_sentinel_alert_rule_anomaly_duplicate.test.mode
  depends_on                 = [azurerm_sentinel_log_analytics_workspace_onboarding.test]
}
`, r.basic(data))
}
