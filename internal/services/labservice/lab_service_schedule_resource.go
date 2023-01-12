package labservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
)

type LabServiceScheduleModel struct {
	Name              string              `tfschema:"name"`
	LabId             string              `tfschema:"lab_id"`
	Notes             string              `tfschema:"notes"`
	RecurrencePattern []RecurrencePattern `tfschema:"recurrence_pattern"`
	StartAt           string              `tfschema:"start_at"`
	StopAt            string              `tfschema:"stop_at"`
	TimeZoneId        string              `tfschema:"time_zone_id"`
}

type RecurrencePattern struct {
	ExpirationDate string                       `tfschema:"expiration_date"`
	Frequency      schedule.RecurrenceFrequency `tfschema:"frequency"`
	Interval       int64                        `tfschema:"interval"`
	WeekDays       []schedule.WeekDay           `tfschema:"week_days"`
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

		"stop_at": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			DiffSuppressFunc: suppress.RFC3339MinuteTime,
			ValidateFunc:     validation.IsRFC3339Time,
		},

		"time_zone_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azValidate.AzureTimeZoneString(),
		},

		"notes": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.ScheduleNotes,
		},

		"recurrence_pattern": {
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

		"start_at": {
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
					StopAt:            model.StopAt,
					TimeZoneId:        model.TimeZoneId,
					RecurrencePattern: expandRecurrencePattern(model.RecurrencePattern),
				},
			}

			if model.Notes != "" {
				properties.Properties.Notes = &model.Notes
			}

			if model.StartAt != "" {
				properties.Properties.StartAt = &model.StartAt
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

			if metadata.ResourceData.HasChange("recurrence_pattern") {
				properties.Properties.RecurrencePattern = expandRecurrencePattern(model.RecurrencePattern)
			}

			if metadata.ResourceData.HasChange("start_at") {
				properties.Properties.StartAt = &model.StartAt
			}

			if metadata.ResourceData.HasChange("stop_at") {
				properties.Properties.StopAt = model.StopAt
			}

			if metadata.ResourceData.HasChange("time_zone_id") {
				properties.Properties.TimeZoneId = model.TimeZoneId
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

			state.StopAt = properties.StopAt
			state.TimeZoneId = properties.TimeZoneId
			state.RecurrencePattern = flattenRecurrencePattern(properties.RecurrencePattern)

			if properties.Notes != nil {
				state.Notes = *properties.Notes
			}

			if properties.StartAt != nil {
				state.StartAt = *properties.StartAt
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

func expandRecurrencePattern(input []RecurrencePattern) *schedule.RecurrencePattern {
	if len(input) == 0 {
		return nil
	}

	recurrencePattern := input[0]
	result := schedule.RecurrencePattern{
		ExpirationDate: recurrencePattern.ExpirationDate,
		Frequency:      recurrencePattern.Frequency,
	}

	if recurrencePattern.Interval != 0 {
		result.Interval = utils.Int64(recurrencePattern.Interval)
	}

	if recurrencePattern.WeekDays != nil {
		result.WeekDays = &recurrencePattern.WeekDays
	}

	return &result
}

func flattenRecurrencePattern(input *schedule.RecurrencePattern) []RecurrencePattern {
	var recurrencePatterns []RecurrencePattern
	if input == nil {
		return recurrencePatterns
	}

	recurrencePattern := RecurrencePattern{
		ExpirationDate: input.ExpirationDate,
		Frequency:      input.Frequency,
	}

	if input.Interval != nil {
		recurrencePattern.Interval = *input.Interval
	}

	if input.WeekDays != nil {
		recurrencePattern.WeekDays = *input.WeekDays
	}

	return append(recurrencePatterns, recurrencePattern)
}
