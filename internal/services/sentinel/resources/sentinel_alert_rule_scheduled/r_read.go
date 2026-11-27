// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_scheduled

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleScheduledRead(d *pluginsdk.ResourceData, meta any) error {
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
			log.Printf("[DEBUG] Sentinel Alert Rule Scheduled %q was not found - removing from state!", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindScheduled); err != nil {
		return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		if rule, ok := model.(alertrules.ScheduledAlertRule); ok {
			d.Set("name", id.RuleId)

			workspaceId := workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
			d.Set("log_analytics_workspace_id", workspaceId.ID())

			if prop := rule.Properties; prop != nil {
				d.Set("description", prop.Description)
				d.Set("display_name", prop.DisplayName)
				if err := d.Set("tactics", helpers.FlattenAlertRuleTactics(prop.Tactics)); err != nil {
					return fmt.Errorf("setting `tactics`: %+v", err)
				}
				if err := d.Set("techniques", prop.Techniques); err != nil {
					return fmt.Errorf("setting `techniques`: %+v", err)
				}

				if err := d.Set("incident", helpers.FlattenAlertRuleIncidentConfiguration(prop.IncidentConfiguration, "create_incident_enabled", false)); err != nil {
					return fmt.Errorf("setting `incident`: %+v", err)
				}

				d.Set("severity", pointer.FromEnum(prop.Severity))
				d.Set("enabled", prop.Enabled)
				d.Set("query", prop.Query)
				d.Set("query_frequency", prop.QueryFrequency)
				d.Set("query_period", prop.QueryPeriod)
				d.Set("trigger_operator", pointer.FromEnum(prop.TriggerOperator))
				d.Set("trigger_threshold", int(pointer.From(prop.TriggerThreshold)))
				d.Set("suppression_enabled", prop.SuppressionEnabled)
				d.Set("suppression_duration", prop.SuppressionDuration)
				d.Set("alert_rule_template_guid", prop.AlertRuleTemplateName)
				d.Set("alert_rule_template_version", prop.TemplateVersion)

				if err := d.Set("event_grouping", helpers.FlattenAlertRuleEventGroupingSetting(prop.EventGroupingSettings)); err != nil {
					return fmt.Errorf("setting `event_grouping`: %+v", err)
				}
				if err := d.Set("alert_details_override", helpers.FlattenAlertRuleAlertDetailsOverride(prop.AlertDetailsOverride)); err != nil {
					return fmt.Errorf("setting `alert_details_override`: %+v", err)
				}
				if err := d.Set("custom_details", pluginsdk.FlattenPtrMapStringString(prop.CustomDetails)); err != nil {
					return fmt.Errorf("setting `custom_details`: %+v", err)
				}
				if err := d.Set("entity_mapping", helpers.FlattenAlertRuleEntityMapping(prop.EntityMappings)); err != nil {
					return fmt.Errorf("setting `entity_mapping`: %+v", err)
				}
				if err := d.Set("sentinel_entity_mapping", helpers.FlattenAlertRuleSentinelEntityMapping(prop.SentinelEntitiesMappings)); err != nil {
					return fmt.Errorf("setting `sentinel_entity_mapping`: %+v", err)
				}
			}
		}
	}

	return nil
}
