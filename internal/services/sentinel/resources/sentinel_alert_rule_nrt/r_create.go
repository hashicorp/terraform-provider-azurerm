// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_nrt

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
)

func resourceSentinelAlertRuleNrtCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
					return fmt.Errorf("checking for existing %q: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_nrt", id.ID())
			}
		}
	}

	param := alertrules.NrtAlertRule{
		Properties: &alertrules.NrtAlertRuleProperties{
			Description:           pointer.To(d.Get("description").(string)),
			DisplayName:           d.Get("display_name").(string),
			Techniques:            helpers.ExpandAlertRuleTechnicals(d.Get("techniques").(*pluginsdk.Set).List()),
			Tactics:               helpers.ExpandAlertRuleTactics(d.Get("tactics").(*pluginsdk.Set).List()),
			IncidentConfiguration: helpers.ExpandAlertRuleIncidentConfiguration(d.Get("incident").([]any), "create_incident_enabled", false),
			Severity:              alertrules.AlertSeverity(d.Get("severity").(string)),
			Enabled:               d.Get("enabled").(bool),
			Query:                 d.Get("query").(string),
			SuppressionEnabled:    d.Get("suppression_enabled").(bool),
			SuppressionDuration:   d.Get("suppression_duration").(string),
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

	// the max number of `sentinel_entity_mapping` and `entity_mapping` together is 5
	if entityMappingCount+sentinelEntityMappingCount > 5 {
		return errors.New("`entity_mapping` and `sentinel_entity_mapping` together can't exceed 5")
	}

	// Service avoid concurrent update of this resource via checking the "etag" to guarantee it is the same value as last Read.
	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving %q: %+v", id, err)
		}

		if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindNRT); err != nil {
			return fmt.Errorf("asserting %q: %+v", id, err)
		}
		if model := resp.Model; model != nil {
			if rule, ok := model.(alertrules.NrtAlertRule); ok {
				param.Etag = rule.Etag
			}
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleNrtRead(d, meta)
}
