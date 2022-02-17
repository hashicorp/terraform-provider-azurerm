package logic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/logic/mgmt/2019-05-01/logic"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
			_, err := parse.IntegrationAccountBatchConfigurationID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.IntegrationAccountBatchConfigurationName(),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
										Type:             pluginsdk.TypeString,
										Required:         true,
										DiffSuppressFunc: suppress.CaseDifferenceV2Only,
										ValidateFunc: validation.StringInSlice([]string{
											string(logic.RecurrenceFrequencySecond),
											string(logic.RecurrenceFrequencyMinute),
											string(logic.RecurrenceFrequencyHour),
											string(logic.RecurrenceFrequencyDay),
											string(logic.RecurrenceFrequencyWeek),
											string(logic.RecurrenceFrequencyMonth),
											string(logic.RecurrenceFrequencyYear),
										}, !features.ThreePointOh()),
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
																Type:             pluginsdk.TypeString,
																Required:         true,
																DiffSuppressFunc: suppress.CaseDifferenceV2Only,
																ValidateFunc: validation.StringInSlice([]string{
																	string(logic.DayOfWeekMonday),
																	string(logic.DayOfWeekTuesday),
																	string(logic.DayOfWeekWednesday),
																	string(logic.DayOfWeekThursday),
																	string(logic.DayOfWeekFriday),
																	string(logic.DayOfWeekSaturday),
																	string(logic.DayOfWeekSunday),
																}, !features.ThreePointOh()),
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
													Type:             pluginsdk.TypeSet,
													Optional:         true,
													DiffSuppressFunc: suppress.CaseDifferenceV2Only,
													Elem: &pluginsdk.Schema{
														Type: pluginsdk.TypeString,
														ValidateFunc: validation.StringInSlice([]string{
															string(logic.DaysOfWeekMonday),
															string(logic.DaysOfWeekTuesday),
															string(logic.DaysOfWeekWednesday),
															string(logic.DaysOfWeekThursday),
															string(logic.DaysOfWeekFriday),
															string(logic.DaysOfWeekSaturday),
															string(logic.DaysOfWeekSunday),
														}, !features.ThreePointOh()),
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

	id := parse.NewIntegrationAccountBatchConfigurationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("integration_account_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.BatchConfigurationName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_logic_app_integration_account_batch_configuration", id.ID())
		}
	}

	parameters := logic.BatchConfiguration{
		Properties: &logic.BatchConfigurationProperties{
			BatchGroupName:  utils.String(d.Get("batch_group_name").(string)),
			ReleaseCriteria: expandIntegrationAccountBatchConfigurationBatchReleaseCriteria(d.Get("release_criteria").([]interface{})),
		},
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := v.(map[string]interface{})
		parameters.Properties.Metadata = &metadata
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.IntegrationAccountName, id.BatchConfigurationName, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceLogicAppIntegrationAccountBatchConfigurationRead(d, meta)
}

func resourceLogicAppIntegrationAccountBatchConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountBatchConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.IntegrationAccountName, id.BatchConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.BatchConfigurationName)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("integration_account_name", id.IntegrationAccountName)

	if props := resp.Properties; props != nil {
		d.Set("batch_group_name", props.BatchGroupName)

		if err := d.Set("release_criteria", flattenIntegrationAccountBatchConfigurationBatchReleaseCriteria(props.ReleaseCriteria)); err != nil {
			return fmt.Errorf("setting `release_criteria`: %+v", err)
		}

		if props.Metadata != nil {
			metadata := props.Metadata.(map[string]interface{})
			d.Set("metadata", metadata)
		}
	}

	return nil
}

func resourceLogicAppIntegrationAccountBatchConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Logic.IntegrationAccountBatchConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IntegrationAccountBatchConfigurationID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.IntegrationAccountName, id.BatchConfigurationName); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIntegrationAccountBatchConfigurationBatchReleaseCriteria(input []interface{}) *logic.BatchReleaseCriteria {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := logic.BatchReleaseCriteria{}

	if batchSize := v["batch_size"].(int); batchSize != 0 {
		result.BatchSize = utils.Int32(int32(batchSize))
	}

	if messageCount := v["message_count"].(int); messageCount != 0 {
		result.MessageCount = utils.Int32(int32(messageCount))
	}

	if recurrence := v["recurrence"].([]interface{}); len(recurrence) != 0 {
		result.Recurrence = expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(recurrence)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input []interface{}) *logic.WorkflowTriggerRecurrence {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := logic.WorkflowTriggerRecurrence{
		Frequency: logic.RecurrenceFrequency(v["frequency"].(string)),
		Interval:  utils.Int32(int32(v["interval"].(int))),
	}

	if startTime := v["start_time"].(string); startTime != "" {
		result.StartTime = utils.String(startTime)
	}

	if endTime := v["end_time"].(string); endTime != "" {
		result.EndTime = utils.String(endTime)
	}

	if timeZone := v["time_zone"].(string); timeZone != "" {
		result.TimeZone = utils.String(timeZone)
	}

	if schedule := v["schedule"].([]interface{}); len(schedule) != 0 {
		result.Schedule = expandIntegrationAccountBatchConfigurationRecurrenceSchedule(schedule)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceSchedule(input []interface{}) *logic.RecurrenceSchedule {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	result := logic.RecurrenceSchedule{}

	if hours := v["hours"].(*pluginsdk.Set).List(); len(hours) != 0 {
		result.Hours = utils.ExpandInt32Slice(hours)
	}

	if minutes := v["minutes"].(*pluginsdk.Set).List(); len(minutes) != 0 {
		result.Minutes = utils.ExpandInt32Slice(minutes)
	}

	if rawWeekDays := v["week_days"].(*pluginsdk.Set).List(); len(rawWeekDays) != 0 {
		weekDays := make([]logic.DaysOfWeek, 0)
		for _, item := range *(utils.ExpandStringSlice(rawWeekDays)) {
			weekDays = append(weekDays, (logic.DaysOfWeek)(item))
		}
		result.WeekDays = &weekDays
	}

	if monthDays := v["month_days"].(*pluginsdk.Set).List(); len(monthDays) != 0 {
		result.MonthDays = utils.ExpandInt32Slice(monthDays)
	}

	if monthlyOccurrence := v["monthly"].(*pluginsdk.Set).List(); len(monthlyOccurrence) != 0 {
		result.MonthlyOccurrences = expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(monthlyOccurrence)
	}

	return &result
}

func expandIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrences(input []interface{}) *[]logic.RecurrenceScheduleOccurrence {
	results := make([]logic.RecurrenceScheduleOccurrence, 0)

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, logic.RecurrenceScheduleOccurrence{
			Day:        logic.DayOfWeek(v["weekday"].(string)),
			Occurrence: utils.Int32(int32(v["week"].(int))),
		})
	}

	return &results
}

func flattenIntegrationAccountBatchConfigurationBatchReleaseCriteria(input *logic.BatchReleaseCriteria) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var batchSize int32
	if input.BatchSize != nil {
		batchSize = *input.BatchSize
	}

	var messageCount int32
	if input.MessageCount != nil {
		messageCount = *input.MessageCount
	}

	return []interface{}{
		map[string]interface{}{
			"batch_size":    batchSize,
			"message_count": messageCount,
			"recurrence":    flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input.Recurrence),
		},
	}
}

func flattenIntegrationAccountBatchConfigurationWorkflowTriggerRecurrence(input *logic.WorkflowTriggerRecurrence) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var endTime string
	if input.EndTime != nil {
		endTime = *input.EndTime
	}

	var frequency logic.RecurrenceFrequency
	if input.Frequency != "" {
		frequency = input.Frequency
	}

	var interval int32
	if input.Interval != nil {
		interval = *input.Interval
	}

	var startTime string
	if input.StartTime != nil {
		startTime = *input.StartTime
	}

	var timeZone string
	if input.TimeZone != nil {
		timeZone = *input.TimeZone
	}

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

func flattenIntegrationAccountBatchConfigurationRecurrenceSchedule(input *logic.RecurrenceSchedule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var weekDays []interface{}
	if input.WeekDays != nil {
		weekDaysCast := make([]string, 0)
		for _, item := range *input.WeekDays {
			weekDaysCast = append(weekDaysCast, (string)(item))
		}
		weekDays = utils.FlattenStringSlice(&weekDaysCast)
	}

	return []interface{}{
		map[string]interface{}{
			"hours":      utils.FlattenInt32Slice(input.Hours),
			"minutes":    utils.FlattenInt32Slice(input.Minutes),
			"month_days": utils.FlattenInt32Slice(input.MonthDays),
			"monthly":    flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input.MonthlyOccurrences),
			"week_days":  weekDays,
		},
	}
}

func flattenIntegrationAccountBatchConfigurationRecurrenceScheduleOccurrence(input *[]logic.RecurrenceScheduleOccurrence) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var day logic.DayOfWeek
		if item.Day != "" {
			day = item.Day
		}

		var occurrence int32
		if item.Occurrence != nil {
			occurrence = *item.Occurrence
		}

		results = append(results, map[string]interface{}{
			"weekday": day,
			"week":    occurrence,
		})
	}

	return results
}
