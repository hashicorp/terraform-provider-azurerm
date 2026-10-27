// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_smart_detector_alert_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-06-01/smartdetectoralertrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorSmartDetectorAlertRuleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.SmartDetectorAlertRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := smartdetectoralertrules.NewSmartDetectorAlertRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id, smartdetectoralertrules.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_monitor_smart_detector_alert_rule", id.ID())
			}
		}
	}

	state := smartdetectoralertrules.AlertRuleStateDisabled
	if d.Get("enabled").(bool) {
		state = smartdetectoralertrules.AlertRuleStateEnabled
	}

	actionRule := smartdetectoralertrules.AlertRule{
		// the location is always global from the portal
		Location: pointer.To(location.Normalize("Global")),
		Properties: &smartdetectoralertrules.AlertRuleProperties{
			Description: pointer.To(d.Get("description").(string)),
			State:       state,
			Severity:    smartdetectoralertrules.Severity(d.Get("severity").(string)),
			Frequency:   d.Get("frequency").(string),
			Detector: smartdetectoralertrules.Detector{
				Id: d.Get("detector_type").(string),
			},
			Scope:        pointer.From(pluginsdk.ExpandStringSlice(d.Get("scope_resource_ids").(*pluginsdk.Set).List())),
			ActionGroups: pointer.From(expandMonitorSmartDetectorAlertRuleActionGroup(d.Get("action_group").([]any))),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("throttling_duration"); ok {
		actionRule.Properties.Throttling = &smartdetectoralertrules.ThrottlingInformation{
			Duration: pointer.To(v.(string)),
		}
	}

	if _, err := client.CreateOrUpdate(ctx, id, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorSmartDetectorAlertRuleRead(d, meta)
}

func expandMonitorSmartDetectorAlertRuleActionGroup(input []any) *smartdetectoralertrules.ActionGroupsInformation {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)
	return &smartdetectoralertrules.ActionGroupsInformation{
		CustomEmailSubject:   pointer.To(v["email_subject"].(string)),
		CustomWebhookPayload: pointer.To(v["webhook_payload"].(string)),
		GroupIds:             pointer.From(pluginsdk.ExpandStringSlice(v["ids"].(*pluginsdk.Set).List())),
	}
}
