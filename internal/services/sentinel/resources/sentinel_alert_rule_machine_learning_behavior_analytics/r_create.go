// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_machine_learning_behavior_analytics

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleMLBehaviorAnalyticsCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := sentinelClient.From(meta).AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	name := d.Get("name").(string)

	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_machine_learning_behavior_analytics", id.ID())
			}
		}
	}

	params := alertrules.MLBehaviorAnalyticsAlertRule{
		Properties: &alertrules.MLBehaviorAnalyticsAlertRuleProperties{
			AlertRuleTemplateName: d.Get("alert_rule_template_guid").(string),
			Enabled:               d.Get("enabled").(bool),
		},
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
		}

		if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMLBehaviorAnalytics); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleMLBehaviorAnalyticsRead(d, meta)
}
