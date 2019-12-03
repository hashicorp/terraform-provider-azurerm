package azurerm

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"action_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Alerting",
					"LogToMetric",
				}, false),
			},
			"authorized_resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"azns_action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_group": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"custom_webhook_payload": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "{}",
							ValidateFunc: validation.ValidateJsonString,
						},
						"email_subject": {
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
						"dimension": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Include",
										}, false),
									},
									"values": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"data_source_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
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
			"frequency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"last_updated_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ResultCount",
				ValidateFunc: validation.StringInSlice([]string{
					"ResultCount",
				}, false),
			},
			"severity": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					0,
					1,
					2,
					3,
					4,
				}),
			},
			"throttling": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"time_window": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"trigger": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_trigger": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_column": {
										Type:     schema.TypeString,
										Required: true,
									},
									"metric_trigger_type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Consecutive",
											"Total",
										}, false),
									},
									"operator": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"GreaterThan",
											"LessThan",
											"Equal",
										}, false),
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Required: true,
									},
								},
							},
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GreaterThan",
								"LessThan",
								"Equal",
							}, false),
						},
						"threshold": {
							Type:     schema.TypeFloat,
							Required: true,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*ArmClient).StopContext, d)
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
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules", *existing.ID)
		}
	}

	actionType := d.Get("action_type").(string)
	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := insights.True
	if enabledRaw == false {
		enabled = insights.False
	}

	location := azure.NormalizeLocation(d.Get("location"))

	var action insights.BasicAction
	switch actionType {
	case "Alerting":
		action = expandMonitorScheduledQueryRulesAlertingAction(d)
	case "LogToMetric":
		action = expandMonitorScheduledQueryRulesLogToMetricAction(d)
	default:
		return fmt.Errorf("Invalid action_type %q. Value must be either 'Alerting' or 'LogToMetric'", actionType)
	}

	source := expandMonitorScheduledQueryRulesSource(d)
	schedule := expandMonitorScheduledQueryRulesSchedule(d)

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
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
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
		d.Set("last_updated_time", *lastUpdated)
	}
	d.Set("provisioning_state", resp.ProvisioningState)

	if resp.Enabled == insights.True {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}

	d.Set("description", *resp.Description)

	switch action := resp.Action.(type) {
	case insights.AlertingAction:
		d.Set("action_type", "Alerting")
		d.Set("azns_action", flattenAzureRmScheduledQueryRulesAznsAction(action.AznsAction))
		severity, err := strconv.Atoi(string(action.Severity))
		if err != nil {
			return fmt.Errorf("Error converting action.Severity %q in query rule %q to int (resource group %q): %+v", action.Severity, name, resourceGroup, err)
		}
		d.Set("severity", severity)
		d.Set("throttling", *action.ThrottlingInMin)
		d.Set("trigger", flattenAzureRmScheduledQueryRulesTrigger(action.Trigger))
	case insights.LogToMetricAction:
		d.Set("action_type", "LogToMetric")
		d.Set("criteria", flattenAzureRmScheduledQueryRulesCriteria(action.Criteria))
	default:
		return fmt.Errorf("Unknown action type in scheduled query rule %q (resource group %q): %T", name, resourceGroup, resp.Action)
	}

	if schedule := resp.Schedule; schedule != nil {
		if schedule.FrequencyInMinutes != nil {
			d.Set("frequency", *schedule.FrequencyInMinutes)
		}
		if schedule.TimeWindowInMinutes != nil {
			d.Set("time_window", *schedule.TimeWindowInMinutes)
		}
	}

	if source := resp.Source; source != nil {
		if source.AuthorizedResources != nil {
			d.Set("authorized_resources", *source.AuthorizedResources)
		}
		if source.DataSourceID != nil {
			d.Set("data_source_id", *source.DataSourceID)
		}
		if source.Query != nil {
			d.Set("query", *source.Query)
		}
		d.Set("query_type", string(source.QueryType))
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorScheduledQueryRulesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*ArmClient).StopContext, d)
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

