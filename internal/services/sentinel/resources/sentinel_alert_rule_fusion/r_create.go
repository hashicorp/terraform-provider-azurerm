// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleFusionCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := SentinelAlertRuleFusionName

	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	// The only one fusion alert is enabled by default, so we do not do existing check here.
	// https://learn.microsoft.com/azure/sentinel/configure-fusion-rules#configure-scheduled-analytics-rules-for-fusion-detections
	params := alertrules.FusionAlertRule{
		Properties: &alertrules.FusionAlertRuleProperties{
			AlertRuleTemplateName: d.Get("alert_rule_template_guid").(string),
			Enabled:               d.Get("enabled").(bool),
			SourceSettings:        expandFusionSourceSettings(d.Get("source").([]any)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleFusionRead(d, meta)
}
