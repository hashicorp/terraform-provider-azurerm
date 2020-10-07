package monitor

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorScheduledQueryRulesAlert() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorScheduledQueryRulesAlertCreateUpdate,
		Read:   resourceArmMonitorScheduledQueryRulesAlertRead,
		Update: resourceArmMonitorScheduledQueryRulesAlertCreateUpdate,
		Delete: resourceArmMonitorScheduledQueryRulesAlertDelete,
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
				ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"authorized_resource_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 100,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
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
							Type:     schema.TypeString,
							Optional: true,
							// TODO remove `Computed: true` in 3.0. This is a breaking change where the Default used to be "{}"
							// We'll keep Computed: true for users who expect the same functionality but will remove it in 3.0
							Computed:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"email_subject": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 1440),
			},
			"query": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4),
			},
			"throttling": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 10000),
			},
			"time_window": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 2880),
			},
			"trigger": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_trigger": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_column": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
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
										ValidateFunc: validateMonitorScheduledQueryRulesAlertThreshold,
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
							Type:         schema.TypeFloat,
							Required:     true,
							ValidateFunc: validateMonitorScheduledQueryRulesAlertThreshold,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorScheduledQueryRulesAlertCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	action := expandMonitorScheduledQueryRulesAlertingAction(d)
	schedule := expandMonitorScheduledQueryRulesAlertSchedule(d)
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)
	if timeWindow < frequency {
		return fmt.Errorf("Error in parameter values for Scheduled Query Rules %q (Resource Group %q): time_window must be greater than or equal to frequency", name, resourceGroup)
	}

	query := d.Get("query").(string)
	_, ok := d.GetOk("metric_trigger")
	if ok {
		if !(strings.Contains(query, "summarize") &&
			strings.Contains(query, "AggregatedValue") &&
			strings.Contains(query, "bin")) {
			return fmt.Errorf("Error in parameter values for Scheduled Query Rules %q (Resource Group %q): query must contain summarize, AggregatedValue, and bin when metric_trigger is specified", name, resourceGroup)
		}
	}

	if d.IsNewResource() {
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

	source := expandMonitorScheduledQueryRulesCommonSource(d)

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
		return fmt.Errorf("Error creating or updating Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Scheduled query rule %q (resource group %q) ID is empty", name, resourceGroup)
	}
	d.SetId(*read.ID)

	return resourceArmMonitorScheduledQueryRulesAlertRead(d, meta)
}

func resourceArmMonitorScheduledQueryRulesAlertRead(d *schema.ResourceData, meta interface{}) error {
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
			log.Printf("[DEBUG] Scheduled Query Rule %q was not found in Resource Group %q", name, resourceGroup)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
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

	action, ok := resp.Action.(insights.AlertingAction)
	if !ok {
		return fmt.Errorf("Wrong action type in Scheduled Query Rule %q (resource group %q): %T", name, resourceGroup, resp.Action)
	}
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

func resourceArmMonitorScheduledQueryRulesAlertDelete(d *schema.ResourceData, meta interface{}) error {
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
			return fmt.Errorf("Error deleting Scheduled Query Rule %q (resource group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesAlertingAction(d *schema.ResourceData) *insights.AlertingAction {
	alertActionRaw := d.Get("action").([]interface{})
	alertAction := expandMonitorScheduledQueryRulesAlertAction(alertActionRaw)
	severityRaw := d.Get("severity").(int)
	severity := strconv.Itoa(severityRaw)
	throttling := d.Get("throttling").(int)

	triggerRaw := d.Get("trigger").([]interface{})
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

	if len(input) == 0 {
		return &result
	}
	for _, item := range input {
		if item == nil {
			continue
		}

		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		actionGroups := v["action_group"].(*schema.Set).List()
		result.ActionGroup = utils.ExpandStringSlice(actionGroups)
		result.EmailSubject = utils.String(v["email_subject"].(string))
		if v := v["custom_webhook_payload"].(string); v != "" {
			result.CustomWebhookPayload = utils.String(v)
		}
	}

	return &result
}

func expandMonitorScheduledQueryRulesAlertMetricTrigger(input []interface{}) *insights.LogMetricTrigger {
	if len(input) == 0 {
		return nil
	}

	result := insights.LogMetricTrigger{}
	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
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
	if len(input) == 0 {
		return &result
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		metricTriggerRaw := v["metric_trigger"].([]interface{})

		result.ThresholdOperator = insights.ConditionalOperator(v["operator"].(string))
		result.Threshold = utils.Float(v["threshold"].(float64))
		result.MetricTrigger = expandMonitorScheduledQueryRulesAlertMetricTrigger(metricTriggerRaw)
	}

	return &result
}
