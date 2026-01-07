// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/notificationrecipientuser"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name api_management_notification_recipient_user -service-package-name apimanagement -properties "notification_name:notification_type,user_id" -compare-values "resource_group_name:api_management_id,service_name:api_management_id" -known-values "subscription_id:data.Subscriptions.Primary" -test-name basic

type ApiManagementNotificationRecipientUserModel struct {
	ApiManagementId  string `tfschema:"api_management_id"`
	NotificationName string `tfschema:"notification_type"`
	UserId           string `tfschema:"user_id"`
}

type ApiManagementNotificationRecipientUserResource struct{}

var (
	_ sdk.Resource             = ApiManagementNotificationRecipientUserResource{}
	_ sdk.ResourceWithIdentity = ApiManagementNotificationRecipientUserResource{}
)

func (r ApiManagementNotificationRecipientUserResource) Identity() resourceids.ResourceId {
	return &notificationrecipientuser.RecipientUserId{}
}

func (r ApiManagementNotificationRecipientUserResource) Arguments() map[string]*pluginsdk.Schema {
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
				string(notificationrecipientuser.NotificationNameAccountClosedPublisher),
				string(notificationrecipientuser.NotificationNameBCC),
				string(notificationrecipientuser.NotificationNameNewApplicationNotificationMessage),
				string(notificationrecipientuser.NotificationNameNewIssuePublisherNotificationMessage),
				string(notificationrecipientuser.NotificationNamePurchasePublisherNotificationMessage),
				string(notificationrecipientuser.NotificationNameQuotaLimitApproachingPublisherNotificationMessage),
				string(notificationrecipientuser.NotificationNameRequestPublisherNotificationMessage),
			}, false),
		},

		"user_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApiManagementNotificationRecipientUserResource) ModelObject() interface{} {
	return &ApiManagementNotificationRecipientUserModel{}
}

func (r ApiManagementNotificationRecipientUserResource) ResourceType() string {
	return "azurerm_api_management_notification_recipient_user"
}

func (r ApiManagementNotificationRecipientUserResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientUserClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ApiManagementNotificationRecipientUserModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			apiManagementId, err := apimanagementservice.ParseServiceID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := notificationrecipientuser.NewRecipientUserID(subscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, notificationrecipientuser.NotificationName(model.NotificationName), model.UserId)

			// CheckEntityExists can not be used, it returns autorest error
			notificationId := notificationrecipientuser.NewNotificationID(subscriptionId, apiManagementId.ResourceGroupName, apiManagementId.ServiceName, notificationrecipientuser.NotificationName(model.NotificationName))
			users, err := client.ListByNotificationComplete(ctx, notificationId)
			if err != nil {
				if !response.WasNotFound(users.LatestHttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			for _, existing := range users.Items {
				if existing.Name != nil && *existing.Name == model.UserId {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) Read() sdk.ResourceFunc {
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

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&model)
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientUserClient

			id, err := notificationrecipientuser.ParseRecipientUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NotificationRecipientUserID
}
