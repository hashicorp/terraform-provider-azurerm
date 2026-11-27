// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_machine_learning_behavior_analytics

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleMLBehaviorAnalyticsRead(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).AlertRulesClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Sentinel Alert Rule MLBehaviorAnalytics %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
	}

	if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMLBehaviorAnalytics); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}
	if model := resp.Model; model != nil {
		if rule, ok := model.(alertrules.MLBehaviorAnalyticsAlertRule); ok {
			d.Set("name", id.RuleId)

			workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			d.Set("log_analytics_workspace_id", workspaceId.ID())

			if prop := rule.Properties; prop != nil {
				d.Set("enabled", prop.Enabled)
				d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
			}
		}
	}

	return nil
}
