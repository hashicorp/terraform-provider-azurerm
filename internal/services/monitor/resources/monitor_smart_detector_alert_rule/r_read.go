// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_smart_detector_alert_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-01-01/actiongroupsapis"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorSmartDetectorAlertRuleRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := smartdetectoralertrules.ParseSmartDetectorAlertRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, smartdetectoralertrules.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SmartDetectorAlertRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("enabled", props.State == smartdetectoralertrules.AlertRuleStateEnabled)
			d.Set("frequency", props.Frequency)
			d.Set("severity", string(props.Severity))
			d.Set("scope_resource_ids", props.Scope)
			d.Set("detector_type", props.Detector.Id)

			throttlingDuration := ""
			if props.Throttling != nil && props.Throttling.Duration != nil {
				throttlingDuration = *props.Throttling.Duration
			}
			d.Set("throttling_duration", throttlingDuration)

			actionGroup, err := flattenMonitorSmartDetectorAlertRuleActionGroup(&props.ActionGroups)
			if err != nil {
				return fmt.Errorf("flatten `action_group`: %+v", err)
			}
			if err := d.Set("action_group", actionGroup); err != nil {
				return fmt.Errorf("setting `action_group`: %+v", err)
			}
		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return nil
}

func flattenMonitorSmartDetectorAlertRuleActionGroup(input *smartdetectoralertrules.ActionGroupsInformation) ([]any, error) {
	if input == nil {
		return []any{}, nil
	}

	var customEmailSubject, CustomWebhookPayload string
	if input.CustomEmailSubject != nil {
		customEmailSubject = *input.CustomEmailSubject
	}
	if input.CustomWebhookPayload != nil {
		CustomWebhookPayload = *input.CustomWebhookPayload
	}

	groupIds := make([]string, 0)
	for _, idRaw := range input.GroupIds {
		id, err := actiongroupsapis.ParseActionGroupIDInsensitively(idRaw)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %v", idRaw, err)
		}
		groupIds = append(groupIds, id.ID())
	}

	return []any{
		map[string]any{
			"ids":             groupIds,
			"email_subject":   customEmailSubject,
			"webhook_payload": CustomWebhookPayload,
		},
	}, nil
}
