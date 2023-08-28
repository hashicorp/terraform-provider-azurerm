// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package labservice

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/lab"
	"github.com/hashicorp/go-azure-sdk/resource-manager/labservices/2022-08-01/schedule"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservice/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServiceScheduleModel struct {
	Name       string       `tfschema:"name"`
	LabId      string       `tfschema:"lab_id"`
	Notes      string       `tfschema:"notes"`
	Recurrence []Recurrence `tfschema:"recurrence"`
	StartTime  string       `tfschema:"start_time"`
	StopTime   string       `tfschema:"stop_time"`
	TimeZone   string       `tfschema:"time_zone"`
}

type Recurrence struct {
	ExpirationDate string                       `tfschema:"expiration_date"`
	Frequency      schedule.RecurrenceFrequency `tfschema:"frequency"`
	Interval       int64                        `tfschema:"interval"`
	WeekDays       []string                     `tfschema:"week_days"`
}

type LabServiceScheduleResource struct{}

var _ sdk.ResourceWithUpdate = LabServiceScheduleResource{}

func (r LabServiceScheduleResource) ResourceType() string {
	return "azurerm_lab_service_schedule"
}

func (r LabServiceScheduleResource) ModelObject() interface{} {
	return &LabServiceScheduleModel{}
}

func (r LabServiceScheduleResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return schedule.ValidateScheduleID
}

func (r LabServiceScheduleResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"lab_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: lab.ValidateLabID,
		},

		"stop_time": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"time_zone": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azValidate.AzureTimeZoneString(),
		},

		"notes": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ScheduleNotes,
		},

		"recurrence": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"expiration_date": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						DiffSuppressFunc: suppress.RFC3339MinuteTime,
						ValidateFunc:     validation.IsRFC3339Time,
					},

					"frequency": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(schedule.RecurrenceFrequencyDaily),
							string(schedule.RecurrenceFrequencyWeekly),
						}, false),
					},

					"interval": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntBetween(1, 365),
					},

					"week_days": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(schedule.WeekDaySunday),
								string(schedule.WeekDayMonday),
								string(schedule.WeekDayTuesday),
								string(schedule.WeekDayWednesday),
								string(schedule.WeekDayThursday),
								string(schedule.WeekDayFriday),
								string(schedule.WeekDaySaturday),
							}, false),
						},
					},
				},
			},
		},

		"start_time": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},
	}
}

func (r LabServiceScheduleResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServiceScheduleResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServiceScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LabService.ScheduleClient
			labId, err := lab.ParseLabID(model.LabId)
			if err != nil {
				return err
			}

			id := schedule.NewScheduleID(labId.SubscriptionId, labId.ResourceGroupName, labId.LabName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &schedule.Schedule{
				Properties: schedule.ScheduleProperties{
					StopAt:            model.StopTime,
					TimeZoneId:        model.TimeZone,
					RecurrencePattern: expandRecurrencePattern(model.Recurrence),
				},
			}

			if model.Notes != "" {
				properties.Properties.Notes = &model.Notes
			}

			if model.StartTime != "" {
				properties.Properties.StartAt = &model.StartTime
			}

			if _, err := client.CreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LabServiceScheduleResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.ScheduleClient

			id, err := schedule.ParseScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LabServiceScheduleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("notes") {
				if model.Notes != "" {
					properties.Properties.Notes = &model.Notes
				} else {
					properties.Properties.Notes = nil
				}
			}

			if metadata.ResourceData.HasChange("recurrence") {
				properties.Properties.RecurrencePattern = expandRecurrencePattern(model.Recurrence)
			}

			if metadata.ResourceData.HasChange("start_time") {
				properties.Properties.StartAt = &model.StartTime
			}

			if metadata.ResourceData.HasChange("stop_time") {
				properties.Properties.StopAt = model.StopTime
			}

			if metadata.ResourceData.HasChange("time_zone") {
				properties.Properties.TimeZoneId = model.TimeZone
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LabServiceScheduleResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.ScheduleClient

			id, err := schedule.ParseScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := LabServiceScheduleModel{
				Name:  id.ScheduleName,
				LabId: lab.NewLabID(id.SubscriptionId, id.ResourceGroupName, id.LabName).ID(),
			}

			properties := &model.Properties

			state.StopTime = properties.StopAt
			state.TimeZone = properties.TimeZoneId
			state.Recurrence = flattenRecurrencePattern(properties.RecurrencePattern)

			if properties.Notes != nil {
				state.Notes = *properties.Notes
			}

			if properties.StartAt != nil {
				state.StartTime = *properties.StartAt
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LabServiceScheduleResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LabService.ScheduleClient

			id, err := schedule.ParseScheduleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandRecurrencePattern(input []Recurrence) *schedule.RecurrencePattern {
	if len(input) == 0 {
		return nil
	}

	recurrencePattern := input[0]
	result := schedule.RecurrencePattern{
		ExpirationDate: recurrencePattern.ExpirationDate,
		Frequency:      recurrencePattern.Frequency,
		WeekDays:       expandWeekDays(recurrencePattern.WeekDays),
	}

	if recurrencePattern.Interval != 0 {
		result.Interval = utils.Int64(recurrencePattern.Interval)
	}

	return &result
}

func flattenRecurrencePattern(input *schedule.RecurrencePattern) []Recurrence {
	var recurrencePatterns []Recurrence
	if input == nil {
		return recurrencePatterns
	}

	recurrencePattern := Recurrence{
		ExpirationDate: input.ExpirationDate,
		Frequency:      input.Frequency,
		WeekDays:       flattenWeekDays(input.WeekDays),
	}

	if input.Interval != nil {
		recurrencePattern.Interval = *input.Interval
	}

	return append(recurrencePatterns, recurrencePattern)
}

func expandWeekDays(input []string) *[]schedule.WeekDay {
	if len(input) == 0 {
		return nil
	}

	result := make([]schedule.WeekDay, 0)

	for _, item := range input {
		result = append(result, schedule.WeekDay(item))
	}

	return &result
}

func flattenWeekDays(input *[]schedule.WeekDay) []string {
	if input == nil {
		return nil
	}

	result := make([]string, 0)

	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}
