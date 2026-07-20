// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package costmanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/scheduledactions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2023-08-01/views"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CostManagementScheduledActionResource struct{}

type CostManagementScheduledActionModel struct {
	Name               string   `tfschema:"name"`
	DisplayName        string   `tfschema:"display_name"`
	ViewId             string   `tfschema:"view_id"`
	EmailSubject       string   `tfschema:"email_subject"`
	EmailAddresses     []string `tfschema:"email_addresses"`
	Message            string   `tfschema:"message"`
	EmailAddressSender string   `tfschema:"email_address_sender"`
	Frequency          string   `tfschema:"frequency"`
	DaysOfWeek         []string `tfschema:"days_of_week"`
	WeeksOfMonth       []string `tfschema:"weeks_of_month"`
	HourOfDay          int64    `tfschema:"hour_of_day"`
	DayOfMonth         int64    `tfschema:"day_of_month"`
	StartDate          string   `tfschema:"start_date"`
	EndDate            string   `tfschema:"end_date"`
}

var _ sdk.Resource = CostManagementScheduledActionResource{}

func (r CostManagementScheduledActionResource) Arguments() map[string]*pluginsdk.Schema {
	resource := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// adding 2026/04 api limitations behind 5.0 flag incase existing resources are still allowed with previous limits
			ValidateFunc: func() pluginsdk.SchemaValidateFunc {
				if features.FivePointOh() {
					return validation.StringLenBetween(1, 25)
				}
				return validation.StringIsNotEmpty
			}(),
		},

		"view_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"email_subject": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"email_addresses": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 20,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"message": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 250),
		},

		"email_address_sender": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"frequency": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(scheduledactions.PossibleValuesForScheduleFrequency(), false),
		},

		"days_of_week": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(scheduledactions.PossibleValuesForDaysOfWeek(), false),
			},
		},

		"weeks_of_month": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(scheduledactions.PossibleValuesForWeeksOfMonth(), false),
			},
		},

		"hour_of_day": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 23),
		},

		"day_of_month": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 31),
		},

		"start_date": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},

		"end_date": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.IsRFC3339Time,
		},
	}

	if !features.FivePointOh() {
		resource["email_subject"].ValidateFunc = validation.StringLenBetween(1, 70)
	}

	return resource
}

func (r CostManagementScheduledActionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CostManagementScheduledActionResource) ModelObject() interface{} {
	return &CostManagementScheduledActionModel{}
}

func (r CostManagementScheduledActionResource) ResourceType() string {
	return "azurerm_cost_management_scheduled_action"
}

func (r CostManagementScheduledActionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return scheduledactions.ValidateScopedScheduledActionID
}

