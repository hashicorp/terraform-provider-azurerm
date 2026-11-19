// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDatadogTagRulesRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Datadog.TagRules
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := tagrules.ParseTagRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
	}

	monitorId := tagrules.NewMonitorID(id.SubscriptionId, id.ResourceGroupName, id.MonitorName)
	d.Set("datadog_monitor_id", monitorId.ID())
	d.Set("name", id.TagRuleName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if err := d.Set("log", flattenLogRules(props.LogRules)); err != nil {
				return fmt.Errorf("setting `log`: %+v", err)
			}
			if err := d.Set("metric", flattenMetricRules(props.MetricRules)); err != nil {
				return fmt.Errorf("setting `metric`: %+v", err)
			}
		}
	}

	return nil
}

func isDefaultSettings(input *tagrules.MonitoringTagRules) bool {
	if input == nil {
		return false
	}

	if input.Properties == nil || input.Properties.LogRules == nil || input.Properties.MetricRules == nil {
		return false
	}

	logRules := input.Properties.LogRules
	metricRules := input.Properties.MetricRules
	return (logRules.SendAadLogs != nil && !*logRules.SendAadLogs) &&
		(logRules.SendSubscriptionLogs != nil && !*logRules.SendSubscriptionLogs) &&
		(logRules.SendResourceLogs != nil && !*logRules.SendResourceLogs) &&
		(logRules.FilteringTags != nil && len(*logRules.FilteringTags) == 0) &&
		(metricRules.FilteringTags != nil && len(*metricRules.FilteringTags) == 0)
}
