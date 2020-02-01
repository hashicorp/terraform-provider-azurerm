package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRulesCreateUpdate(d *schema.ResourceData, meta interface{}, action insights.BasicAction, schedule *insights.Schedule) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor Scheduled Query Rules %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules_alert", *existing.ID)
		}
	}

	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := insights.True
	if !enabledRaw {
		enabled = insights.False
	}

	location := azure.NormalizeLocation(d.Get("location"))

	source := expandMonitorScheduledQueryRulesSource(d)

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.LogSearchRuleResource{
		Location: utils.String(location),
		LogSearchRule: &insights.LogSearchRule{
			Description: utils.String(description),
			Enabled:     enabled,
			Source:      source,
			Schedule:    schedule,
			Action:      action,
		},
		Tags: expandedTags,
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating or updating scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Scheduled query rule %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorScheduledQueryRulesRead(d, meta)
}

func resourceArmMonitorScheduledQueryRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["scheduledqueryrules"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if lastUpdated := resp.LastUpdatedTime; lastUpdated != nil {
		d.Set("last_updated_time", lastUpdated.Format(time.RFC3339))
	}
	d.Set("provisioning_state", resp.ProvisioningState)

	if resp.Enabled == insights.True {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	d.Set("description", resp.Description)

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

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorScheduledQueryRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["scheduledqueryrules"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesSource(d *schema.ResourceData) *insights.Source {
	authorizedResourceIDs := d.Get("authorized_resource_ids").(*schema.Set).List()
	dataSourceID := d.Get("data_source_id").(string)
	query, ok := d.Get("query").(string)
	source := insights.Source{
		AuthorizedResources: utils.ExpandStringSlice(authorizedResourceIDs),
		DataSourceID:        utils.String(dataSourceID),
		QueryType:           insights.ResultCount,
	}
	if ok {
		source.Query = utils.String(query)
	}

	return &source
}
