// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_automation_rule

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2024-09-01/automationrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSentinelAutomationRuleCreateOrUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Sentinel.AutomationRulesClient
	tenantId := meta.(*clients.Client).Account.TenantId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	workspaceId, err := automationrules.ParseWorkspaceID(d.Get("log_analytics_workspace_id").(string))
	if err != nil {
		return err
	}
	id := automationrules.NewAutomationRuleID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, name)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_sentinel_automation_rule", id.ID())
			}
		}
	}

	actions, err := expandAutomationRuleActions(d, tenantId)
	if err != nil {
		return err
	}
	params := automationrules.AutomationRule{
		Properties: automationrules.AutomationRuleProperties{
			DisplayName: d.Get("display_name").(string),
			Order:       int64(d.Get("order").(int)),
			TriggeringLogic: automationrules.AutomationRuleTriggeringLogic{
				IsEnabled:    d.Get("enabled").(bool),
				TriggersOn:   automationrules.TriggersOn(d.Get("triggers_on").(string)),
				TriggersWhen: automationrules.TriggersWhen(d.Get("triggers_when").(string)),
			},
			Actions: actions,
		},
	}

	conditions, err := expandAutomationRuleConditionsFromJSON(d.Get("condition_json").(string))
	if err != nil {
		return fmt.Errorf("expanding `condition_json`: %v", err)
	}
	params.Properties.TriggeringLogic.Conditions = conditions

	if expiration := d.Get("expiration").(string); expiration != "" {
		t, _ := time.Parse(time.RFC3339, expiration)
		params.Properties.TriggeringLogic.SetExpirationTimeUtcAsTime(t)
	}

	if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSentinelAutomationRuleRead(d, meta)
}
