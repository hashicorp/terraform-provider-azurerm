// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = EmailCommunicationServiceResource{}

type EmailCommunicationServiceResource struct{}

type EmailCommunicationServiceResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	DataLocation      string            `tfschema:"data_location"`
	Tags              map[string]string `tfschema:"tags"`
}

func (EmailCommunicationServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CommunicationServiceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"data_location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Africa",
				"Asia Pacific",
				"Australia",
				"Brazil",
				"Canada",
				"Europe",
				"France",
				"Germany",
				"India",
				"Japan",
				"Korea",
				"Norway",
				"Switzerland",
				"UAE",
				"UK",
				"United States",
				"usgov",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (EmailCommunicationServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (EmailCommunicationServiceResource) ModelObject() interface{} {
	return &EmailCommunicationServiceResourceModel{}
}

func (EmailCommunicationServiceResource) ResourceType() string {
	return "azurerm_email_communication_service"
}

func (r EmailCommunicationServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.EmailServicesClient

			var model EmailCommunicationServiceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := emailservices.NewEmailServiceID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := emailservices.EmailServiceResource{
				// The location is always `global` from the Azure Portal
				Location: location.Normalize("global"),
				Properties: &emailservices.EmailServiceProperties{
					DataLocation: model.DataLocation,
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r EmailCommunicationServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.EmailServicesClient

			var model EmailCommunicationServiceResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := emailservices.ParseEmailServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			emailService := *existing.Model

			props := pointer.From(emailService.Properties)

			if metadata.ResourceData.HasChange("data_location") {
				props.DataLocation = model.DataLocation
			}

			if metadata.ResourceData.HasChange("tags") {
				emailService.Tags = pointer.To(model.Tags)
			}

			emailService.Properties = &props

			if err := client.CreateOrUpdateThenPoll(ctx, *id, emailService); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (EmailCommunicationServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.EmailServicesClient

			state := EmailCommunicationServiceResourceModel{}

			id, err := emailservices.ParseEmailServiceID(metadata.ResourceData.Id())
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

			state.Name = id.EmailServiceName
			state.ResourceGroupName = id.ResourceGroupName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DataLocation = props.DataLocation
				}

				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (EmailCommunicationServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.EmailServicesClient

			id, err := emailservices.ParseEmailServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (EmailCommunicationServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return emailservices.ValidateEmailServiceID
}
