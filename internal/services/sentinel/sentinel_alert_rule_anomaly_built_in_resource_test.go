// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel_test

import (
	"context"
	"fmt"
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

type SentinelAlertRuleAnomalyBuiltInResource struct{}

func (r SentinelAlertRuleAnomalyBuiltInResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
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
	return utils.Bool(resp != nil && resp.Enabled != nil && *resp.Enabled == true), nil
}

func TestAccSentinelAlertRuleAnomalyBuiltIn_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_sentinel_alert_rule_anomaly_built_in", "test")
	r := SentinelAlertRuleAnomalyBuiltInResource{}

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

func (SentinelAlertRuleAnomalyBuiltInResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
resource "azurerm_sentinel_alert_rule_anomaly_built_in" "test" {
  display_name               = "Potential data staging"
  log_analytics_workspace_id = azurerm_sentinel_log_analytics_workspace_onboarding.test.workspace_id
  enabled                    = true
  mode                       = "Production"
}
`, SecurityInsightsSentinelOnboardingStateResource{}.basic(data))
}
