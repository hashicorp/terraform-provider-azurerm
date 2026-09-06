// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package logic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/logic/2019-05-01/integrationaccountbatchconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLogicAppIntegrationAccountBatchConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceLogicAppIntegrationAccountBatchConfigurationCreateUpdate,
		Read:   resourceLogicAppIntegrationAccountBatchConfigurationRead,
		Update: resourceLogicAppIntegrationAccountBatchConfigurationCreateUpdate,
		Delete: resourceLogicAppIntegrationAccountBatchConfigurationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountBatchConfigurationName(),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"integration_account_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountName(),
			},

			"batch_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountBatchConfigurationBatchGroupName(),
			},

			"release_criteria": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"batch_size": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 83886080),
							AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
						},

						"message_count": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 8000),
							AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
						},

						"recurrence": {
							Type:         pluginsdk.TypeList,
							Optional:     true,
							MaxItems:     1,
							AtLeastOneOf: []string{"release_criteria.0.batch_size", "release_criteria.0.message_count", "release_criteria.0.recurrence"},
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"frequency": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForRecurrenceFrequency(), false),
									},

									"interval": {
										Type:         pluginsdk.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},

									"end_time": {
										Type:             pluginsdk.TypeString,
										Optional:         true,
										DiffSuppressFunc: suppress.RFC3339Time,
										ValidateFunc:     validation.IsRFC3339Time,
									},

									"schedule": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"hours": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeInt,
														ValidateFunc: validation.IntBetween(0, 23),
													},
												},

												"minutes": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeInt,
														ValidateFunc: validation.IntBetween(0, 59),
													},
												},

												"month_days": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeInt,
														ValidateFunc: validation.All(
															validation.IntBetween(-31, 31),
															validation.IntNotInSlice([]int{0}),
														),
													},
													ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
												},

												"monthly": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													Elem: &pluginsdk.Resource{
														Schema: map[string]*pluginsdk.Schema{
															"weekday": {
																Type:         pluginsdk.TypeString,
																Required:     true,
																ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDayOfWeek(), false),
															},

															"week": {
																Type:     pluginsdk.TypeInt,
																Required: true,
																ValidateFunc: validation.All(
																	validation.IntBetween(-5, 5),
																	validation.IntNotInSlice([]int{0}),
																),
															},
														},
													},
													ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.week_days"},
												},

												"week_days": {
													Type:     pluginsdk.TypeSet,
													Optional: true,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeString,
														ValidateFunc: validation.StringInSlice(integrationaccountbatchconfigurations.PossibleValuesForDaysOfWeek(), false),
													},
													ConflictsWith: []string{"release_criteria.0.recurrence.0.schedule.0.month_days", "release_criteria.0.recurrence.0.schedule.0.monthly"},
												},
											},
										},
									},

									"start_time": {
										Type:             pluginsdk.TypeString,
										Optional:         true,
										DiffSuppressFunc: suppress.RFC3339Time,
										ValidateFunc:     validation.IsRFC3339Time,
									},

									"time_zone": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validate.BatchConfigurationRecurrenceTimeZone(),
									},
								},
							},
						},
					},
				},
			},

			"metadata": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
			frequency := strings.ToLower(diff.Get("release_criteria.0.recurrence.0.frequency").(string))

			_, hasWeekDays := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.week_days")
			if hasWeekDays && frequency != "week" {
				return fmt.Errorf("`week_days` can only be set when frequency is `Week`")
			}

			_, hasMonthDays := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.month_days")
			if hasMonthDays && frequency != "month" {
				return fmt.Errorf("`month_days` can only be set when frequency is `Month`")
			}

			_, hasMonthlyOccurrences := diff.GetOk("release_criteria.0.recurrence.0.schedule.0.monthly")
			if hasMonthlyOccurrences && frequency != "month" {
				return fmt.Errorf("`monthly` can only be set when frequency is `Month`")
			}

			return nil
		}),
	}
}

func resourceLogicAppIntegrationAccountBatchConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := integrationaccountbatchconfigurations.NewBatchConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_logic_app_integration_account_batch_configuration", id.ID())
			}
		}
	}

	parameters := integrationaccountbatchconfigurations.BatchConfiguration{
		Properties: integrationaccountbatchconfigurations.BatchConfigurationProperties{
			BatchGroupName:  d.Get("batch_group_name").(string),
			ReleaseCriteria: expandIntegrationAccountBatchConfigurationBatchReleaseCriteria(d.Get("release_criteria").([]interface{})),
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		parameters.Properties.Metadata = &v
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountBatchConfigurationRead(d, meta)
}

func resourceLogicAppIntegrationAccountBatchConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.BatchConfigurationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if model := resp.Model; model != nil {
		props := model.Properties
		d.Set("batch_group_name", props.BatchGroupName)

		if err := d.Set("release_criteria", flattenIntegrationAccountBatchConfigurationBatchReleaseCriteria(props.ReleaseCriteria)); err != nil {
			return fmt.Errorf("setting `release_criteria`: %+v", err)
		}

		if props.Metadata != nil {
			d.Set("metadata", props.Metadata)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountBatchConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationaccountbatchconfigurations.ParseBatchConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountBatchConfigurationBatchReleaseCriteria(input []interface{}) integrationaccountbatchconfigurations.BatchReleaseCriteria {
	result := integrationaccountbatchconfigurations.BatchReleaseCriteria{}
	if len(input) == 0 {
		return result
	}
	v := input[0].(map[string]interface{})

	if batchSize := v["batch_size"].(int); batchSize != 0 {
		result.BatchSize = pointer.To(int64(batchSize))
	}

	if messageCount := v["message_count"].(int); messageCount != 0 {
		result.MessageCount = pointer.To(int64(messageCount))
	}

	if recurrence := v["recurrence"].([]interface{}); len(recurrence) != 0 {
		result.Recurrence = expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(recurrence)
	}

	return result
}

func expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input []interface{}) *integrationaccountbatchconfigurations.WorkflowTriggerRecurrence {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := integrationaccountbatchconfigurations.WorkflowTriggerRecurrence{
		Frequency: pointer.ToEnum[integrationaccountbatchconfigurations.RecurrenceFrequency](v["frequency"].(string)),
		Interval:  pointer.To(int64(v["interval"].(int))),
	}

	if startTime := v["start_time"].(string); startTime != "" {
		result.StartTime = pointer.To(startTime)
	}

	if endTime := v["end_time"].(string); endTime != "" {
		result.EndTime = pointer.To(endTime)
	}

	if timeZone := v["time_zone"].(string); timeZone != "" {
		result.TimeZone = pointer.To(timeZone)
	}

	if schedule := v["schedule"].([]interface{}); len(schedule) != 0 {
		result.Schedule = expandIntegrationAccountBatchConfigurationRecurrenceSchedule(schedule)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceSchedule(input []interface{}) *integrationaccountbatchconfigurations.RecurrenceSchedule {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := integrationaccountbatchconfigurations.RecurrenceSchedule{}

	if hours := v["hours"].(*pluginsdk.Set).List(); len(hours) != 0 {
		result.Hours = helpers.ExpandInt64Slice(hours)
	}

	if minutes := v["minutes"].(*pluginsdk.Set).List(); len(minutes) != 0 {
		result.Minutes = helpers.ExpandInt64Slice(minutes)
	}

	if rawWeekDays := v["week_days"].(*pluginsdk.Set).List(); len(rawWeekDays) != 0 {
		weekDays := make([]integrationaccountbatchconfigurations.DaysOfWeek, 0)
		for _, item := range *helpers.ExpandStringSlice(rawWeekDays) {
			weekDays = append(weekDays, integrationaccountbatchconfigurations.DaysOfWeek(item))
		}
		result.WeekDays = &weekDays
	}

	if monthDays := v["month_days"].(*pluginsdk.Set).List(); len(monthDays) != 0 {
		result.MonthDays = helpers.ExpandInt64Slice(monthDays)
	}

	if monthlyOccurrence := v["monthly"].(*pluginsdk.Set).List(); len(monthlyOccurrence) != 0 {
		result.MonthlyOccurrences = expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(monthlyOccurrence)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(input []interface{}) *[]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence {
	results := make([]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence{
			Day:        pointer.ToEnum[integrationaccountbatchconfigurations.DayOfWeek](v["weekday"].(string)),
			Occurrence: pointer.To(int64(v["week"].(int))),
		})
	}

	return &results
}

func flattenIntegrationAccountBatchConfigurationBatchReleaseCriteria(input integrationaccountbatchconfigurations.BatchReleaseCriteria) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"batch_size":    pointer.From(input.BatchSize),
			"message_count": pointer.From(input.MessageCount),
			"recurrence":    flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input.Recurrence),
		},
	}
}

func flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input *integrationaccountbatchconfigurations.WorkflowTriggerRecurrence) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	endTime := pointer.From(input.EndTime)

	var frequency integrationaccountbatchconfigurations.RecurrenceFrequency
	if input.Frequency != nil && *input.Frequency != "" {
		frequency = *input.Frequency
	}

	interval := pointer.From(input.Interval)

	startTime := pointer.From(input.StartTime)

	timeZone := pointer.From(input.TimeZone)

	return []interface{}{
		map[string]interface{}{
			"end_time":   endTime,
			"frequency":  frequency,
			"interval":   interval,
			"schedule":   flattenIntegrationAccountBatchConfigurationRecurrenceSchedule(input.Schedule),
			"start_time": startTime,
			"time_zone":  timeZone,
		},
	}
}

func flattenIntegrationAccountBatchConfigurationRecurrenceSchedule(input *integrationaccountbatchconfigurations.RecurrenceSchedule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var weekDays []interface{}
	if input.WeekDays != nil {
		weekDaysCast := make([]string, 0)
		for _, item := range *input.WeekDays {
			weekDaysCast = append(weekDaysCast, string(item))
		}
		weekDays = helpers.FlattenStringSlice(&weekDaysCast)
	}

	return []interface{}{
		map[string]interface{}{
			"hours":      helpers.FlattenInt64Slice(input.Hours),
			"minutes":    helpers.FlattenInt64Slice(input.Minutes),
			"month_days": helpers.FlattenInt64Slice(input.MonthDays),
			"monthly":    flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input.MonthlyOccurrences),
			"week_days":  weekDays,
		},
	}
}

func flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input *[]integrationaccountbatchconfigurations.RecurrenceScheduleOccurrence) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var day integrationaccountbatchconfigurations.DayOfWeek
		if item.Day != nil && *item.Day != "" {
			day = *item.Day
		}

		occurrence := pointer.From(item.Occurrence)

		results = append(results, map[string]interface{}{
			"weekday": day,
			"week":    occurrence,
		})
	}

	return results
}
