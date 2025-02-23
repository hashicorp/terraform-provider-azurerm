// Copyright (c) HashiCorp, Inc.
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CostManagementScheduledActionResource struct{}

var _ sdk.Resource = CostManagementScheduledActionResource{}

func (r CostManagementScheduledActionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotWhiteSpace,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
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
			ValidateFunc: validation.StringLenBetween(1, 70),
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
}

func (r CostManagementScheduledActionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r CostManagementScheduledActionResource) ModelObject() interface{} {
	return nil
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

			viewId, err := views.ParseScopedViewID(metadata.ResourceData.Get("view_id").(string))
			if err != nil {
				return err
			}
			id := scheduledactions.NewScopedScheduledActionID(viewId.Scope, metadata.ResourceData.Get("name").(string))

			existing, err := client.GetByScope(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError(r.ResourceType(), id.ID())
			}

			var daysOfWeek []scheduledactions.DaysOfWeek
			if len(metadata.ResourceData.Get("days_of_week").([]interface{})) > 0 {
				daysOfWeek = make([]scheduledactions.DaysOfWeek, 0)
				for _, value := range metadata.ResourceData.Get("days_of_week").([]interface{}) {
					daysOfWeek = append(daysOfWeek, scheduledactions.DaysOfWeek(value.(string)))
				}
			}

			var weeksOfMonth []scheduledactions.WeeksOfMonth
			if len(metadata.ResourceData.Get("weeks_of_month").([]interface{})) > 0 {
				weeksOfMonth = make([]scheduledactions.WeeksOfMonth, 0)
				for _, value := range metadata.ResourceData.Get("weeks_of_month").([]interface{}) {
					weeksOfMonth = append(weeksOfMonth, scheduledactions.WeeksOfMonth(value.(string)))
				}
			}
			schedule := scheduledactions.ScheduleProperties{
				Frequency:    scheduledactions.ScheduleFrequency(metadata.ResourceData.Get("frequency").(string)),
				WeeksOfMonth: &weeksOfMonth,
				DaysOfWeek:   &daysOfWeek,
				HourOfDay:    utils.Int64(int64(metadata.ResourceData.Get("hour_of_day").(int))),
				DayOfMonth:   utils.Int64(int64(metadata.ResourceData.Get("day_of_month").(int))),
				StartDate:    metadata.ResourceData.Get("start_date").(string),
				EndDate:      metadata.ResourceData.Get("end_date").(string),
			}

			props := scheduledactions.ScheduledAction{
				Kind: pointer.To(scheduledactions.ScheduledActionKindEmail),
				Properties: &scheduledactions.ScheduledActionProperties{
					DisplayName: metadata.ResourceData.Get("display_name").(string),
					Status:      scheduledactions.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactions.FileDestination{
						FileFormats: &[]scheduledactions.FileFormat{},
					},
					NotificationEmail: utils.String(metadata.ResourceData.Get("email_address_sender").(string)),
					Notification: scheduledactions.NotificationProperties{
						Subject: metadata.ResourceData.Get("email_subject").(string),
						Message: utils.String(metadata.ResourceData.Get("message").(string)),
						To:      *utils.ExpandStringSlice(metadata.ResourceData.Get("email_addresses").([]interface{})),
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

			metadata.ResourceData.Set("name", id.ScheduledActionName)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("display_name", props.DisplayName)
					metadata.ResourceData.Set("email_address_sender", props.NotificationEmail)
					viewId, err := views.ParseScopedViewID(props.ViewId)
					if err != nil {
						return err
					}
					metadata.ResourceData.Set("view_id", viewId.ID())

					metadata.ResourceData.Set("email_subject", props.Notification.Subject)
					metadata.ResourceData.Set("email_addresses", props.Notification.To)
					metadata.ResourceData.Set("message", props.Notification.Message)

					metadata.ResourceData.Set("frequency", props.Schedule.Frequency)
					metadata.ResourceData.Set("days_of_week", props.Schedule.DaysOfWeek)
					metadata.ResourceData.Set("weeks_of_month", props.Schedule.WeeksOfMonth)
					metadata.ResourceData.Set("hour_of_day", props.Schedule.HourOfDay)
					metadata.ResourceData.Set("day_of_month", props.Schedule.DayOfMonth)
					metadata.ResourceData.Set("start_date", props.Schedule.StartDate)
					metadata.ResourceData.Set("end_date", props.Schedule.EndDate)
				}
			}

			return nil
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
					model.Properties.DisplayName = metadata.ResourceData.Get("display_name").(string)
				}

				if metadata.ResourceData.HasChange("view_id") {
					id, err := views.ParseScopedViewID(metadata.ResourceData.Get("view_id").(string))
					if err != nil {
						return err
					}
					model.Properties.ViewId = id.ID()
				}

				if metadata.ResourceData.HasChange("email_address_sender") {
					model.Properties.NotificationEmail = utils.String(metadata.ResourceData.Get("email_address_sender").(string))
				}

				if metadata.ResourceData.HasChange("email_subject") {
					model.Properties.Notification.Subject = metadata.ResourceData.Get("email_subject").(string)
				}

				if metadata.ResourceData.HasChange("email_addresses") {
					model.Properties.Notification.To = *utils.ExpandStringSlice(metadata.ResourceData.Get("email_addresses").([]interface{}))
				}

				if metadata.ResourceData.HasChange("message") {
					model.Properties.Notification.Message = utils.String(metadata.ResourceData.Get("message").(string))
				}

				if metadata.ResourceData.HasChange("frequency") {
					model.Properties.Schedule.Frequency = scheduledactions.ScheduleFrequency(metadata.ResourceData.Get("frequency").(string))
				}

				if metadata.ResourceData.HasChange("days_of_week") {
					var daysOfWeek []scheduledactions.DaysOfWeek
					if len(metadata.ResourceData.Get("days_of_week").([]interface{})) > 0 {
						daysOfWeek = make([]scheduledactions.DaysOfWeek, 0)
						for _, value := range metadata.ResourceData.Get("days_of_week").([]interface{}) {
							daysOfWeek = append(daysOfWeek, scheduledactions.DaysOfWeek(value.(string)))
						}
					}
					model.Properties.Schedule.DaysOfWeek = &daysOfWeek
				}

				if metadata.ResourceData.HasChange("weeks_of_month") {
					var weeksOfMonth []scheduledactions.WeeksOfMonth
					if len(metadata.ResourceData.Get("weeks_of_month").([]interface{})) > 0 {
						weeksOfMonth = make([]scheduledactions.WeeksOfMonth, 0)
						for _, value := range metadata.ResourceData.Get("weeks_of_month").([]interface{}) {
							weeksOfMonth = append(weeksOfMonth, scheduledactions.WeeksOfMonth(value.(string)))
						}
					}
					model.Properties.Schedule.WeeksOfMonth = &weeksOfMonth
				}

				if metadata.ResourceData.HasChange("start_date") {
					model.Properties.Schedule.StartDate = metadata.ResourceData.Get("start_date").(string)
				}

				if metadata.ResourceData.HasChange("end_date") {
					model.Properties.Schedule.EndDate = metadata.ResourceData.Get("end_date").(string)
				}

				if metadata.ResourceData.HasChange("hour_of_day") {
					model.Properties.Schedule.HourOfDay = utils.Int64(int64(metadata.ResourceData.Get("hour_of_day").(int)))
				}

				if metadata.ResourceData.HasChange("day_of_month") {
					model.Properties.Schedule.DayOfMonth = utils.Int64(int64(metadata.ResourceData.Get("day_of_month").(int)))
				}

				if _, err = client.CreateOrUpdateByScope(ctx, *id, *model, scheduledactions.CreateOrUpdateByScopeOperationOptions{}); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