func expandMonitorScheduledQueryRulesAlertingAction(d *schema.ResourceData) *insights.AlertingAction {
	aznsActionRaw := d.Get("azns_action").(*schema.Set).List()
	aznsAction := expandMonitorScheduledQueryRulesAznsAction(aznsActionRaw)
	severityRaw := d.Get("severity").(int)
	severity := strconv.Itoa(severityRaw)
	throttling := d.Get("throttling").(int)

	triggerRaw := d.Get("trigger").(*schema.Set).List()
	trigger := expandMonitorScheduledQueryRulesTrigger(triggerRaw)

	action := insights.AlertingAction{
		AznsAction:      aznsAction,
		Severity:        insights.AlertSeverity(severity),
		ThrottlingInMin: utils.Int32(int32(throttling)),
		Trigger:         trigger,
		OdataType:       insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesAlertingAction,
	}

	return &action
}

func expandMonitorScheduledQueryRulesAznsAction(input []interface{}) *insights.AzNsActionGroup {
	result := insights.AzNsActionGroup{}

	for _, item := range input {
		v := item.(map[string]interface{})
		actionGroups := v["action_group"].(*schema.Set).List()

		result.ActionGroup = utils.ExpandStringSlice(actionGroups)
		result.EmailSubject = utils.String(v["email_subject"].(string))
		result.CustomWebhookPayload = utils.String(v["custom_webhook_payload"].(string))
	}

	return &result
}

func expandMonitorScheduledQueryRulesCriteria(input []interface{}) *[]insights.Criteria {
	criteria := make([]insights.Criteria, 0)
	for _, item := range input {
		v := item.(map[string]interface{})

		dimensions := make([]insights.Dimension, 0)
		for _, dimension := range v["dimension"].(*schema.Set).List() {
			dVal := dimension.(map[string]interface{})
			dimensions = append(dimensions, insights.Dimension{
				Name:     utils.String(dVal["name"].(string)),
				Operator: utils.String(dVal["operator"].(string)),
				Values:   utils.ExpandStringSlice(dVal["values"].([]interface{})),
			})
		}

		criteria = append(criteria, insights.Criteria{
			MetricName: utils.String(v["metric_name"].(string)),
			Dimensions: &dimensions,
		})
	}
	return &criteria
}

func expandMonitorScheduledQueryRulesLogToMetricAction(d *schema.ResourceData) *insights.LogToMetricAction {
	criteriaRaw := d.Get("criteria").(*schema.Set).List()
	criteria := expandMonitorScheduledQueryRulesCriteria(criteriaRaw)

	action := insights.LogToMetricAction{
		Criteria:  criteria,
		OdataType: insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesLogToMetricAction,
	}

	return &action
}

func expandMonitorScheduledQueryRulesSchedule(d *schema.ResourceData) *insights.Schedule {
	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)

	schedule := insights.Schedule{
		FrequencyInMinutes:  utils.Int32(int32(frequency)),
		TimeWindowInMinutes: utils.Int32(int32(timeWindow)),
	}

	return &schedule
}

func expandMonitorScheduledQueryRulesMetricTrigger(input []interface{}) *insights.LogMetricTrigger {
	result := insights.LogMetricTrigger{}

	for _, item := range input {
		v := item.(map[string]interface{})
		result.ThresholdOperator = insights.ConditionalOperator(v["operator"].(string))
		result.Threshold = utils.Float(v["threshold"].(float64))
		result.MetricTriggerType = insights.MetricTriggerType(v["metric_trigger_type"].(string))
		result.MetricColumn = utils.String(v["metric_column"].(string))
	}

	return &result
}

func expandMonitorScheduledQueryRulesSource(d *schema.ResourceData) *insights.Source {
	authorizedResources := d.Get("authorized_resources").(*schema.Set).List()
	dataSourceID := d.Get("data_source_id").(string)
	query := d.Get("query").(string)

	source := insights.Source{
		AuthorizedResources: utils.ExpandStringSlice(authorizedResources),
		DataSourceID:        utils.String(dataSourceID),
		Query:               utils.String(query),
		QueryType:           insights.ResultCount,
	}

	return &source
}

