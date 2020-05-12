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

func resourceArmMonitorActionRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMonitorActionRuleCreateUpdate,
		Read:   resourceArmMonitorActionRuleRead,
		Update: resourceArmMonitorActionRuleCreateUpdate,
		Delete: resourceArmMonitorActionRuleDelete,

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

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(alertsmanagement.TypeActionGroup),
					string(alertsmanagement.TypeDiagnostics),
					string(alertsmanagement.TypeSuppression),
				}, false),
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

			"action_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validate.ActionGroupID,
				ConflictsWith: []string{"suppression"},
			},

			"suppression": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"action_group_id"},
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
											ValidateFunc: validation.StringInSlice(weekDays, false),
										},
									},

									"recurrence_monthly": {
										Type:          schema.TypeSet,
										Optional:      true,
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

			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alert_context": schemaActionRuleAlertContextCondtion(),

						"alert_rule_id": schemaActionRuleAlertRuleIDCondtion(),

						"description": schemaActionRuleDescriptionCondtion(),

						"monitor": schemaActionRuleMonitorCondtion(),

						"monitor_service": schemaActionRuleMonitorServiceCondtion(),

						"severity": schemaActionRuleSeverityCondtion(),

						"target_resource_type": schemaActionRuleTargetResourceTypeCondtion(),
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMonitorActionRuleCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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
			return tf.ImportAsExistsError("azurerm_monitor_action_rule", *existing.ID)
		}
	}

	actionRule := alertsmanagement.ActionRule{
		// the location is always global from the portal
		Location: utils.String(location.Normalize("Global")),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	t := alertsmanagement.Type(d.Get("type").(string))
	description := d.Get("description").(string)
	scope := expandArmActionRuleScope(d.Get("scope").([]interface{}))
	conditions := expandArmActionRuleConditions(d.Get("condition").([]interface{}))

	actionRuleStatus := alertsmanagement.Enabled
	if !d.Get("enabled").(bool) {
		actionRuleStatus = alertsmanagement.Disabled
	}

	switch t {
	case alertsmanagement.TypeSuppression:
		suppressionConfig, err := expandArmActionRuleSuppressionConfig(d.Get("suppression").([]interface{}))
		if err != nil {
			return err
		}
		if suppressionConfig == nil {
			return fmt.Errorf("`suppression` field must be set when type is `Suppression`.")
		}
		actionRule.Properties = &alertsmanagement.Suppression{
			SuppressionConfig: suppressionConfig,
			Scope:             scope,
			Conditions:        conditions,
			Description:       utils.String(description),
			Status:            actionRuleStatus,
			Type:              t,
		}
	case alertsmanagement.TypeActionGroup:
		actionGroupId := d.Get("action_group_id").(string)
		if actionGroupId == "" {
			return fmt.Errorf("`action_group_id` field must be set when type is `ActionGroup`.")
		}
		actionRule.Properties = &alertsmanagement.ActionGroup{
			ActionGroupID: utils.String(actionGroupId),
			Scope:         scope,
			Conditions:    conditions,
			Description:   utils.String(description),
			Status:        actionRuleStatus,
			Type:          t,
		}
	case alertsmanagement.TypeDiagnostics:
		actionRule.Properties = &alertsmanagement.Diagnostics{
			Scope:       scope,
			Conditions:  conditions,
			Description: utils.String(description),
			Status:      actionRuleStatus,
			Type:        t,
		}
	}

	if _, err := client.CreateUpdate(ctx, resourceGroup, name, actionRule); err != nil {
		return fmt.Errorf("creating/updatinge AlertsManagement ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.GetByName(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving AlertsManagement ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for AlertsManagement ActionRule %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)
	return resourceArmMonitorActionRuleRead(d, meta)
}

func resourceArmMonitorActionRuleRead(d *schema.ResourceData, meta interface{}) error {
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
		switch props := resp.Properties.(type) {
		case alertsmanagement.Suppression:
			if err := d.Set("suppression", flattenArmActionRuleSuppression(props.SuppressionConfig)); err != nil {
				return fmt.Errorf("setting suppression: %+v", err)
			}
			if err := d.Set("scope", flattenArmActionRuleScope(props.Scope)); err != nil {
				return fmt.Errorf("setting scope: %+v", err)
			}
			if err := d.Set("condition", flattenArmActionRuleConditions(props.Conditions)); err != nil {
				return fmt.Errorf("setting condition: %+v", err)
			}
			d.Set("type", string(props.Type))
			d.Set("description", props.Description)
			d.Set("enabled", props.Status == alertsmanagement.Enabled)
		case alertsmanagement.ActionGroup:
			if err := d.Set("scope", flattenArmActionRuleScope(props.Scope)); err != nil {
				return fmt.Errorf("setting scope: %+v", err)
			}
			if err := d.Set("condition", flattenArmActionRuleConditions(props.Conditions)); err != nil {
				return fmt.Errorf("setting condition: %+v", err)
			}
			d.Set("type", string(props.Type))
			d.Set("description", props.Description)
			d.Set("enabled", props.Status == alertsmanagement.Enabled)
			d.Set("action_group_id", props.ActionGroupID)
		case alertsmanagement.Diagnostics:
			if err := d.Set("scope", flattenArmActionRuleScope(props.Scope)); err != nil {
				return fmt.Errorf("setting scope: %+v", err)
			}
			if err := d.Set("condition", flattenArmActionRuleConditions(props.Conditions)); err != nil {
				return fmt.Errorf("setting condition: %+v", err)
			}
			d.Set("type", string(props.Type))
			d.Set("description", props.Description)
			d.Set("enabled", props.Status == alertsmanagement.Enabled)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMonitorActionRuleDelete(d *schema.ResourceData, meta interface{}) error {
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

func expandArmActionRuleScope(input []interface{}) *alertsmanagement.Scope {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &alertsmanagement.Scope{
		ScopeType: alertsmanagement.ScopeType(v["type"].(string)),
		Values:    utils.ExpandStringSlice(v["resource_ids"].(*schema.Set).List()),
	}
}

func expandArmActionRuleConditions(input []interface{}) *alertsmanagement.Conditions {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &alertsmanagement.Conditions{
		AlertContext:       expandArmActionRuleCondition(v["alert_context"].([]interface{})),
		AlertRuleID:        expandArmActionRuleCondition(v["alert_rule_id"].([]interface{})),
		Description:        expandArmActionRuleCondition(v["description"].([]interface{})),
		MonitorCondition:   expandArmActionRuleCondition(v["monitor"].([]interface{})),
		MonitorService:     expandArmActionRuleCondition(v["monitor_service"].([]interface{})),
		Severity:           expandArmActionRuleCondition(v["severity"].([]interface{})),
		TargetResourceType: expandArmActionRuleCondition(v["target_resource_type"].([]interface{})),
	}
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

func flattenArmActionRuleScope(input *alertsmanagement.Scope) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var scopeType alertsmanagement.ScopeType
	if input.ScopeType != "" {
		scopeType = input.ScopeType
	}
	return []interface{}{
		map[string]interface{}{
			"type":         scopeType,
			"resource_ids": utils.FlattenStringSlice(input.Values),
		},
	}
}

func flattenArmActionRuleConditions(input *alertsmanagement.Conditions) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	return []interface{}{
		map[string]interface{}{
			"alert_context":        flattenArmActionRuleCondition(input.AlertContext),
			"alert_rule_id":        flattenArmActionRuleCondition(input.AlertRuleID),
			"description":          flattenArmActionRuleCondition(input.Description),
			"monitor":              flattenArmActionRuleCondition(input.MonitorCondition),
			"monitor_service":      flattenArmActionRuleCondition(input.MonitorService),
			"severity":             flattenArmActionRuleCondition(input.Severity),
			"target_resource_type": flattenArmActionRuleCondition(input.TargetResourceType),
		},
	}
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

func FlattenInt32Slice(input *[]int32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
