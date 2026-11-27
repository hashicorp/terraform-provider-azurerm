// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleFusionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	workspaceId := alertrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	d.Set("log_analytics_workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindFusion); err != nil {
			return fmt.Errorf("asserting alert rule of %s: %+v", id, err)
		}

		if rule, ok := model.(alertrules.FusionAlertRule); ok {
			if prop := rule.Properties; prop != nil {
				d.Set("enabled", prop.Enabled)
				d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
				if err := d.Set("source", flattenFusionSourceSettings(prop.SourceSettings)); err != nil {
					return fmt.Errorf("setting `source`: %v", err)
				}
			}
		}
	}

	return nil
}

const SentinelAlertRuleFusionName = "BuiltInFusion"
