// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"fmt"
	"strings"

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

func resourceMonitorScheduledQueryRulesAlertCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	action := expandMonitorScheduledQueryRulesAlertingAction(d)
	schedule := expandMonitorScheduledQueryRulesAlertSchedule(d)
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)
	if timeWindow < frequency {
		return fmt.Errorf("in parameter values for %s: time_window must be greater than or equal to frequency", id)
	}

	query := d.Get("query").(string)
	_, ok := d.GetOk("metric_trigger")
	if ok {
		if !strings.Contains(query, "summarize") || !strings.Contains(query, "AggregatedValue") || !strings.Contains(query, "bin") {
			return fmt.Errorf("in parameter values for %s: query must contain summarize, AggregatedValue, and bin when metric_trigger is specified", id)
		}
	}

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(monitorScheduledQueryRulesAlertResourceName, id.ID())
			}
		}
	}

	autoMitigate := d.Get("auto_mitigation_enabled").(bool)
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
			Description:  pointer.To(description),
			Enabled:      pointer.To(enabled),
			Source:       source,
			Schedule:     schedule,
			Action:       action,
			AutoMitigate: pointer.To(autoMitigate),
		},
		Tags: pluginsdk.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceMonitorScheduledQueryRulesAlertRead(d, meta)
}