func expandMonitorScheduledQueryRulesTrigger(input []interface{}) *insights.TriggerCondition {
	result := insights.TriggerCondition{}

	for _, item := range input {
		v := item.(map[string]interface{})
		metricTriggerRaw := v["metric_trigger"].(*schema.Set).List()

		result.ThresholdOperator = insights.ConditionalOperator(v["operator"].(string))
		result.Threshold = utils.Float(v["threshold"].(float64))
		result.MetricTrigger = expandMonitorScheduledQueryRulesMetricTrigger(metricTriggerRaw)
	}

	return &result
}

func flattenAzureRmScheduledQueryRulesAznsAction(input *insights.AzNsActionGroup) []interface{} {
	result := make([]interface{}, 0)
	v := make(map[string]interface{})

	if input != nil {
		if input.ActionGroup != nil {
			v["action_group"] = *input.ActionGroup
		}
		v["email_subject"] = input.EmailSubject
		v["custom_webhook_payload"] = input.CustomWebhookPayload
	}
	result = append(result, v)

	return result
}

func flattenAzureRmScheduledQueryRulesCriteria(input *[]insights.Criteria) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, criteria := range *input {
			v := make(map[string]interface{})

			/*if err = d.Set("azure_function_receiver", flattenMonitorActionGroupAzureFunctionReceiver(group.AzureFunctionReceivers)); err != nil {
				return fmt.Errorf("Error setting `azure_function_receiver`: %+v", err)
			}*/
			v["dimension"] = flattenAzureRmScheduledQueryRulesDimension(criteria.Dimensions)
			v["metric_name"] = *criteria.MetricName

			result = append(result, v)
		}
	}

	return result
}

func flattenAzureRmScheduledQueryRulesDimension(input *[]insights.Dimension) []interface{} {
	result := make([]interface{}, 0)

	if input != nil {
		for _, dimension := range *input {
			v := make(map[string]interface{})

			if dimension.Name != nil {
				v["name"] = *dimension.Name
			}

			if dimension.Operator != nil {
				v["operator"] = *dimension.Operator
			}

			if dimension.Values != nil {
				v["values"] = *dimension.Values
			}

			result = append(result, v)
		}
	}

	return result
}

func flattenAzureRmScheduledQueryRulesMetricTrigger(input *insights.LogMetricTrigger) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	result["operator"] = string(input.ThresholdOperator)

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	result["metric_trigger_type"] = string(input.MetricTriggerType)

	if input.MetricColumn != nil {
		result["metric_column"] = *input.MetricColumn
	}
	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesSchedule(input *insights.Schedule) []interface{} {
	result := make(map[string]interface{})

	if input == nil {
		return []interface{}{}
	}

	if input.FrequencyInMinutes != nil {
		result["frequency_in_minutes"] = int(*input.FrequencyInMinutes)
	}

	if input.TimeWindowInMinutes != nil {
		result["time_window_in_minutes"] = int(*input.TimeWindowInMinutes)
	}

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesSource(input *insights.Source) []interface{} {
	result := make(map[string]interface{})

	if input.AuthorizedResources != nil {
		result["authorized_resources"] = *input.AuthorizedResources
	}
	if input.DataSourceID != nil {
		result["data_source_id"] = *input.DataSourceID
	}
	if input.Query != nil {
		result["query"] = *input.Query
	}
	result["query_type"] = string(input.QueryType)

	return []interface{}{result}
}

func flattenAzureRmScheduledQueryRulesTrigger(input *insights.TriggerCondition) []interface{} {
	result := make(map[string]interface{})

	result["operator"] = string(input.ThresholdOperator)

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	if input.MetricTrigger != nil {
		result["metric_trigger"] = flattenAzureRmScheduledQueryRulesMetricTrigger(input.MetricTrigger)
	}

	return []interface{}{result}
}
