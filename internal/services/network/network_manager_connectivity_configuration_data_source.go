package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type ManagerConnectivityConfigurationDataSource struct{}

var _ sdk.DataSource = ManagerConnectivityConfigurationDataSource{}

func (r ManagerConnectivityConfigurationDataSource) ResourceType() string {
	return "azurerm_network_manager_connectivity_configuration"
}

func (r ManagerConnectivityConfigurationDataSource) ModelObject() interface{} {
	return &ManagerConnectivityConfigurationModel{}
}

func (r ManagerConnectivityConfigurationDataSource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: connectivityconfigurations.ValidateNetworkManagerID,
		},
	}
}

func (r ManagerConnectivityConfigurationDataSource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"applies_to_group": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_connectivity": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"global_mesh_enabled": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},

					"network_group_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"use_hub_gateway": {
						Type:     pluginsdk.TypeBool,
						Computed: true,
					},
				},
			},
		},

		"connectivity_topology": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"delete_existing_peering_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"hub": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"resource_type": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"global_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}

func (r ManagerConnectivityConfigurationDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ConnectivityConfigurations

			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			networkManagerId, err := connectivityconfigurations.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := connectivityconfigurations.NewConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			resp, err := client.Get(ctx, id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("%s does not exist", id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			state := ManagerConnectivityConfigurationModel{
				Name:             model.Name,
				NetworkManagerId: networkManagerId.ID(),
			}

			if model := resp.Model; model != nil {
				if properties := model.Properties; properties != nil {
					state.AppliesToGroups = flattenConnectivityGroupItemModel(properties.AppliesToGroups)
					state.ConnectivityTopology = properties.ConnectivityTopology
					state.DeleteExistingPeeringEnabled = flattenDeleteExistingPeering(properties.DeleteExistingPeering)
					state.GlobalMeshEnabled = flattenConnectivityConfIsGlobal(properties.IsGlobal)
					state.Hub = flattenHubModel(properties.Hubs)
					state.Description = pointer.From(properties.Description)
				}
			}

			metadata.SetID(id)

			return metadata.Encode(&state)
		},
	}
}
