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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/privateendpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.ResourceWithUpdate   = DiscoveryBookshelfPrivateEndpointConnectionResource{}
	_ sdk.ResourceWithIdentity = DiscoveryBookshelfPrivateEndpointConnectionResource{}
)

type DiscoveryBookshelfPrivateEndpointConnectionResource struct{}

type DiscoveryBookshelfPrivateEndpointConnectionModel struct {
	BookshelfId       string `tfschema:"bookshelf_id"`
	PrivateEndpointId string `tfschema:"private_endpoint_id"`
	Description       string `tfschema:"description"`
	Status            string `tfschema:"status"`

	ActionsRequired string   `tfschema:"actions_required"`
	GroupIds        []string `tfschema:"group_ids"`
	Name            string   `tfschema:"name"`
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
		"bookshelf_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: bookshelfprivateendpointconnections.ValidateBookshelfID,
		},

		"private_endpoint_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: privateendpoints.ValidatePrivateEndpointID,
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
			Computed:     true,
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

		"name": {
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

			privateEndpointId, err := privateendpoints.ParsePrivateEndpointID(model.PrivateEndpointId)
			if err != nil {
				return err
			}

			connections, err := client.ListByBookshelfComplete(ctx, *bookshelfId)
			if err != nil {
				return fmt.Errorf("listing Private Endpoint Connections for %s: %+v", *bookshelfId, err)
			}

			var existing *bookshelfprivateendpointconnections.BookshelfPrivateEndpointConnection
			for _, item := range connections.Items {
				if item.Properties == nil || item.Properties.PrivateEndpoint == nil || item.Properties.PrivateEndpoint.Id == nil {
					continue
				}

				if candidate, err := privateendpoints.ParsePrivateEndpointIDInsensitively(*item.Properties.PrivateEndpoint.Id); err == nil && resourceids.Match(candidate, privateEndpointId) {
					existing = &item
					break
				}
			}

			if existing == nil || existing.Name == nil {
				return fmt.Errorf("no Private Endpoint Connection for %s was found on %s - the Private Endpoint must target the Bookshelf before its connection can be managed, and must be recreated if its connection was previously deleted", *privateEndpointId, *bookshelfId)
			}

			id := bookshelfprivateendpointconnections.NewBookshelfPrivateEndpointConnectionID(bookshelfId.SubscriptionId, bookshelfId.ResourceGroupName, bookshelfId.BookshelfName, *existing.Name)

			payload := *existing
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

				if privateEndpoint := props.PrivateEndpoint; privateEndpoint != nil && privateEndpoint.Id != nil {
					privateEndpointId, err := privateendpoints.ParsePrivateEndpointIDInsensitively(*privateEndpoint.Id)
					if err != nil {
						return err
					}
					state.PrivateEndpointId = privateEndpointId.ID()
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
