// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package discovery

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/discovery/2026-06-01/bookshelfprivatelinkresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.DataSource = DiscoveryBookshelfPrivateLinkResourcesDataSource{}

type DiscoveryBookshelfPrivateLinkResourcesDataSource struct{}

type DiscoveryBookshelfPrivateLinkResourcesModel struct {
	BookshelfId          string                                       `tfschema:"bookshelf_id"`
	PrivateLinkResources []DiscoveryBookshelfPrivateLinkResourceModel `tfschema:"private_link_resources"`
}

type DiscoveryBookshelfPrivateLinkResourceModel struct {
	Name              string   `tfschema:"name"`
	GroupId           string   `tfschema:"group_id"`
	RequiredMembers   []string `tfschema:"required_members"`
	RequiredZoneNames []string `tfschema:"required_zone_names"`
}

func (r DiscoveryBookshelfPrivateLinkResourcesDataSource) ResourceType() string {
	return "azurerm_discovery_bookshelf_private_link_resources"
}

func (r DiscoveryBookshelfPrivateLinkResourcesDataSource) ModelObject() any {
	return &DiscoveryBookshelfPrivateLinkResourcesModel{}
}

func (r DiscoveryBookshelfPrivateLinkResourcesDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"bookshelf_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: bookshelfprivatelinkresources.ValidateBookshelfID,
		},
	}
}

func (r DiscoveryBookshelfPrivateLinkResourcesDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"private_link_resources": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"required_members": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"required_zone_names": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},
	}
}

func (r DiscoveryBookshelfPrivateLinkResourcesDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Discovery.BookshelfPrivateLinkResources

			var model DiscoveryBookshelfPrivateLinkResourcesModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := bookshelfprivatelinkresources.ParseBookshelfID(model.BookshelfId)
			if err != nil {
				return err
			}

			resp, err := client.ListByBookshelfComplete(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving Private Link Resources for %s: %+v", *id, err)
			}

			state := DiscoveryBookshelfPrivateLinkResourcesModel{
				BookshelfId:          id.ID(),
				PrivateLinkResources: make([]DiscoveryBookshelfPrivateLinkResourceModel, 0),
			}

			for _, item := range resp.Items {
				privateLinkResource := DiscoveryBookshelfPrivateLinkResourceModel{
					Name:              pointer.From(item.Name),
					RequiredMembers:   make([]string, 0),
					RequiredZoneNames: make([]string, 0),
				}

				if props := item.Properties; props != nil {
					privateLinkResource.GroupId = pointer.From(props.GroupId)
					privateLinkResource.RequiredMembers = pointer.From(props.RequiredMembers)
					privateLinkResource.RequiredZoneNames = pointer.From(props.RequiredZoneNames)
				}

				state.PrivateLinkResources = append(state.PrivateLinkResources, privateLinkResource)
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
