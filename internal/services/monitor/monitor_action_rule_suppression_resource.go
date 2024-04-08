// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/actionrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorActionRuleSuppression() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorActionRuleSuppressionCreateUpdate,
		Read:   resourceMonitorActionRuleSuppressionRead,
		Update: resourceMonitorActionRuleSuppressionCreateUpdate,
		Delete: resourceMonitorActionRuleSuppressionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		DeprecationMessage: `This resource has been deprecated in favour of the 'azurerm_monitor_alert_processing_rule_suppression' resource and will be removed in v4.0 of the AzureRM Provider`,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := actionrules.ParseActionRuleID(id)
			return err
		}, importMonitorActionRule(actionrules.ActionRuleTypeSuppression)),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ActionRuleName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"suppression": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"recurrence_type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(actionrules.SuppressionTypeAlways),
								string(actionrules.SuppressionTypeOnce),
								string(actionrules.SuppressionTypeDaily),
								string(actionrules.SuppressionTypeWeekly),
								string(actionrules.SuppressionTypeMonthly),
							}, false),
						},

						"schedule": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"start_date_utc": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"end_date_utc": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},

									"recurrence_weekly": {
										Type:          pluginsdk.TypeSet,
										Optional:      true,
										MinItems:      1,
										ConflictsWith: []string{"suppression.0.schedule.0.recurrence_monthly"},
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeString,
											ValidateFunc: validation.IsDayOfTheWeek(false),
										},
									},

									"recurrence_monthly": {
										Type:          pluginsdk.TypeSet,
										Optional:      true,
										MinItems:      1,
										ConflictsWith: []string{"suppression.0.schedule.0.recurrence_weekly"},
										Elem: &pluginsdk.Schema{
											Type:         pluginsdk.TypeInt,
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
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"condition": schemaActionRuleConditions(),

			"scope": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(actionrules.ScopeTypeResourceGroup),
								string(actionrules.ScopeTypeResource),
							}, false),
						},

						"resource_ids": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMonitorActionRuleSuppressionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := actionrules.NewActionRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.GetByName(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_action_rule_suppression", id.ID())
		}
	}

	actionRuleStatus := actionrules.ActionRuleStatusEnabled
	if !d.Get("enabled").(bool) {
		actionRuleStatus = actionrules.ActionRuleStatusDisabled
	}

	suppressionConfig, err := expandActionRuleSuppressionConfig(d.Get("suppression").([]interface{}))
	if err != nil {
		return err
	}

	actionRule := actionrules.ActionRule{
		// the location is always global from the portal
		Location: location.Normalize("Global"),
		Properties: &actionrules.Suppression{
			SuppressionConfig: suppressionConfig,
			Scope:             expandActionRuleScope(d.Get("scope").([]interface{})),
			Conditions:        expandActionRuleConditions(d.Get("condition").([]interface{})),
			Description:       utils.String(d.Get("description").(string)),
			Status:            pointer.To(actionRuleStatus),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateUpdate(ctx, id, actionRule); err != nil {
		return fmt.Errorf("creating/updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorActionRuleSuppressionRead(d, meta)
}

func resourceMonitorActionRuleSuppressionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actionrules.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.GetByName(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ActionRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		props, ok := model.Properties.(actionrules.Suppression)
		if !ok {
			return fmt.Errorf("%s is not of type `Suppression`", id)
		}
		d.Set("description", props.Description)
		if props.Status != nil {
			d.Set("enabled", *props.Status == actionrules.ActionRuleStatusEnabled)
		}
		if err := d.Set("suppression", flattenActionRuleSuppression(props.SuppressionConfig)); err != nil {
			return fmt.Errorf("setting suppression: %+v", err)
		}
		if err := d.Set("scope", flattenActionRuleScope(props.Scope)); err != nil {
			return fmt.Errorf("setting scope: %+v", err)
		}
		if err := d.Set("condition", flattenActionRuleConditions(props.Conditions)); err != nil {
			return fmt.Errorf("setting condition: %+v", err)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceMonitorActionRuleSuppressionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ActionRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := actionrules.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func expandActionRuleSuppressionConfig(input []interface{}) (*actionrules.SuppressionConfig, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})
	recurrenceType := actionrules.SuppressionType(v["recurrence_type"].(string))
	schedule, err := expandActionRuleSuppressionSchedule(v["schedule"].([]interface{}), recurrenceType)
	if err != nil {
		return nil, err
	}
	if recurrenceType != actionrules.SuppressionTypeAlways && schedule == nil {
		return nil, fmt.Errorf("`schedule` block must be set when `recurrence_type` is Once, Daily, Weekly or Monthly.")
	}
	return &actionrules.SuppressionConfig{
		RecurrenceType: recurrenceType,
		Schedule:       schedule,
	}, nil
}

func expandActionRuleSuppressionSchedule(input []interface{}, suppressionType actionrules.SuppressionType) (*actionrules.SuppressionSchedule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	v := input[0].(map[string]interface{})

	var recurrence []interface{}
	switch suppressionType {
	case actionrules.SuppressionTypeWeekly:
		if recurrenceWeekly, ok := v["recurrence_weekly"]; ok {
			recurrence = expandActionRuleSuppressionScheduleRecurrenceWeekly(recurrenceWeekly.(*pluginsdk.Set).List())
		}
		if len(recurrence) == 0 {
			return nil, fmt.Errorf("`recurrence_weekly` must be set and should have at least one element when `recurrence_type` is Weekly.")
		}
	case actionrules.SuppressionTypeMonthly:
		if recurrenceMonthly, ok := v["recurrence_monthly"]; ok {
			recurrence = recurrenceMonthly.(*pluginsdk.Set).List()
		}
		if len(recurrence) == 0 {
			return nil, fmt.Errorf("`recurrence_monthly` must be set and should have at least one element when `recurrence_type` is Monthly.")
		}
	}

	startDateUTC, _ := time.Parse(time.RFC3339, v["start_date_utc"].(string))
	endDateUTC, _ := time.Parse(time.RFC3339, v["end_date_utc"].(string))
	return &actionrules.SuppressionSchedule{
		StartDate:        utils.String(startDateUTC.Format(scheduleDateLayout)),
		EndDate:          utils.String(endDateUTC.Format(scheduleDateLayout)),
		StartTime:        utils.String(startDateUTC.Format(scheduleTimeLayout)),
		EndTime:          utils.String(endDateUTC.Format(scheduleTimeLayout)),
		RecurrenceValues: utils.ExpandInt64Slice(recurrence),
	}, nil
}

func expandActionRuleSuppressionScheduleRecurrenceWeekly(input []interface{}) []interface{} {
	result := make([]interface{}, 0, len(input))
	for _, v := range input {
		result = append(result, weekDayMap[v.(string)])
	}
	return result
}

func flattenActionRuleSuppression(input *actionrules.SuppressionConfig) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var recurrenceType actionrules.SuppressionType
	if input.RecurrenceType != "" {
		recurrenceType = input.RecurrenceType
	}
	return []interface{}{
		map[string]interface{}{
			"recurrence_type": string(recurrenceType),
			"schedule":        flattenActionRuleSuppressionSchedule(input.Schedule, recurrenceType),
		},
	}
}

func flattenActionRuleSuppressionSchedule(input *actionrules.SuppressionSchedule, recurrenceType actionrules.SuppressionType) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	startDateUTCStr := ""
	endDateUTCStr := ""
	recurrenceWeekly := []interface{}{}
	recurrenceMonthly := []interface{}{}

	if input.StartDate != nil && input.StartTime != nil {
		date, _ := time.ParseInLocation(scheduleDateTimeLayout, fmt.Sprintf("%s %s", *input.StartDate, *input.StartTime), time.UTC)
		startDateUTCStr = date.Format(time.RFC3339)
	}
	if input.EndDate != nil && input.EndTime != nil {
		date, _ := time.ParseInLocation(scheduleDateTimeLayout, fmt.Sprintf("%s %s", *input.EndDate, *input.EndTime), time.UTC)
		endDateUTCStr = date.Format(time.RFC3339)
	}

	if recurrenceType == actionrules.SuppressionTypeWeekly {
		recurrenceWeekly = flattenActionRuleSuppressionScheduleRecurrenceWeekly(input.RecurrenceValues)
	}
	if recurrenceType == actionrules.SuppressionTypeMonthly {
		recurrenceMonthly = utils.FlattenInt64Slice(input.RecurrenceValues)
	}
	return []interface{}{
		map[string]interface{}{
			"start_date_utc":     startDateUTCStr,
			"end_date_utc":       endDateUTCStr,
			"recurrence_weekly":  recurrenceWeekly,
			"recurrence_monthly": recurrenceMonthly,
		},
	}
}

func flattenActionRuleSuppressionScheduleRecurrenceWeekly(input *[]int64) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, weekDays[int(item)])
		}
	}
	return result
}