func (r CostManagementScheduledActionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			var config CostManagementScheduledActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			viewId, err := views.ParseScopedViewID(config.ViewId)
			if err != nil {
				return err
			}
			id := scheduledactions.NewScopedScheduledActionID(viewId.Scope, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.GetByScope(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			var daysOfWeek []scheduledactions.DaysOfWeek
			if len(config.DaysOfWeek) > 0 {
				daysOfWeek = make([]scheduledactions.DaysOfWeek, 0)
				for _, value := range config.DaysOfWeek {
					daysOfWeek = append(daysOfWeek, scheduledactions.DaysOfWeek(value))
				}
			}

			var weeksOfMonth []scheduledactions.WeeksOfMonth
			if len(config.WeeksOfMonth) > 0 {
				weeksOfMonth = make([]scheduledactions.WeeksOfMonth, 0)
				for _, value := range config.WeeksOfMonth {
					weeksOfMonth = append(weeksOfMonth, scheduledactions.WeeksOfMonth(value))
				}
			}

			schedule := scheduledactions.ScheduleProperties{
				Frequency:    scheduledactions.ScheduleFrequency(config.Frequency),
				WeeksOfMonth: &weeksOfMonth,
				DaysOfWeek:   &daysOfWeek,
				HourOfDay:    pointer.To(config.HourOfDay),
				DayOfMonth:   pointer.To(config.DayOfMonth),
				StartDate:    config.StartDate,
				EndDate:      config.EndDate,
			}

			props := scheduledactions.ScheduledAction{
				Kind: pointer.To(scheduledactions.ScheduledActionKindEmail),
				Properties: &scheduledactions.ScheduledActionProperties{
					DisplayName: config.DisplayName,
					Status:      scheduledactions.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactions.FileDestination{
						FileFormats: &[]scheduledactions.FileFormat{},
					},
					NotificationEmail: pointer.To(config.EmailAddressSender),
					Notification: scheduledactions.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}

			if _, err = client.CreateOrUpdateByScope(ctx, id, props, scheduledactions.CreateOrUpdateByScopeOperationOptions{}); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r CostManagementScheduledActionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.GetByScope(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("reading %s: %+v", *id, err)
			}

			state := CostManagementScheduledActionModel{
				Name: id.ScheduledActionName,
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = props.DisplayName
					state.EmailAddressSender = pointer.From(props.NotificationEmail)

					viewId, err := views.ParseScopedViewID(props.ViewId)
					if err != nil {
						return err
					}
					state.ViewId = viewId.ID()

					state.EmailSubject = props.Notification.Subject
					state.EmailAddresses = props.Notification.To
					state.Message = pointer.From(props.Notification.Message)

					state.Frequency = string(props.Schedule.Frequency)
					if props.Schedule.DaysOfWeek != nil {
						daysOfWeek := make([]string, 0)
						for _, v := range *props.Schedule.DaysOfWeek {
							daysOfWeek = append(daysOfWeek, string(v))
						}
						state.DaysOfWeek = daysOfWeek
					}
					if props.Schedule.WeeksOfMonth != nil {
						weeksOfMonth := make([]string, 0)
						for _, v := range *props.Schedule.WeeksOfMonth {
							weeksOfMonth = append(weeksOfMonth, string(v))
						}
						state.WeeksOfMonth = weeksOfMonth
					}
					state.HourOfDay = pointer.From(props.Schedule.HourOfDay)
					state.DayOfMonth = pointer.From(props.Schedule.DayOfMonth)
					state.StartDate = props.Schedule.StartDate
					state.EndDate = props.Schedule.EndDate
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r CostManagementScheduledActionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err = client.DeleteByScope(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r CostManagementScheduledActionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionsClient

			id, err := scheduledactions.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config CostManagementScheduledActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Update operation requires latest eTag to be set in the request.
			existing, err := client.GetByScope(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading %s: %+v", *id, err)
			}
			if model := existing.Model; model != nil {
				if model.ETag == nil {
					return fmt.Errorf("updating %s: eTag was nil", *id)
				}

				if model.Properties == nil {
					return fmt.Errorf("updating %s: properties was nil", *id)
				}

				if metadata.ResourceData.HasChange("display_name") {
					model.Properties.DisplayName = config.DisplayName
				}

				if metadata.ResourceData.HasChange("view_id") {
					viewId, err := views.ParseScopedViewID(config.ViewId)
					if err != nil {
						return err
					}
					model.Properties.ViewId = viewId.ID()
				}

				if metadata.ResourceData.HasChange("email_address_sender") {
					model.Properties.NotificationEmail = pointer.To(config.EmailAddressSender)
				}

				if metadata.ResourceData.HasChange("email_subject") {
					model.Properties.Notification.Subject = config.EmailSubject
				}

				if metadata.ResourceData.HasChange("email_addresses") {
					model.Properties.Notification.To = config.EmailAddresses
				}

				if metadata.ResourceData.HasChange("message") {
					model.Properties.Notification.Message = pointer.To(config.Message)
				}

				if metadata.ResourceData.HasChange("frequency") {
					model.Properties.Schedule.Frequency = scheduledactions.ScheduleFrequency(config.Frequency)
				}

				if metadata.ResourceData.HasChange("days_of_week") {
					var daysOfWeek []scheduledactions.DaysOfWeek
					if len(config.DaysOfWeek) > 0 {
						daysOfWeek = make([]scheduledactions.DaysOfWeek, 0)
						for _, value := range config.DaysOfWeek {
							daysOfWeek = append(daysOfWeek, scheduledactions.DaysOfWeek(value))
						}
					}
					model.Properties.Schedule.DaysOfWeek = &daysOfWeek
				}

				if metadata.ResourceData.HasChange("weeks_of_month") {
					var weeksOfMonth []scheduledactions.WeeksOfMonth
					if len(config.WeeksOfMonth) > 0 {
						weeksOfMonth = make([]scheduledactions.WeeksOfMonth, 0)
						for _, value := range config.WeeksOfMonth {
							weeksOfMonth = append(weeksOfMonth, scheduledactions.WeeksOfMonth(value))
						}
					}
					model.Properties.Schedule.WeeksOfMonth = &weeksOfMonth
				}

				if metadata.ResourceData.HasChange("start_date") {
					model.Properties.Schedule.StartDate = config.StartDate
				}

				if metadata.ResourceData.HasChange("end_date") {
					model.Properties.Schedule.EndDate = config.EndDate
				}

				if metadata.ResourceData.HasChange("hour_of_day") {
					model.Properties.Schedule.HourOfDay = pointer.To(config.HourOfDay)
				}

				if metadata.ResourceData.HasChange("day_of_month") {
					model.Properties.Schedule.DayOfMonth = pointer.To(config.DayOfMonth)
				}

				if _, err = client.CreateOrUpdateByScope(ctx, *id, *model, scheduledactions.CreateOrUpdateByScopeOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
