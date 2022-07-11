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

			apiManagementId, err := parse.ApiManagementID(model.ApiManagementId)
			if err != nil {
				return err
			}

			id := parse.NewNotificationRecipientEmailID(subscriptionId, apiManagementId.ResourceGroup, apiManagementId.ServiceName, model.NotificationName, model.Email)

			// CheckEntityExists can not be used, it returns autorest error
			emails, err := client.ListByNotification(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName))
			if err != nil {
				if !utils.ResponseWasNotFound(emails.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if emails.Value != nil {
				for _, existing := range *emails.Value {
					if existing.RecipientEmailContractProperties != nil && existing.RecipientEmailContractProperties.Email != nil && *existing.RecipientEmailContractProperties.Email == model.Email {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}
				}
			}

			if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName), id.RecipientEmailName); err != nil {
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
			id, err := parse.NotificationRecipientEmailID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			// CheckEntityExists can not be used, it returns autorest error
			emails, err := client.ListByNotification(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName))
			if err != nil {
				if !utils.ResponseWasNotFound(emails.Response) {
					return fmt.Errorf("retrieving %s: %+v", id, err)
				}
			}

			found := false
			if emails.Value != nil {
				for _, existing := range *emails.Value {
					if existing.RecipientEmailContractProperties != nil && existing.RecipientEmailContractProperties.Email != nil && *existing.RecipientEmailContractProperties.Email == id.RecipientEmailName {
						found = true
					}
				}
			}
			if !found {
				return metadata.MarkAsGone(id)
			}

			model := ApiManagementNotificationRecipientEmailModel{
				ApiManagementId:  parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName).ID(),
				NotificationName: id.NotificationName,
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

			id, err := parse.NotificationRecipientEmailID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			_, err = client.Delete(ctx, id.ResourceGroup, id.ServiceName, apimanagement.NotificationName(id.NotificationName), id.RecipientEmailName)
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
