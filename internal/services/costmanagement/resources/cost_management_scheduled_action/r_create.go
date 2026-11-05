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
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			var config CostManagementScheduledActionModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			viewId, err := viewoperationgroup.ParseScopedViewID(config.ViewId)
			if err != nil {
				return err
			}
			id := scheduledactionoperationgroup.NewScopedScheduledActionID(viewId.Scope, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ScheduledActionsGetByScope(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return tf.ImportAsExistsError(r.ResourceType(), id.ID())
				}
			}

			var daysOfWeek []scheduledactionoperationgroup.DaysOfWeek
			if len(config.DaysOfWeek) > 0 {
				daysOfWeek = make([]scheduledactionoperationgroup.DaysOfWeek, 0)
				for _, value := range config.DaysOfWeek {
					daysOfWeek = append(daysOfWeek, scheduledactionoperationgroup.DaysOfWeek(value))
				}
			}

			var weeksOfMonth []scheduledactionoperationgroup.WeeksOfMonth
			if len(config.WeeksOfMonth) > 0 {
				weeksOfMonth = make([]scheduledactionoperationgroup.WeeksOfMonth, 0)
				for _, value := range config.WeeksOfMonth {
					weeksOfMonth = append(weeksOfMonth, scheduledactionoperationgroup.WeeksOfMonth(value))
				}
			}

			schedule := scheduledactionoperationgroup.ScheduleProperties{
				Frequency:    scheduledactionoperationgroup.ScheduleFrequency(config.Frequency),
				WeeksOfMonth: &weeksOfMonth,
				DaysOfWeek:   &daysOfWeek,
				HourOfDay:    pointer.To(config.HourOfDay),
				DayOfMonth:   pointer.To(config.DayOfMonth),
				StartDate:    config.StartDate,
				EndDate:      config.EndDate,
			}

			props := scheduledactionoperationgroup.ScheduledAction{
				Kind: pointer.To(scheduledactionoperationgroup.ScheduledActionKindEmail),
				Properties: &scheduledactionoperationgroup.ScheduledActionProperties{
					DisplayName: config.DisplayName,
					Status:      scheduledactionoperationgroup.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactionoperationgroup.FileDestination{
						FileFormats: &[]scheduledactionoperationgroup.FileFormat{},
					},
					NotificationEmail: pointer.To(config.EmailAddressSender),
					Notification: scheduledactionoperationgroup.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}

			if _, err = client.ScheduledActionsCreateOrUpdateByScope(ctx, id, props, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
