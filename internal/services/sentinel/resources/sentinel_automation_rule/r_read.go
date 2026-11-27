// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAutomationRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationrules.ParseAutomationRuleID(d.Id())
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

	d.Set("name", id.AutomationRuleId)
	d.Set("log_analytics_workspace_id", automationrules.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID())
	if model := resp.Model; model != nil {
		prop := model.Properties
		d.Set("display_name", prop.DisplayName)
		d.Set("order", prop.Order)

		tl := prop.TriggeringLogic
		d.Set("enabled", tl.IsEnabled)
		d.Set("triggers_on", string(tl.TriggersOn))
		d.Set("triggers_when", string(tl.TriggersWhen))
		d.Set("expiration", tl.ExpirationTimeUtc)

		conditionJSON, err := flattenAutomationRuleConditionsToJSON(tl.Conditions)
		if err != nil {
			return fmt.Errorf("flattening `condition_json`: %v", err)
		}
		d.Set("condition_json", conditionJSON)

		actionIncident, actionPlaybook, actionIncidentTask := flattenAutomationRuleActions(prop.Actions)

		if err := d.Set("action_incident", actionIncident); err != nil {
			return fmt.Errorf("setting `action_incident`: %v", err)
		}
		if err := d.Set("action_playbook", actionPlaybook); err != nil {
			return fmt.Errorf("setting `action_playbook`: %v", err)
		}
		if err := d.Set("action_incident_task", actionIncidentTask); err != nil {
			return fmt.Errorf("setting `action_incident_task`: %v", err)
		}
	}

	return nil
}
