// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ApiManagementNotificationRecipientEmailModel struct {
	ApiManagementId  string `tfschema:"api_management_id"`
	NotificationName string `tfschema:"notification_type"`
	Email            string `tfschema:"email"`
}

type ApiManagementNotificationRecipientEmailResource struct{}

var _ sdk.Resource = ApiManagementNotificationRecipientEmailResource{}

func (r ApiManagementNotificationRecipientEmailResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: apimanagementservice.ValidateServiceID,
		},

		"notification_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(notificationrecipientemail.NotificationNameAccountClosedPublisher),
				string(notificationrecipientemail.NotificationNameBCC),
				string(notificationrecipientemail.NotificationNameNewApplicationNotificationMessage),
				string(notificationrecipientemail.NotificationNameNewIssuePublisherNotificationMessage),
				string(notificationrecipientemail.NotificationNamePurchasePublisherNotificationMessage),
				string(notificationrecipientemail.NotificationNameQuotaLimitApproachingPublisherNotificationMessage),
				string(notificationrecipientemail.NotificationNameRequestPublisherNotificationMessage),
			}, false),
		},

		"email": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementNotificationRecipientEmailResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementNotificationRecipientEmailResource) ModelObject() interface{} {
	return &ApiManagementNotificationRecipientEmailModel{}
}

func (r ApiManagementNotificationRecipientEmailResource) ResourceType() string {
	return "azurerm_api_management_notification_recipient_email"
}

func (r ApiManagementNotificationRecipientEmailResource) Create() sdk.ResourceFunc {
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

			if _, err = client.CreateOrUpdate(ctx, id); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientEmailResource) Read() sdk.ResourceFunc {
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

func (r ApiManagementNotificationRecipientEmailResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientEmailClient

			id, err := notificationrecipientemail.ParseRecipientEmailID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientEmailResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NotificationRecipientEmailID
}
