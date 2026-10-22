// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_log

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMonitorScheduledQueryRulesLogRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		props := model.Properties

		d.Set("description", props.Description)
		if props.Enabled != nil && *props.Enabled == scheduledqueryrules.EnabledTrue {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}

		action, ok := props.Action.(scheduledqueryrules.LogToMetricAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", *id, props.Action)
		}

		if err = d.Set("criteria", flattenAzureRmScheduledQueryRulesLogCriteria(action.Criteria)); err != nil {
			return fmt.Errorf("setting `criteria`: %+v", err)
		}

		if props.Source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", pluginsdk.FlattenSlice(props.Source.AuthorizedResources))
		}

		d.Set("data_source_id", props.Source.DataSourceId)

		if err = d.Set("tags", pluginsdk.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}
