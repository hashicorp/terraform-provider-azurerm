// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package communication

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/emailservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = EmailCommunicationServiceResource{}

type EmailCommunicationServiceResource struct{}

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
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}

func (EmailCommunicationServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (EmailCommunicationServiceResource) ModelObject() interface{} {
	return nil
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

			id := emailservices.NewEmailServiceID(subscriptionId, metadata.ResourceData.Get("resource_group_name").(string), metadata.ResourceData.Get("name").(string))

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
					DataLocation: metadata.ResourceData.Get("data_location").(string),
				},
				Tags: tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
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

			id, err := emailservices.ParseEmailServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			param := emailservices.EmailServiceResource{
				// The location is always `global` from the Azure Portal
				Location: location.Normalize("global"),
				Properties: &emailservices.EmailServiceProperties{
					DataLocation: metadata.ResourceData.Get("data_location").(string),
				},
				Tags: tags.Expand(metadata.ResourceData.Get("tags").(map[string]interface{})),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, param); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (EmailCommunicationServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.EmailServicesClient

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
			metadata.ResourceData.Set("name", id.EmailServiceName)
			metadata.ResourceData.Set("resource_group_name", id.ResourceGroupName)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("data_location", props.DataLocation)
				}

				if err := tags.FlattenAndSet(metadata.ResourceData, model.Tags); err != nil {
					return err
				}
			}

			return nil
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
