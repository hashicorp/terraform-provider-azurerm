package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-05-05/alertsmanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMonitorActionRuleSuppression() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorActionRuleSuppressionCreateUpdate,
		Read:   resourceArmMonitorActionRuleSuppressionRead,
		Update: resourceArmMonitorActionRuleSuppressionCreateUpdate,
		Delete: resourceArmMonitorActionRuleSuppressionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.ActionRuleID(id)
			return err
		}),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ActionRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"suppression": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recurrence_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(alertsmanagement.Always),
								string(alertsmanagement.Once),
								string(alertsmanagement.Daily),
								string(alertsmanagement.Weekly),
								string(alertsmanagement.Monthly),
							}, false),
						},

						"schedule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_date": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ActionRuleScheduleDate,
									},

									"end_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ActionRuleScheduleTime,
									},

									"start_date": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ActionRuleScheduleDate,
									},

									"start_time": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.ActionRuleScheduleTime,
									},

									"recurrence_weekly": {
										Type:          schema.TypeSet,
										Optional:      true,
										MinItems:      1,
										ConflictsWith: []string{"suppression.0.schedule.0.recurrence_monthly"},
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.IsDayOfTheWeek(false),
										},
									},

									"recurrence_monthly": {
										Type:          schema.TypeSet,
										Optional:      true,
										MinItems:      1,
										ConflictsWith: []string{"suppression.0.schedule.0.recurrence_weekly"},
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(1, 31),
										},
									},
								},
							},
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

			"scope": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(alertsmanagement.ScopeTypeResourceGroup),
								string(alertsmanagement.ScopeTypeResource),
							}, false),
						},

						"resource_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
					},
				},
			},

			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_context": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
								string(alertsmanagement.Contains),
								string(alertsmanagement.DoesNotContain),
							}, false),
							validation.StringIsNotEmpty,
						),

						"alert_rule_id": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
								string(alertsmanagement.Contains),
								string(alertsmanagement.DoesNotContain),
							}, false),
							validation.StringIsNotEmpty,
						),

						"description": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
								string(alertsmanagement.Contains),
								string(alertsmanagement.DoesNotContain),
							}, false),
							validation.StringIsNotEmpty,
						),

						"monitor": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
							}, false),
							validation.StringInSlice([]string{
								string(alertsmanagement.Fired),
								string(alertsmanagement.Resolved),
							}, false),
						),

						"monitor_service": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
							}, false),
							// the supported type list is not consistent with the swagger and sdk
							// https://github.com/Azure/azure-rest-api-specs/issues/9076
							// directly use string constant
							validation.StringInSlice([]string{
								"ActivityLog Administrative",
								"ActivityLog Autoscale",
								"ActivityLog Policy",
								"ActivityLog Recommendation",
								"ActivityLog Security",
								"Application Insights",
								"Azure Backup",
								"Data Box Edge",
								"Data Box Gateway",
								"Health Platform",
								"Log Analytics",
								"Platform",
								"Resource Health",
							}, false),
						),

						"severity": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
							}, false),
							validation.StringInSlice([]string{
								string(alertsmanagement.Sev0),
								string(alertsmanagement.Sev1),
								string(alertsmanagement.Sev2),
								string(alertsmanagement.Sev3),
								string(alertsmanagement.Sev4),
							}, false),
						),

						"target_resource_type": schemaActionRuleCondition(
							validation.StringInSlice([]string{
								string(alertsmanagement.Equals),
								string(alertsmanagement.NotEquals),
							}, false),
							validation.StringIsNotEmpty,
						),
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorActionRuleSuppressionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	name := d.Get("name").(string)

	if d.IsNewResource() {
		existing, err := client.GetByName(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for present of existing Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_action_rule_suppression", *existing.ID)
		}
	}

	actionRuleStatus := alertsmanagement.Enabled
	if !d.Get("enabled").(bool) {
		actionRuleStatus = alertsmanagement.Disabled
	}

	suppressionConfig, err := expandArmActionRuleSuppressionConfig(d.Get("suppression").([]interface{}))
	if err != nil {
		return err
	}

	actionRule := alertsmanagement.ActionRule{
		// the location is always global from the portal
		Location: utils.String(location.Normalize("Global")),
		Properties: &alertsmanagement.Suppression{
			SuppressionConfig: suppressionConfig,
			Scope:             expandArmActionRuleScope(d.Get("scope").([]interface{})),
			Conditions:        expandArmActionRuleConditions(d.Get("condition").([]interface{})),
			Description:       utils.String(d.Get("description").(string)),
			Status:            actionRuleStatus,
			Type:              alertsmanagement.TypeSuppression,
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateUpdate(ctx, resourceGroup, name, actionRule); err != nil {
		return fmt.Errorf("creating/updatinge Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetByName(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Monitor ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmMonitorActionRuleSuppressionRead(d, meta)
}

func resourceArmMonitorActionRuleSuppressionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ActionRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByName(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Action Rule %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving monitor ActionRule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if resp.Properties != nil {
		props, _ := resp.Properties.AsSuppression()
		d.Set("description", props.Description)
		d.Set("enabled", props.Status == alertsmanagement.Enabled)
		if err := d.Set("suppression", flattenArmActionRuleSuppression(props.SuppressionConfig)); err != nil {
			return fmt.Errorf("setting suppression: %+v", err)
		}
		if err := d.Set("scope", flattenArmActionRuleScope(props.Scope)); err != nil {
			return fmt.Errorf("setting scope: %+v", err)
		}
		if err := d.Set("condition", flattenArmActionRuleConditions(props.Conditions)); err != nil {
			return fmt.Errorf("setting condition: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorActionRuleSuppressionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ActionRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Name); err != nil {
		return fmt.Errorf("deleting monitor ActionRule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandArmActionRuleSuppressionConfig(input []interface{}) (*alertsmanagement.SuppressionConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})
	recurrenceType := alertsmanagement.SuppressionType(v["recurrence_type"].(string))
	schedule, err := expandArmActionRuleSuppressionSchedule(v["schedule"].([]interface{}), recurrenceType)
	if err != nil {
		return nil, err
	}
	if recurrenceType != alertsmanagement.Always && schedule == nil {
		return nil, fmt.Errorf("`schedule` block must be set when `recurrence_type` is Once, Daily, Weekly or Monthly.")
	}
	return &alertsmanagement.SuppressionConfig{
		RecurrenceType: recurrenceType,
		Schedule:       schedule,
	}, nil
}

func expandArmActionRuleSuppressionSchedule(input []interface{}, suppressionType alertsmanagement.SuppressionType) (*alertsmanagement.SuppressionSchedule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})

	var recurrence []interface{}
	switch suppressionType {
	case alertsmanagement.Weekly:
		if recurrenceWeekly, ok := v["recurrence_weekly"]; ok {
			recurrence = expandArmActionRuleSuppressionScheduleRecurrenceWeekly(recurrenceWeekly.(*schema.Set).List())
		}
		if len(recurrence) == 0 {
			return nil, fmt.Errorf("`recurrence_weekly` must be set and should have at least one element when `recurrence_type` is Weekly.")
		}
	case alertsmanagement.Monthly:
		if recurrenceMonthly, ok := v["recurrence_monthly"]; ok {
			recurrence = recurrenceMonthly.(*schema.Set).List()
		}
		if len(recurrence) == 0 {
			return nil, fmt.Errorf("`recurrence_monthly` must be set and should have at least one element when `recurrence_type` is Monthly.")
		}
	}

	recurrenceValues := make([]int32, len(recurrence))
	for i, item := range recurrence {
		recurrenceValues[i] = int32(item.(int))
	}

	return &alertsmanagement.SuppressionSchedule{
		StartDate:        utils.String(v["start_date"].(string)),
		EndDate:          utils.String(v["end_date"].(string)),
		StartTime:        utils.String(v["start_time"].(string)),
		EndTime:          utils.String(v["end_time"].(string)),
		RecurrenceValues: &recurrenceValues,
	}, nil
}

func expandArmActionRuleSuppressionScheduleRecurrenceWeekly(input []interface{}) []interface{} {
	result := make([]interface{}, 0, len(input))
	for _, v := range input {
		result = append(result, weekDayMap[v.(string)])
	}
	return result
}

func flattenArmActionRuleSuppression(input *alertsmanagement.SuppressionConfig) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var recurrenceType alertsmanagement.SuppressionType
	if input.RecurrenceType != "" {
		recurrenceType = input.RecurrenceType
	}
	return []interface{}{
		map[string]interface{}{
			"recurrence_type": string(recurrenceType),
			"schedule":        flattenArmActionRuleSuppressionSchedule(input.Schedule, recurrenceType),
		},
	}
}

func flattenArmActionRuleSuppressionSchedule(input *alertsmanagement.SuppressionSchedule, recurrenceType alertsmanagement.SuppressionType) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	startDate := ""
	startTime := ""
	endDate := ""
	endTime := ""
	recurrenceWeekly := []interface{}{}
	recurrenceMonthly := []interface{}{}

	if input.StartDate != nil {
		startDate = *input.StartDate
	}
	if input.StartTime != nil {
		startTime = *input.StartTime
	}
	if input.EndDate != nil {
		endDate = *input.EndDate
	}
	if input.EndTime != nil {
		endTime = *input.EndTime
	}
	if recurrenceType == alertsmanagement.Weekly {
		recurrenceWeekly = flattenArmActionRuleSuppressionScheduleRecurrenceWeekly(input.RecurrenceValues)
	}
	if recurrenceType == alertsmanagement.Monthly {
		recurrenceMonthly = FlattenInt32Slice(input.RecurrenceValues)
	}
	return []interface{}{
		map[string]interface{}{
			"start_date":         startDate,
			"start_time":         startTime,
			"end_date":           endDate,
			"end_time":           endTime,
			"recurrence_weekly":  recurrenceWeekly,
			"recurrence_monthly": recurrenceMonthly,
		},
	}
}

func flattenArmActionRuleSuppressionScheduleRecurrenceWeekly(input *[]int32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, weekDays[int(item)])
		}
	}
	return result
}
