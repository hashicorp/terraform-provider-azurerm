package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRules() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesRead,
		Update: resourceArmMonitorScheduledQueryRulesCreateUpdate,
		Delete: resourceArmMonitorScheduledQueryRulesDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_source_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"authorized_resources": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
			"queryType": {
				Type:     schema.TypeString,
				Required: true,
				Default:  "ResultCount",
				ValidateFunc: validation.StringInSlice([]string{
					"ResultCount",
				}, false),
			},

			"frequency_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"time_window_in_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"odata.type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.AlertingAction",
					"Microsoft.WindowsAzure.Management.Monitoring.Alerts.Models.Microsoft.AppInsights.Nexus.DataContracts.Resources.ScheduledQueryRules.LogToMetricAction",
				}, false),
			},
			"severity": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"0",
					"1",
					"2",
					"3",
					"4",
				}, false),
			},
			"throttling": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"azns_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"actionGroup": {
							Type:         schema.TypeSet,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"customWebhookPayload": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.URLIsHTTPOrHTTPS,
						},
						"emailSubject": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"criteria": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimensions": {
							Type:         schema.TypeSet,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"metricName": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"thresholdOperator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GreaterThan",
								"LessThan",
								"Equal",
							}, false),
						},
						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"lastUpdatedTime": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"provisioningState": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

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
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules", *existing.ID)
		}
	}

	enabled := d.Get("enabled").(bool)
	description := d.Get("description").(string)
	displayName := d.Get("displayName").(string)
	sourceRaw := d.Get("source").(*schema.Set).List()
	scheduleRaw := d.Get("schedule").(*schema.Set).List()
	actionRaw := d.Get("action").(*schema.Set).List()

	t := d.Get("tags").(map[string]interface{})
	expandedTags := tags.Expand(t)

	parameters := insights.LogSearchRuleResource{
		Location: &azure.NormalizeLocation(d.Get("location").(string)),
		LogSearchRule: &insights.LogSearchRule{
			Enabled:     utils.Bool(enabled),
			Description: utils.String(description),
			DisplayName: utils.String(displayName),
			Source:      expandMonitorScheduledQueryRulesSource(sourceRaw),
			Schedule:    expandMonitorScheduledQueryRulesSchedule(scheduleRaw),
			Action:      expandMonitorScheduledQueryRulesAction(actionRaw),
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
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

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
	if rule := resp.ScheduledQueryRules; rule != nil {
		d.Set("enabled", rule.Enabled)
		d.Set("description", rule.Description)
		d.Set("displayName", rule.DisplayName)
		if err := d.Set("source", flattenMonitorScheduledQueryRulesSource(rule.Source)); err != nil {
			return fmt.Errorf("Error setting `source`: %+v", err)
		}
		if err := d.Set("schedule", flattenMonitorScheduledQueryRulesSchedule(rule.Schedule)); err != nil {
			return fmt.Errorf("Error setting `schedule`: %+v", err)
		}
		if err := d.Set("action", flattenMonitorScheduledQueryRulesAction(rule.Action)); err != nil {
			return fmt.Errorf("Error setting `action`: %+v", err)
		}
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorScheduledQueryRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).monitor.ScheduledQueryRulesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["ScheduledQueryRules"]

	if resp, err := client.Delete(ctx, resourceGroup, name); err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting scheduled query rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesSource(input []interface{}) (*insights.Source, error) {
	conditions := make([]insights.ScheduledQueryRulesLeafCondition, 0)
	v := input[0].(map[string]interface{})

	if category := v["category"].(string); category != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("category"),
			Equals: utils.String(category),
		})
	}
	if op := v["operation_name"].(string); op != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("operationName"),
			Equals: utils.String(op),
		})
	}
	if caller := v["caller"].(string); caller != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("caller"),
			Equals: utils.String(caller),
		})
	}
	if level := v["level"].(string); level != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("level"),
			Equals: utils.String(level),
		})
	}
	if resourceProvider := v["resource_provider"].(string); resourceProvider != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceProvider"),
			Equals: utils.String(resourceProvider),
		})
	}
	if resourceType := v["resource_type"].(string); resourceType != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceType"),
			Equals: utils.String(resourceType),
		})
	}
	if resourceGroup := v["resource_group"].(string); resourceGroup != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceGroup"),
			Equals: utils.String(resourceGroup),
		})
	}
	if id := v["resource_id"].(string); id != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("resourceId"),
			Equals: utils.String(id),
		})
	}
	if status := v["status"].(string); status != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("status"),
			Equals: utils.String(status),
		})
	}
	if subStatus := v["sub_status"].(string); subStatus != "" {
		conditions = append(conditions, insights.ScheduledQueryRulesLeafCondition{
			Field:  utils.String("subStatus"),
			Equals: utils.String(subStatus),
		})
	}

	return &insights.ScheduledQueryRulesAllOfCondition{
		AllOf: &conditions,
	}
}

