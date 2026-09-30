// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivateendpointconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

//go:generate go run ../../tools/generator-tests resourceidentity -properties name -compare-values subscription_id:bookshelf_id,resource_group_name:bookshelf_id,bookshelf_name:bookshelf_id

var (
	_ sdk.ResourceWithUpdate   = DiscoveryBookshelfPrivateEndpointConnectionResource{}
	_ sdk.ResourceWithIdentity = DiscoveryBookshelfPrivateEndpointConnectionResource{}
)

type DiscoveryBookshelfPrivateEndpointConnectionResource struct{}

type DiscoveryBookshelfPrivateEndpointConnectionModel struct {
	Name        string `tfschema:"name"`
	BookshelfId string `tfschema:"bookshelf_id"`
	Description string `tfschema:"description"`
	Status      string `tfschema:"status"`

	ActionsRequired   string   `tfschema:"actions_required"`
	GroupIds          []string `tfschema:"group_ids"`
	PrivateEndpointId string   `tfschema:"private_endpoint_id"`
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) ResourceType() string {
	return "azurerm_discovery_bookshelf_private_endpoint_connection"
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) ModelObject() interface{} {
	return &DiscoveryBookshelfPrivateEndpointConnectionModel{}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return bookshelfprivateendpointconnections.ValidateBookshelfPrivateEndpointConnectionID
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Identity() resourceids.ResourceId {
	return &bookshelfprivateendpointconnections.BookshelfPrivateEndpointConnectionId{}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"bookshelf_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: bookshelfprivateendpointconnections.ValidateBookshelfID,
		},

		"status": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(bookshelfprivateendpointconnections.PrivateEndpointServiceConnectionStatusApproved),
				string(bookshelfprivateendpointconnections.PrivateEndpointServiceConnectionStatusRejected),
			}, false),
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"actions_required": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"group_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_endpoint_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.BookshelfPrivateEndpointConnections

			var model DiscoveryBookshelfPrivateEndpointConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			bookshelfId, err := bookshelfprivateendpointconnections.ParseBookshelfID(model.BookshelfId)
			if err != nil {
				return err
			}

			id := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionID(bookshelfId.SubscriptionId, bookshelfId.ResourceGroupName, bookshelfId.BookshelfName, model.Name)

			// The Private Endpoint Connection is created by Azure when a Private Endpoint targets the
			// Bookshelf, so this resource adopts that connection in order to set its status rather than
			// creating it. As such there is no import check - the connection is always expected to exist.
			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found - a Private Endpoint targeting %s must exist before its connection can be managed", id, *bookshelfId)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			payload := *existing.Model
			payload.Properties = &bookshelfprivateendpointconnections.PrivateEndpointConnectionProperties{
				PrivateLinkServiceConnectionState: bookshelfprivateendpointconnections.PrivateLinkServiceConnectionState{
					Status: pointer.To(bookshelfprivateendpointconnections.PrivateEndpointServiceConnectionStatus(model.Status)),
				},
			}

			if model.Description != "" {
				payload.Properties.PrivateLinkServiceConnectionState.Description = pointer.To(model.Description)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.BookshelfPrivateEndpointConnections

			id, err := bookshelfprivateendpointconnections.ParseBookshelfPrivateEndpointConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			state := DiscoveryBookshelfPrivateEndpointConnectionModel{
				Name:        id.PrivateEndpointConnectionName,
				BookshelfId: bookshelfprivateendpointconnections.NewBookshelfID(id.SubscriptionId, id.ResourceGroupName, id.BookshelfName).ID(),
			}

			if props := resp.Model.Properties; props != nil {
				state.GroupIds = pointer.From(props.GroupIds)

				if privateEndpoint := props.PrivateEndpoint; privateEndpoint != nil {
					state.PrivateEndpointId = pointer.From(privateEndpoint.Id)
				}

				connectionState := props.PrivateLinkServiceConnectionState
				state.ActionsRequired = pointer.From(connectionState.ActionsRequired)
				state.Description = pointer.From(connectionState.Description)
				state.Status = string(pointer.From(connectionState.Status))
			}

			if err := metadata.Encode(&state); err != nil {
				return fmt.Errorf("encoding: %+v", err)
			}

			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, id)
		},
	}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.BookshelfPrivateEndpointConnections

			id, err := bookshelfprivateendpointconnections.ParseBookshelfPrivateEndpointConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model DiscoveryBookshelfPrivateEndpointConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			payload := *existing.Model
			if payload.Properties == nil {
				payload.Properties = &bookshelfprivateendpointconnections.PrivateEndpointConnectionProperties{}
			}

			if metadata.ResourceData.HasChange("status") {
				payload.Properties.PrivateLinkServiceConnectionState.Status = pointer.To(bookshelfprivateendpointconnections.PrivateEndpointServiceConnectionStatus(model.Status))
			}

			if metadata.ResourceData.HasChange("description") {
				payload.Properties.PrivateLinkServiceConnectionState.Description = pointer.To(model.Description)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r DiscoveryBookshelfPrivateEndpointConnectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.BookshelfPrivateEndpointConnections

			id, err := bookshelfprivateendpointconnections.ParseBookshelfPrivateEndpointConnectionID(metadata.ResourceData.Id())
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
