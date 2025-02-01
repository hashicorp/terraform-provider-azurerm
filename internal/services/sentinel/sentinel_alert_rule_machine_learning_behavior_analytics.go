// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package sentinel

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleMLBehaviorAnalytics() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSentinelAlertRuleMLBehaviorAnalyticsCreateUpdate,
		Read:   resourceSentinelAlertRuleMLBehaviorAnalyticsRead,
		Update: resourceSentinelAlertRuleMLBehaviorAnalyticsCreateUpdate,
		Delete: resourceSentinelAlertRuleMLBehaviorAnalyticsDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := alertrules.ParseAlertRuleID(id)
			return err
		}, importSentinelAlertRule(alertrules.AlertRuleKindMLBehaviorAnalytics)),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"log_analytics_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: alertrules.ValidateWorkspaceID,
			},

			"alert_rule_template_guid": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsUUID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceSentinelAlertRuleMLBehaviorAnalyticsCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	workspaceID, err := alertrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := alertrules.NewAlertRuleID(workspaceID.SubscriptionId, workspaceID.ResourceGroupName, workspaceID.WorkspaceName, name)

	if d.IsNewResource() {
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

		if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMLBehaviorAnalytics); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleMLBehaviorAnalyticsRead(d, meta)
}

func resourceSentinelAlertRuleMLBehaviorAnalyticsRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Sentinel Alert Rule MLBehaviorAnalytics %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
	}

	if err := assertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMLBehaviorAnalytics); err != nil {
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

func resourceSentinelAlertRuleMLBehaviorAnalyticsDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Sentinel.AlertRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertrules.ParseAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting Sentinel Alert Rule MLBehaviorAnalytics %q: %+v", id, err)
	}

	return nil
}
