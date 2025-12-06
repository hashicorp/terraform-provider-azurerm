// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/senderusernames"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EmailCommunicationServiceDomainSenderUsernameResource{}

type EmailCommunicationServiceDomainSenderUsernameResource struct{}

type EmailCommunicationServiceDomainSenderUsernameModel struct {
	Name                 string `tfschema:"name"`
	EMailServiceDomainID string `tfschema:"email_service_domain_id"`
	DisplayName          string `tfschema:"display_name"`
}

func (EmailCommunicationServiceDomainSenderUsernameResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringLenBetween(1, 253),
		},

		"email_service_domain_id": commonschema.ResourceIDReferenceRequiredForceNew(&senderusernames.DomainId{}),

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 253),
		},
	}
}

func (EmailCommunicationServiceDomainSenderUsernameResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return senderusernames.ValidateSenderUsernameID
}

func (EmailCommunicationServiceDomainSenderUsernameResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (EmailCommunicationServiceDomainSenderUsernameResource) ModelObject() interface{} {
	return &EmailCommunicationServiceDomainSenderUsernameModel{}
}

func (EmailCommunicationServiceDomainSenderUsernameResource) ResourceType() string {
	return "azurerm_email_communication_service_domain_sender_username"
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.SenderUsernamesClient

			var model EmailCommunicationServiceDomainSenderUsernameModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			eMailServiceDomainID, err := senderusernames.ParseDomainID(model.EMailServiceDomainID)
			if err != nil {
				return err
			}

			id := senderusernames.NewSenderUsernameID(subscriptionId, eMailServiceDomainID.ResourceGroupName, eMailServiceDomainID.EmailServiceName, eMailServiceDomainID.DomainName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := senderusernames.SenderUsernameResource{
				Properties: &senderusernames.SenderUsernameProperties{
					// service API response indicates that the username must be same with the resource name
					Username: model.Name,
				},
			}

			if v := model.DisplayName; v != "" {
				parameters.Properties.DisplayName = pointer.To(v)
			}

			if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EmailCommunicationServiceDomainSenderUsernameResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.SenderUsernamesClient

			var model EmailCommunicationServiceDomainSenderUsernameModel

			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := senderusernames.ParseSenderUsernameID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			senderUsername := *existing.Model

			props := pointer.From(senderUsername.Properties)

			if metadata.ResourceData.HasChange("display_name") {
				props.DisplayName = pointer.To(model.DisplayName)
			}

			senderUsername.Properties = &props

			if _, err := client.CreateOrUpdate(ctx, *id, senderUsername); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (EmailCommunicationServiceDomainSenderUsernameResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.SenderUsernamesClient

			id, err := senderusernames.ParseSenderUsernameID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := EmailCommunicationServiceDomainSenderUsernameModel{
				Name:                 id.SenderUsernameName,
				EMailServiceDomainID: senderusernames.NewDomainID(id.SubscriptionId, id.ResourceGroupName, id.EmailServiceName, id.DomainName).ID(),
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DisplayName = pointer.From(props.DisplayName)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (EmailCommunicationServiceDomainSenderUsernameResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.SenderUsernamesClient

			id, err := senderusernames.ParseSenderUsernameID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
