// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_management_scheduled_action

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			id, err := scheduledactionoperationgroup.ParseScopedScheduledActionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ScheduledActionsGetByScope(ctx, *id)
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

					viewId, err := viewoperationgroup.ParseScopedViewID(props.ViewId)
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
