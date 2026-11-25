// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMonitorScheduledQueryRulesAlertFlatten(d *pluginsdk.ResourceData, id *scheduledqueryrules.ScheduledQueryRuleId, model *scheduledqueryrules.LogSearchRuleResource) error {
	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties
		d.Set("auto_mitigation_enabled", props.AutoMitigate)
		d.Set("description", props.Description)
		if props.Enabled != nil && *props.Enabled == scheduledqueryrules.EnabledTrue {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}

		action, ok := props.Action.(scheduledqueryrules.AlertingAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", *id, props.Action)
		}
		if err := d.Set("action", flattenAzureRmScheduledQueryRulesAlertAction(action.AznsAction)); err != nil {
			return fmt.Errorf("setting `action`: %+v", err)
		}
		severity, err := strconv.Atoi(string(action.Severity))
		if err != nil {
			return fmt.Errorf("converting action.Severity %q to int in %s: %+v", action.Severity, *id, err)
		}
		d.Set("severity", severity)
		d.Set("throttling", action.ThrottlingInMin)
		if err = d.Set("trigger", flattenAzureRmScheduledQueryRulesAlertTrigger(action.Trigger)); err != nil {
			return fmt.Errorf("setting `trigger`: %+v", err)
		}

		if schedule := props.Schedule; schedule != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}

		d.Set("authorized_resource_ids", pluginsdk.FlattenSlice(props.Source.AuthorizedResources))
		d.Set("data_source_id", props.Source.DataSourceId)
		d.Set("query", props.Source.Query)
		d.Set("query_type", pointer.FromEnum(props.Source.QueryType))

		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}