func expandMonitorScheduledQueryRulesSchedule(input []interface{}) (*insights.Schedule, error) {
	actions := make([]insights.ScheduledQueryRulesActionGroup, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		if agID := v["action_group_id"].(string); agID != "" {
			props := make(map[string]*string)
			if pVal, ok := v["webhook_properties"]; ok {
				for pk, pv := range pVal.(map[string]interface{}) {
					props[pk] = utils.String(pv.(string))
				}
			}

			actions = append(actions, insights.ScheduledQueryRulesActionGroup{
				ActionGroupID:     utils.String(agID),
				WebhookProperties: props,
			})
		}
	}
	return &insights.ScheduledQueryRulesActionList{
		ActionGroups: &actions,
	}
}

func flattenAzureRmScheduledQueryRulesAction(input insights.BasicAction) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{result}
	}

	alertingAction, ok := input.(*insights.AlertingAction)
	if ok {
		result["action"] = flattenAzureRmScheduledQueryRulesAlertingAction(alertingAction)
	}

	logToMetricAction, ok := input.(*insights.LogToMetricAction)
	if ok {
		result["action"] = flattenAzureRmScheduledQueryRulesLogToMetricAction(logToMetricAction)
	}

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesAlertingAction(action *insights.AlertingAction) map[string]interface{} {
	result := make(map[string]interface{})

	result["severity"] = action.Severity

	// https://github.com/Azure/azure-sdk-for-go/blob/7a9d2769e4a581b0b1bc609c71b59af043e05c98/services/preview/monitor/mgmt/2019-06-01/insights/models.go#L1771-L1779
	if action.AznsAction != nil {
		result["azns_action"] = *action.AznsAction
	}

	if action.ThrottlingInMin != nil {
		result["throttling"] = *action.ThrottlingInMin
	}

	// https://github.com/Azure/azure-sdk-for-go/blob/7a9d2769e4a581b0b1bc609c71b59af043e05c98/services/preview/monitor/mgmt/2019-06-01/insights/models.go#L5608-L5616
	if action.Trigger != nil {
		result["trigger"] = *action.Trigger
	}

	return result
}

func flattenAzureRmScheduledQueryRulesLogToMetricAction(action *insights.LogToMetricAction) map[string]interface{} {
	result := make(map[string]interface{})

	// https://github.com/Azure/azure-sdk-for-go/blob/7a9d2769e4a581b0b1bc609c71b59af043e05c98/services/preview/monitor/mgmt/2019-06-01/insights/models.go#L1929-L1935
	if action.Criteria != nil {
		result["criteria"] = *action.Criteria
	}

	return result
}

func flattenAzureRmScheduledQueryRulesSchedule(input *insights.Schedule) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	if input.FrequencyInMinutes != nil {
		result["frequency_in_minutes"] = *input.FrequencyInMinutes
	}

	if input.TimeWindowInMinutes != nil {
		result["time_window_in_minutes"] = *input.TimeWindowInMinutes
	}

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesSource(input *insights.Source) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	if input.Query != nil {
		result["query"] = *input.Query
	}

	if input.DataSourceID != nil {
		result["data_source_id"] = *input.DataSourceID
	}

	if input.QueryType != "" {
		result["query_type"] = input.QueryType
	}

	if input.AuthorizedResources != nil {
		v := make(map[string][]string)
		resources := []string{}
		for _, authorized := range *input.AuthorizedResources {
			if authorized != "" {
				resources = append(resources, authorized)
			}
		}
		result["authorized_resources"] = resources
	}

	return []interface{}{result}
}

/*
// LogSearchRule log Search Rule Definition
type LogSearchRule struct {
	// Description - The description of the Log Search rule.
	Description *string `json:"description,omitempty"`
	// Enabled - The flag which indicates whether the Log Search rule is enabled. Value should be true or false. Possible values include: 'True', 'False'
	Enabled Enabled `json:"enabled,omitempty"`
	// Source - Data Source against which rule will Query Data
	Source *Source `json:"source,omitempty"`
	// Schedule - Schedule (Frequency, Time Window) for rule. Required for action type - AlertingAction
	Schedule *Schedule `json:"schedule,omitempty"`
	// Action - Action needs to be taken on rule execution.
	Action BasicAction `json:"action,omitempty"`
}
*/
