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
	"github.com/hashicorp/go-azure-sdk/resource-manager/communication/2023-03-31/communicationservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = CommunicationServiceResource{}
var _ sdk.ResourceWithStateMigration = CommunicationServiceResource{}

type CommunicationServiceResource struct{}

type CommunicationServiceResourceModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	DataLocation      string            `tfschema:"data_location"`
	Tags              map[string]string `tfschema:"tags"`

	PrimaryConnectionString   string `tfschema:"primary_connection_string"`
	SecondaryConnectionString string `tfschema:"secondary_connection_string"`
	PrimaryKey                string `tfschema:"primary_key"`
	SecondaryKey              string `tfschema:"secondary_key"`
}

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
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"primary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},
	}
}

func (CommunicationServiceResource) ModelObject() interface{} {
	return &CommunicationServiceResourceModel{}
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

			var model CommunicationServiceResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id := communicationservices.NewCommunicationServiceID(subscriptionId, model.ResourceGroupName, model.Name)

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

func (r CommunicationServiceResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			var model CommunicationServiceResourceModel

			if err := metadata.Decode(&model); err != nil {
				return err
			}

			id, err := communicationservices.ParseCommunicationServiceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil || existing.Model == nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			commService := *existing.Model

			props := pointer.From(commService.Properties)

			if metadata.ResourceData.HasChange("data_location") {
				props.DataLocation = model.DataLocation
			}

			existing.Model.Properties = &props

			if metadata.ResourceData.HasChange("tags") {
				commService.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, commService); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			return nil
		},
	}
}

func (CommunicationServiceResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Communication.ServiceClient

			state := CommunicationServiceResourceModel{}

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

			state.Name = id.CommunicationServiceName
			state.ResourceGroupName = id.ResourceGroupName

			keysResp, err := client.ListKeys(ctx, *id)
			if err != nil {
				return fmt.Errorf("listing keys for %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					state.DataLocation = props.DataLocation
				}

				state.Tags = pointer.From(model.Tags)
			}

			if model := keysResp.Model; model != nil {
				state.PrimaryConnectionString = pointer.From(model.PrimaryConnectionString)
				state.SecondaryConnectionString = pointer.From(model.SecondaryConnectionString)
				state.PrimaryKey = pointer.From(model.PrimaryKey)
				state.SecondaryKey = pointer.From(model.SecondaryKey)
			}

			return metadata.Encode(&state)
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
