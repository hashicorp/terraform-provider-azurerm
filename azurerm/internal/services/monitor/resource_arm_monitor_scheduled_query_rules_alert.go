package monitor

import (
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRulesAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesAlertCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesRead,
		Update: resourceArmMonitorScheduledQueryRulesAlertCreateUpdate,
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

			"authorized_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
			"action": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_group": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
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
										Type:         schema.TypeFloat,
										Required:     true,
										ValidateFunc: validation.NoZeroValues,
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

func resourceArmMonitorScheduledQueryRulesAlertCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	var action insights.BasicAction
	action = expandMonitorScheduledQueryRulesAlertingAction(d)
	schedule := expandMonitorScheduledQueryRulesAlertSchedule(d)

	return resourceArmMonitorScheduledQueryRulesCreateUpdate(d, meta, action, schedule)
}

func expandMonitorScheduledQueryRulesAlertingAction(d *schema.ResourceData) *insights.AlertingAction {
	alertActionRaw := d.Get("action").(*schema.Set).List()
	alertAction := expandMonitorScheduledQueryRulesAlertAction(alertActionRaw)
	severityRaw := d.Get("severity").(int)
	severity := strconv.Itoa(severityRaw)
	throttling := d.Get("throttling").(int)

	triggerRaw := d.Get("trigger").(*schema.Set).List()
	trigger := expandMonitorScheduledQueryRulesAlertTrigger(triggerRaw)

	action := insights.AlertingAction{
		AznsAction:      alertAction,
		Severity:        insights.AlertSeverity(severity),
		ThrottlingInMin: utils.Int32(int32(throttling)),
		Trigger:         trigger,
		OdataType:       insights.OdataTypeMicrosoftWindowsAzureManagementMonitoringAlertsModelsMicrosoftAppInsightsNexusDataContractsResourcesScheduledQueryRulesAlertingAction,
	}

	return &action
}

func expandMonitorScheduledQueryRulesAlertAction(input []interface{}) *insights.AzNsActionGroup {
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

func expandMonitorScheduledQueryRulesAlertMetricTrigger(input []interface{}) *insights.LogMetricTrigger {
	if len(input) == 0 {
		return nil
	}

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

func expandMonitorScheduledQueryRulesAlertSchedule(d *schema.ResourceData) *insights.Schedule {
	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)

	schedule := insights.Schedule{
		FrequencyInMinutes:  utils.Int32(int32(frequency)),
		TimeWindowInMinutes: utils.Int32(int32(timeWindow)),
	}

	return &schedule
}

func expandMonitorScheduledQueryRulesAlertTrigger(input []interface{}) *insights.TriggerCondition {
	result := insights.TriggerCondition{}

	for _, item := range input {
		v := item.(map[string]interface{})
		metricTriggerRaw := v["metric_trigger"].(*schema.Set).List()

		result.ThresholdOperator = insights.ConditionalOperator(v["operator"].(string))
		result.Threshold = utils.Float(v["threshold"].(float64))
		result.MetricTrigger = expandMonitorScheduledQueryRulesAlertMetricTrigger(metricTriggerRaw)
	}

	return &result
}

func flattenAzureRmScheduledQueryRulesAlertAction(input *insights.AzNsActionGroup) []interface{} {
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

func flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input *insights.LogMetricTrigger) []interface{} {
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

func flattenAzureRmScheduledQueryRulesAlertTrigger(input *insights.TriggerCondition) []interface{} {
	result := make(map[string]interface{})

	result["operator"] = string(input.ThresholdOperator)

	if input.Threshold != nil {
		result["threshold"] = *input.Threshold
	}

	if input.MetricTrigger != nil {
		result["metric_trigger"] = flattenAzureRmScheduledQueryRulesAlertMetricTrigger(input.MetricTrigger)
	}

	return []interface{}{result}
}
