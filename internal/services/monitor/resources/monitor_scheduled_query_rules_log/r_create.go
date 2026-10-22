// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_log

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorScheduledQueryRulesLogCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	action := expandMonitorScheduledQueryRulesLogToMetricAction(d)
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules_alert", id.ID())
			}
		}
	}

	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := scheduledqueryrules.EnabledTrue
	if !enabledRaw {
		enabled = scheduledqueryrules.EnabledFalse
	}

	location := location.Normalize(d.Get("location").(string))

	source := helpers.ExpandMonitorScheduledQueryRulesCommonSource(d)

	t := d.Get("tags").(map[string]any)

	parameters := scheduledqueryrules.LogSearchRuleResource{
		Location: location,
		Properties: scheduledqueryrules.LogSearchRule{
			Description: pointer.To(description),
			Enabled:     pointer.To(enabled),
			Source:      source,
			Action:      action,
		},
		Tags: pluginsdk.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorScheduledQueryRulesLogRead(d, meta)
}

func expandMonitorScheduledQueryRulesLogCriteria(input []any) []scheduledqueryrules.Criteria {
	criteria := make([]scheduledqueryrules.Criteria, 0)
	if len(input) == 0 {
		return criteria
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]any)
		if !ok {
			continue
		}

		dimensions := make([]scheduledqueryrules.Dimension, 0)
		for _, dimension := range v["dimension"].(*pluginsdk.Set).List() {
			if dimension == nil {
				continue
			}
			dVal, ok := dimension.(map[string]any)
			if !ok {
				continue
			}

			dimensions = append(dimensions, scheduledqueryrules.Dimension{
				Name:     dVal["name"].(string),
				Operator: scheduledqueryrules.Operator(dVal["operator"].(string)),
				Values:   helpers.ExpandStringValues(dVal["values"].([]any)),
			})
		}

		criteria = append(criteria, scheduledqueryrules.Criteria{
			MetricName: v["metric_name"].(string),
			Dimensions: &dimensions,
		})
	}
	return criteria
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *pluginsdk.ResourceData) *scheduledqueryrules.LogToMetricAction {
	criteriaRaw := d.Get("criteria").([]any)
	criteria := expandMonitorScheduledQueryRulesLogCriteria(criteriaRaw)

	action := scheduledqueryrules.LogToMetricAction{
		Criteria: criteria,
	}

	return &action
}
