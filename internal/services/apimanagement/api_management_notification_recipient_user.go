package apimanagement

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApiManagementNotificationRecipientUserModel struct {
	ApiManagementId  string `tfschema:"api_management_id"`
	NotificationName string `tfschema:"notification_type"`
	UserId           string `tfschema:"user_id"`
}

type ApiManagementNotificationRecipientUserResource struct{}

var _ sdk.Resource = ApiManagementNotificationRecipientUserResource{}

func (r ApiManagementNotificationRecipientUserResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"api_management_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ApiManagementID,
		},

		"notification_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(apimanagement.NotificationNameAccountClosedPublisher),
				string(apimanagement.NotificationNameBCC),
				string(apimanagement.NotificationNameNewApplicationNotificationMessage),
				string(apimanagement.NotificationNameNewIssuePublisherNotificationMessage),
				string(apimanagement.NotificationNamePurchasePublisherNotificationMessage),
				string(apimanagement.NotificationNameQuotaLimitApproachingPublisherNotificationMessage),
				string(apimanagement.NotificationNameRequestPublisherNotificationMessage),
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

			apiManagementId, err := parse.ApiManagementID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := parse.NewNotificationRecipientUserID(subscriptionId, apiManagementId.ResourceGroup, apiManagementId.ServiceName, model.NotificationName, model.UserId)

			// CheckEntityExists can not be used, it returns autorest error
			users, err := client.ListByNotification(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName))
			if err != nil {
				if !utils.ResponseWasNotFound(users.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if users.Value != nil {
				for _, existing := range *users.Value {
					if existing.Name != nil && *existing.Name == model.UserId {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName), id.RecipientUserName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ApiManagement.NotificationRecipientUserClient
			id, err := parse.NotificationRecipientUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// CheckEntityExists can not be used, it returns autorest error
			users, err := client.ListByNotification(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName))
			if err != nil {
				if !utils.ResponseWasNotFound(users.Response) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			found := false
			if users.Value != nil {
				for _, existing := range *users.Value {
					if existing.Name != nil && *existing.Name == id.RecipientUserName {
						found = true
					}
				}
			}
			if !found {
				return metadata.MarkAsGone(id)
			}

			model := ApiManagementNotificationRecipientUserModel{
				ApiManagementId:  parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName).ID(),
				NotificationName: id.NotificationName,
				UserId:           id.RecipientUserName,
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

			id, err := parse.NotificationRecipientUserID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.Delete(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName), id.RecipientUserName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (r ApiManagementNotificationRecipientUserResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NotificationRecipientUserID
}
