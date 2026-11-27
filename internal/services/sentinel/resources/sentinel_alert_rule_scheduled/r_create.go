// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_scheduled

import (
	"errors"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	sentinelClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/rickb777/period"
)

func resourceSentinelAlertRuleScheduledCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for existing Sentinel Alert Rule Scheduled %q: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_scheduled", id.ID())
			}
		}
	}

	// Sanity checks

	// query frequency must <= query period: ensure there is no gaps in the overall query coverage.
	queryFreq := d.Get("query_frequency").(string)
	queryFreqDuration := period.MustParse(queryFreq).DurationApprox()

	queryPeriod := d.Get("query_period").(string)
	queryPeriodDuration := period.MustParse(queryPeriod).DurationApprox()
	if queryFreqDuration > queryPeriodDuration {
		return fmt.Errorf("`query_frequency`(%v) should not be larger than `query period`(%v), which introduce gaps in the overall query coverage", queryFreq, queryPeriod)
	}

	// query frequency must <= suppression duration: otherwise suppression has no effect.
	suppressionDuration := d.Get("suppression_duration").(string)
	suppressionEnabled := d.Get("suppression_enabled").(bool)
	if suppressionEnabled {
		suppressionDurationDuration := period.MustParse(suppressionDuration).DurationApprox()
		if queryFreqDuration > suppressionDurationDuration {
			return fmt.Errorf("`query_frequency`(%v) should not be larger than `suppression_duration`(%v), which makes suppression pointless", queryFreq, suppressionDuration)
		}
	}

	incident := helpers.ExpandAlertRuleIncidentConfiguration(d.Get("incident").([]any), "create_incident_enabled", false)

	param := alertrules.ScheduledAlertRule{
		Properties: &alertrules.ScheduledAlertRuleProperties{
			Description:           pointer.To(d.Get("description").(string)),
			DisplayName:           d.Get("display_name").(string),
			Tactics:               helpers.ExpandAlertRuleTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			Techniques:            helpers.ExpandAlertRuleTechnicals(d.Get("techniques").(*pluginsdk.Set).List()),
			IncidentConfiguration: incident,
			Severity:              pointer.ToEnum[alertrules.AlertSeverity](d.Get("severity").(string)),
			Enabled:               d.Get("enabled").(bool),
			Query:                 pointer.To(d.Get("query").(string)),
			QueryFrequency:        pointer.To(queryFreq),
			QueryPeriod:           pointer.To(queryPeriod),
			SuppressionEnabled:    suppressionEnabled,
			SuppressionDuration:   suppressionDuration,
			TriggerOperator:       pointer.ToEnum[alertrules.TriggerOperator](d.Get("trigger_operator").(string)),
			TriggerThreshold:      pointer.To(int64(d.Get("trigger_threshold").(int))),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.Properties.AlertRuleTemplateName = pointer.To(v.(string))
	}
	if v, ok := d.GetOk("alert_rule_template_version"); ok {
		param.Properties.TemplateVersion = pointer.To(v.(string))
	}
	if v, ok := d.GetOk("event_grouping"); ok {
		param.Properties.EventGroupingSettings = helpers.ExpandAlertRuleEventGroupingSetting(v.([]any))
	}
	if v, ok := d.GetOk("alert_details_override"); ok {
		param.Properties.AlertDetailsOverride = helpers.ExpandAlertRuleAlertDetailsOverride(v.([]any))
	}
	if v, ok := d.GetOk("custom_details"); ok {
		param.Properties.CustomDetails = pluginsdk.ExpandPtrMapStringString(v.(map[string]any))
	}

	entityMappingCount := 0
	sentinelEntityMappingCount := 0
	if v, ok := d.GetOk("entity_mapping"); ok {
		param.Properties.EntityMappings = helpers.ExpandAlertRuleEntityMapping(v.([]any))
		entityMappingCount = len(*param.Properties.EntityMappings)
	}
	if v, ok := d.GetOk("sentinel_entity_mapping"); ok {
		param.Properties.SentinelEntitiesMappings = helpers.ExpandAlertRuleSentinelEntityMapping(v.([]any))
		sentinelEntityMappingCount = len(*param.Properties.SentinelEntitiesMappings)
	}

	// the max number of `sentinel_entity_mapping` and `entity_mapping` together is 10
	if entityMappingCount+sentinelEntityMappingCount > 10 {
		return errors.New("`entity_mapping` and `sentinel_entity_mapping` together can't exceed 10")
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Scheduled %q: %+v", id, err)
		}

		if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindScheduled); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Scheduled %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleScheduledRead(d, meta)
}
