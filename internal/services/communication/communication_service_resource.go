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
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = CommunicationServiceResource{}
var _ sdk.ResourceWithStateMigration = CommunicationServiceResource{}

type CommunicationServiceResource struct{}

func (CommunicationServiceResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.ServiceV0ToV1{},
		},
	}
}

func (CommunicationServiceResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CommunicationServiceName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"data_location": {
			Type: pluginsdk.TypeString,
			// TODO: should this become Required and remove the default in 4.0?
			Optional: true,
			ForceNew: true,
			Default:  "United States",
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

func (CommunicationServiceResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"primary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"secondary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (CommunicationServiceResource) ModelObject() interface{} {
	return nil
}

func (CommunicationServiceResource) ResourceType() string {
	return "azurerm_communication_service"
}

func (r CommunicationServiceResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.Communication.ServiceClient

			id := communicationservices.NewCommunicationServiceID(subscriptionId, metadata.ResourceData.Get("resource_group_name").(string), metadata.ResourceData.Get("name").(string))

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			param := communicationservices.CommunicationServiceResource{
				// The location is always `global` from the Azure Portal
				Location: location.Normalize("global"),
				Properties: &communicationservices.CommunicationServiceProperties{
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

func (r CommunicationServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			param := communicationservices.CommunicationServiceResource{
				// The location is always `global` from the Azure Portal
				Location: location.Normalize("global"),
				Properties: &communicationservices.CommunicationServiceProperties{
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

func (CommunicationServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
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

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			metadata.ResourceData.Set("name", id.CommunicationServiceName)
			metadata.ResourceData.Set("resource_group_name", id.ResourceGroupName)

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					metadata.ResourceData.Set("data_location", props.DataLocation)
				}

				if err := tags.FlattenAndSet(metadata.ResourceData, model.Tags); err != nil {
					return err
				}
			}

			if model := keysResp.Model; model != nil {
				metadata.ResourceData.Set("primary_connection_string", model.PrimaryConnectionString)
				metadata.ResourceData.Set("secondary_connection_string", model.SecondaryConnectionString)
				metadata.ResourceData.Set("primary_key", model.PrimaryKey)
				metadata.ResourceData.Set("secondary_key", model.SecondaryKey)
			}

			return nil
		},
	}
}

func (CommunicationServiceResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
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

func (CommunicationServiceResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return communicationservices.ValidateCommunicationServiceID
}
