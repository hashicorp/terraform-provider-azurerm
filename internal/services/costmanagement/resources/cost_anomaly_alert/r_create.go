// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cost_anomaly_alert

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/scheduledactionoperationgroup"
	"github.com/hashicorp/go-azure-sdk/resource-manager/costmanagement/2025-03-01/viewoperationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.CostManagement.ScheduledActionOperationGroupClient

			var config AnomalyAlertModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := config.SubscriptionId
			if subscriptionId == "" {
				subscriptionId = commonids.NewSubscriptionID(metadata.Client.Account.SubscriptionId).ID()
			}
			id := scheduledactionoperationgroup.NewScopedScheduledActionID(subscriptionId, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ScheduledActionsGetByScope(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			viewId := viewoperationgroup.NewScopedViewID(subscriptionId, "ms:DailyAnomalyByResourceGroup")

			schedule := scheduledactionoperationgroup.ScheduleProperties{
				Frequency: scheduledactionoperationgroup.ScheduleFrequencyDaily,
			}
			schedule.SetEndDateAsTime(time.Now().AddDate(1, 0, 0))
			schedule.SetStartDateAsTime(time.Now())

			notificationEmail := config.EmailAddresses[0]
			if config.NotificationEmail != "" {
				notificationEmail = config.NotificationEmail
			}

			param := scheduledactionoperationgroup.ScheduledAction{
				Kind: pointer.To(scheduledactionoperationgroup.ScheduledActionKindInsightAlert),
				Properties: &scheduledactionoperationgroup.ScheduledActionProperties{
					DisplayName: config.DisplayName,
					Status:      scheduledactionoperationgroup.ScheduledActionStatusEnabled,
					ViewId:      viewId.ID(),
					FileDestination: &scheduledactionoperationgroup.FileDestination{
						FileFormats: &[]scheduledactionoperationgroup.FileFormat{},
					},
					NotificationEmail: &notificationEmail,
					Notification: scheduledactionoperationgroup.NotificationProperties{
						Subject: config.EmailSubject,
						Message: pointer.To(config.Message),
						To:      config.EmailAddresses,
					},
					Schedule: schedule,
				},
			}
			if _, err := client.ScheduledActionsCreateOrUpdateByScope(ctx, id, param, scheduledactionoperationgroup.DefaultScheduledActionsCreateOrUpdateByScopeOperationOptions()); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
