package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

type ManagerConnectivityConfigurationModel struct {
	Name                         string                       `tfschema:"name"`
	NetworkManagerId             string                       `tfschema:"network_manager_id"`
	AppliesToGroups              []ConnectivityGroupItemModel `tfschema:"applies_to_group"`
	ConnectivityTopology         network.ConnectivityTopology `tfschema:"connectivity_topology"`
	DeleteExistingPeeringEnabled bool                         `tfschema:"delete_existing_peering_enabled"`
	Description                  string                       `tfschema:"description"`
	Hub                          []HubModel                   `tfschema:"hub"`
	GlobalMeshEnabled            bool                         `tfschema:"global_mesh_enabled"`
}

type ConnectivityGroupItemModel struct {
	GroupConnectivity network.GroupConnectivity `tfschema:"group_connectivity"`
	GlobalMeshEnabled bool                      `tfschema:"global_mesh_enabled"`
	NetworkGroupId    string                    `tfschema:"network_group_id"`
	UseHubGateway     bool                      `tfschema:"use_hub_gateway"`
}

type HubModel struct {
	ResourceId   string `tfschema:"resource_id"`
	ResourceType string `tfschema:"resource_type"`
}

type ManagerConnectivityConfigurationResource struct{}

var _ sdk.ResourceWithUpdate = ManagerConnectivityConfigurationResource{}

func (r ManagerConnectivityConfigurationResource) ResourceType() string {
	return "azurerm_network_manager_connectivity_configuration"
}

func (r ManagerConnectivityConfigurationResource) ModelObject() interface{} {
	return &ManagerConnectivityConfigurationModel{}
}

func (r ManagerConnectivityConfigurationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.NetworkManagerConnectivityConfigurationID
}

func (r ManagerConnectivityConfigurationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"network_manager_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NetworkManagerID,
		},

		"applies_to_group": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"group_connectivity": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(network.GroupConnectivityNone),
							string(network.GroupConnectivityDirectlyConnected),
						}, false),
					},

					"global_mesh_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"network_group_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"use_hub_gateway": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},
				},
			},
		},

		"connectivity_topology": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(network.ConnectivityTopologyHubAndSpoke),
				string(network.ConnectivityTopologyMesh),
			}, false),
		},

		"delete_existing_peering_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"resource_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},

					"resource_type": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"global_mesh_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerConnectivityConfigurationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Network.ManagerConnectivityConfigurationsClient
			networkManagerId, err := parse.NetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := parse.NewNetworkManagerConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroup, networkManagerId.Name, model.Name)
			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conf := &network.ConnectivityConfiguration{
				ConnectivityConfigurationProperties: &network.ConnectivityConfigurationProperties{
					AppliesToGroups:       expandConnectivityGroupItemModel(model.AppliesToGroups),
					ConnectivityTopology:  model.ConnectivityTopology,
					DeleteExistingPeering: expandDeleteExistingPeering(model.DeleteExistingPeeringEnabled),
					IsGlobal:              expandConnectivityConfIsGlobal(model.GlobalMeshEnabled),
					Hubs:                  expandHubModel(model.Hub),
				},
			}

			if model.Description != "" {
				conf.ConnectivityConfigurationProperties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, *conf, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfigurationsClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ConnectivityConfigurationProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("applies_to_group") {
				properties.AppliesToGroups = expandConnectivityGroupItemModel(model.AppliesToGroups)
			}

			if metadata.ResourceData.HasChange("connectivity_topology") {
				properties.ConnectivityTopology = model.ConnectivityTopology
			}

			if metadata.ResourceData.HasChange("delete_existing_peering_enabled") {
				properties.DeleteExistingPeering = expandDeleteExistingPeering(model.DeleteExistingPeeringEnabled)
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Description = utils.String(model.Description)
			}

			if metadata.ResourceData.HasChange("hub") {
				properties.Hubs = expandHubModel(model.Hub)
			}

			if metadata.ResourceData.HasChange("global_mesh_enabled") {
				properties.IsGlobal = expandConnectivityConfIsGlobal(model.GlobalMeshEnabled)
			}

			if _, err := client.CreateOrUpdate(ctx, existing, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfigurationsClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName)
			if err != nil {
				if utils.ResponseWasNotFound(existing.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.ConnectivityConfigurationProperties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := ManagerConnectivityConfigurationModel{
				Name:                         id.ConnectivityConfigurationName,
				NetworkManagerId:             parse.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName).ID(),
				AppliesToGroups:              flattenConnectivityGroupItemModel(properties.AppliesToGroups),
				ConnectivityTopology:         properties.ConnectivityTopology,
				DeleteExistingPeeringEnabled: flattenDeleteExistingPeering(properties.DeleteExistingPeering),
				GlobalMeshEnabled:            flattenConnectivityConfIsGlobal(properties.IsGlobal),
				Hub:                          flattenHubModel(properties.Hubs),
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ManagerConnectivityConfigurationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ManagerConnectivityConfigurationsClient

			id, err := parse.NetworkManagerConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.NetworkManagerName, id.ConnectivityConfigurationName, utils.Bool(true))
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func expandDeleteExistingPeering(input bool) network.DeleteExistingPeering {
	if input {
		return network.DeleteExistingPeeringTrue
	}
	return network.DeleteExistingPeeringFalse
}

func expandConnectivityConfIsGlobal(input bool) network.IsGlobal {
	if input {
		return network.IsGlobalTrue
	}
	return network.IsGlobalFalse
}

func expandConnectivityGroupItemModel(inputList []ConnectivityGroupItemModel) *[]network.ConnectivityGroupItem {
	var outputList []network.ConnectivityGroupItem
	for _, v := range inputList {
		input := v
		output := network.ConnectivityGroupItem{
			GroupConnectivity: input.GroupConnectivity,
			IsGlobal:          expandConnectivityConfIsGlobal(input.GlobalMeshEnabled),
			NetworkGroupID:    utils.String(input.NetworkGroupId),
			UseHubGateway:     expandUseHubGateWay(input.UseHubGateway),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func expandUseHubGateWay(input bool) network.UseHubGateway {
	if input {
		return network.UseHubGatewayTrue
	}
	return network.UseHubGatewayFalse
}

func expandHubModel(inputList []HubModel) *[]network.Hub {
	var outputList []network.Hub
	for _, v := range inputList {
		input := v
		output := network.Hub{
			ResourceID:   utils.String(input.ResourceId),
			ResourceType: utils.String(input.ResourceType),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenDeleteExistingPeering(input network.DeleteExistingPeering) bool {
	return input == network.DeleteExistingPeeringTrue
}

func flattenConnectivityConfIsGlobal(input network.IsGlobal) bool {
	return input == network.IsGlobalTrue
}

func flattenConnectivityGroupItemModel(inputList *[]network.ConnectivityGroupItem) []ConnectivityGroupItemModel {
	var outputList []ConnectivityGroupItemModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := ConnectivityGroupItemModel{
			GroupConnectivity: input.GroupConnectivity,
			UseHubGateway:     flattenUseHubGateWay(input.UseHubGateway),
			GlobalMeshEnabled: flattenConnectivityConfIsGlobal(input.IsGlobal),
		}

		if input.NetworkGroupID != nil {
			output.NetworkGroupId = *input.NetworkGroupID
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenUseHubGateWay(input network.UseHubGateway) bool {
	return input == network.UseHubGatewayTrue
}

func flattenHubModel(inputList *[]network.Hub) []HubModel {
	var outputList []HubModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := HubModel{}

		if input.ResourceID != nil {
			output.ResourceId = *input.ResourceID
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return outputList
}
