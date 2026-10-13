// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_notification_recipient_user

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientUserClient
			id, err := notificationrecipientuser.ParseRecipientUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// CheckEntityExists can not be used, it returns autorest error
			notificationId := notificationrecipientuser.NewNotificationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.NotificationName)
			users, err := client.ListByNotificationComplete(ctx, notificationId)
			if err != nil {
				if !response.WasNotFound(users.LatestHttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			found := false
			for _, existing := range users.Items {
				if existing.Name != nil && *existing.Name == id.UserId {
					found = true
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			model := ApiManagementNotificationRecipientUserModel{
				ApiManagementId:  apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
				NotificationName: string(id.NotificationName),
				UserId:           id.UserId,
			}

			return metadata.Encode(&model)
		},
	}
}
