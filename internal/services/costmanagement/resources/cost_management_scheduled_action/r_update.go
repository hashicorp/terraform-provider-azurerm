// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_management_scheduled_action

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			id, err := scheduledactionoperationgroup.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config CostManagementScheduledActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// Update operation requires latest eTag to be set in the request.
			existing, err := client.ScheduledActionsGetByScope(ctx, *id)
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
					viewId, err := viewoperationgroup.ParseScopedViewID(config.ViewId)
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
					model.Properties.Schedule.Frequency = scheduledactionoperationgroup.ScheduleFrequency(config.Frequency)
				}

				if metadata.ResourceData.HasChange("days_of_week") {
					var daysOfWeek []scheduledactionoperationgroup.DaysOfWeek
					if len(config.DaysOfWeek) > 0 {
						daysOfWeek = make([]scheduledactionoperationgroup.DaysOfWeek, 0)
						for _, value := range config.DaysOfWeek {
							daysOfWeek = append(daysOfWeek, scheduledactionoperationgroup.DaysOfWeek(value))
						}
					}
					model.Properties.Schedule.DaysOfWeek = &daysOfWeek
				}

				if metadata.ResourceData.HasChange("weeks_of_month") {
					var weeksOfMonth []scheduledactionoperationgroup.WeeksOfMonth
					if len(config.WeeksOfMonth) > 0 {
						weeksOfMonth = make([]scheduledactionoperationgroup.WeeksOfMonth, 0)
						for _, value := range config.WeeksOfMonth {
							weeksOfMonth = append(weeksOfMonth, scheduledactionoperationgroup.WeeksOfMonth(value))
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

				if _, err = client.ScheduledActionsCreateOrUpdateByScope(ctx, *id, *model, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions()); err != nil {
					return fmt.Errorf("updating %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
