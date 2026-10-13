// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_notification_recipient_email

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientemail"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientEmailClient
			id, err := notificationrecipientemail.ParseRecipientEmailID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// CheckEntityExists can not be used, it returns autorest error
			notificationId := notificationrecipientemail.NewNotificationID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.NotificationName)
			emails, err := client.ListByNotificationComplete(ctx, notificationId)
			if err != nil {
				if !response.WasNotFound(emails.LatestHttpResponse) {
					return fmt.Errorf("retrieving %s: %+v", notificationId, err)
				}
			}

			found := false
			for _, existing := range emails.Items {
				if existing.Properties != nil && existing.Properties.Email != nil && *existing.Properties.Email == id.RecipientEmailName {
					found = true
				}
			}

			if !found {
				return metadata.MarkAsGone(id)
			}

			model := ApiManagementNotificationRecipientEmailModel{
				ApiManagementId:  apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName).ID(),
				NotificationName: string(id.NotificationName),
				Email:            id.RecipientEmailName,
			}

			return metadata.Encode(&model)
		},
	}
}
