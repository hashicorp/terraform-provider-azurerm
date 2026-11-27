// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type DataSource struct{}

var _ sdk.DataSource = DataSource{}

func (r DataSource) ResourceType() string {
	return "azurerm_network_manager"
}

func (r DataSource) Description() string {
	return "Use this data source to access information about a Network Manager."
}

func (r DataSource) ModelObject() any {
	return &ManagerModel{}
}

func (r DataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Network Manager.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupNameForDataSource(),
	}
}

func (r DataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.LocationComputed(),

		"scope": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"management_group_ids": {
						Description: "A list of Management Group IDs used as the scope for the Network Manager.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"subscription_ids": {
						Description: "A list of subscription IDs used as the scope for the Network Manager.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"scope_accesses": {
			Description: "A list of configuration deployment type configured on the Network Manager.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"cross_tenant_scopes": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tenant_id": {
						Description: "The tenant ID of the cross tenant scope.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
					"subscriptions": {
						Description: "A list of subscriptions used as cross tenant scope for the Network Manager.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"management_groups": {
						Description: "A list of Management Groups used as cross tenant scope for the Network Manager.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"description": {
			Description: "A description of the Network Manager.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.TagsDataSource(),
	}
}

func (r DataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagers
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ManagerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := networkmanagers.NewNetworkManagerID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", id)
			}

			properties := existing.Model.Properties
			var description string
			var scope []ManagerScopeModel
			var ScopeAccesses []string
			if properties.Description != nil {
				description = *properties.Description
			}
			scope = flattenNetworkManagerScope(properties.NetworkManagerScopes)
			ScopeAccesses = flattenNetworkManagerScopeAccesses(properties.NetworkManagerScopeAccesses)

			state := ManagerModel{
				CrossTenantScopes: flattenNetworkManagerCrossTenantScopes(properties.NetworkManagerScopes.CrossTenantScopes),
				Description:       description,
				Location:          location.NormalizeNilable(existing.Model.Location),
				Name:              id.NetworkManagerName,
				ResourceGroupName: id.ResourceGroupName,
				ScopeAccesses:     ScopeAccesses,
				Scope:             scope,
				Tags:              pluginsdk.FlattenPtrMapStringString(existing.Model.Tags),
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
