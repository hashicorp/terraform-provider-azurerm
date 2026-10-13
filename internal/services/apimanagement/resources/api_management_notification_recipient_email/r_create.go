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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientEmailClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementNotificationRecipientEmailModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			apiManagementId, err := apimanagementservice.ParseServiceID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := notificationrecipientemail.NewRecipientEmailID(subscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, notificationrecipientemail.NotificationName(model.NotificationName), model.Email)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				// CheckEntityExists can not be used, it returns autorest error
				notificationId := notificationrecipientemail.NewNotificationID(subscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, notificationrecipientemail.NotificationName(model.NotificationName))
				emails, err := client.ListByNotificationComplete(ctx, notificationId)
				if err != nil {
					if !response.WasNotFound(emails.LatestHttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				for _, existing := range emails.Items {
					if existing.Properties != nil && existing.Properties.Email != nil && *existing.Properties.Email == model.Email {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
