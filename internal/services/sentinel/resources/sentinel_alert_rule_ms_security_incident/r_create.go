// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_ms_security_incident

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAlertRuleMsSecurityIncidentCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
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
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_alert_rule_ms_security_incident", id.ID())
			}
		}
	}

	param := alertrules.MicrosoftSecurityIncidentCreationAlertRule{
		Properties: &alertrules.MicrosoftSecurityIncidentCreationAlertRuleProperties{
			ProductFilter:    alertrules.MicrosoftSecurityProductName(d.Get("product_filter").(string)),
			DisplayName:      d.Get("display_name").(string),
			Description:      pointer.To(d.Get("description").(string)),
			Enabled:          d.Get("enabled").(bool),
			SeveritiesFilter: expandAlertRuleMsSecurityIncidentSeverityFilter(d.Get("severity_filter").(*pluginsdk.Set).List()),
		},
	}

	if v, ok := d.GetOk("alert_rule_template_guid"); ok {
		param.Properties.AlertRuleTemplateName = pointer.To(v.(string))
	}

	if dnf, ok := d.GetOk("display_name_filter"); ok {
		param.Properties.DisplayNamesFilter = pluginsdk.ExpandStringSlice(dnf.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("display_name_exclude_filter"); ok {
		param.Properties.DisplayNamesExcludeFilter = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if !d.IsNewResource() {
		resp, err := client.Get(ctx, id)
		if err != nil {
			return fmt.Errorf("retrieving Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
		}

		if err := helpers.AssertAlertRuleKind(resp.Model, alertrules.AlertRuleKindMicrosoftSecurityIncidentCreation); err != nil {
			return fmt.Errorf("asserting alert rule of %q: %+v", id, err)
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, param); err != nil {
		return fmt.Errorf("creating Sentinel Alert Rule Ms Security Incident %q: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAlertRuleMsSecurityIncidentRead(d, meta)
}

func expandAlertRuleMsSecurityIncidentSeverityFilter(input []any) *[]alertrules.AlertSeverity {
	result := make([]alertrules.AlertSeverity, 0)

	for _, e := range input {
		result = append(result, alertrules.AlertSeverity(e.(string)))
	}

	return &result
}
