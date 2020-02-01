package monitor

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmMonitorScheduledQueryRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Scheduled Query Rule %q was not found", name)
		}
		return fmt.Errorf("Error reading Scheduled Query Rule: %+v", err)
	}

	d.SetId(*resp.ID)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("description", resp.Description)
	if resp.Enabled == insights.True {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	switch action := resp.Action.(type) {
	case insights.AlertingAction:
		if err = d.Set("action", flattenAzureRmScheduledQueryRulesAlertAction(action.AznsAction)); err != nil {
			return fmt.Errorf("Error setting `action`: %+v", err)
		}
		severity, err := strconv.Atoi(string(action.Severity))
		if err != nil {
			return fmt.Errorf("Error converting action.Severity %q in query rule %q to int (resource group %q): %+v", action.Severity, name, resourceGroup, err)
		}
		d.Set("severity", severity)
		d.Set("throttling", action.ThrottlingInMin)
		if err = d.Set("trigger", flattenAzureRmScheduledQueryRulesAlertTrigger(action.Trigger)); err != nil {
			return fmt.Errorf("Error setting `trigger`: %+v", err)
		}
	case insights.LogToMetricAction:
		if err = d.Set("criteria", flattenAzureRmScheduledQueryRulesLogCriteria(action.Criteria)); err != nil {
			return fmt.Errorf("Error setting `criteria`: %+v", err)
		}
	default:
		return fmt.Errorf("Unknown action type in scheduled query rule %q (resource group %q): %T", name, resourceGroup, resp.Action)
	}

	if schedule := resp.Schedule; schedule != nil {
		if schedule.FrequencyInMinutes != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
		}
		if schedule.TimeWindowInMinutes != nil {
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}
	}

	if source := resp.Source; source != nil {
		if source.AuthorizedResources != nil {
			d.Set("authorized_resource_ids", utils.FlattenStringSlice(source.AuthorizedResources))
		}
		if source.DataSourceID != nil {
			d.Set("data_source_id", source.DataSourceID)
		}
		if source.Query != nil {
			d.Set("query", source.Query)
		}
		d.Set("query_type", string(source.QueryType))
	}

	if lastUpdated := resp.LastUpdatedTime; lastUpdated != nil {
		d.Set("last_updated_time", lastUpdated.Format(time.RFC3339))
	}
	d.Set("provisioning_state", resp.ProvisioningState)

	return tags.FlattenAndSet(d, resp.Tags)
}
