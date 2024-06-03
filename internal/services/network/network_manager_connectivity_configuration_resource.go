// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/connectivityconfigurations"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ManagerConnectivityConfigurationModel struct {
	Name                         string                                          `tfschema:"name"`
	NetworkManagerId             string                                          `tfschema:"network_manager_id"`
	AppliesToGroups              []ConnectivityGroupItemModel                    `tfschema:"applies_to_group"`
	ConnectivityTopology         connectivityconfigurations.ConnectivityTopology `tfschema:"connectivity_topology"`
	DeleteExistingPeeringEnabled bool                                            `tfschema:"delete_existing_peering_enabled"`
	Description                  string                                          `tfschema:"description"`
	Hub                          []HubModel                                      `tfschema:"hub"`
	GlobalMeshEnabled            bool                                            `tfschema:"global_mesh_enabled"`
}

type ConnectivityGroupItemModel struct {
	GroupConnectivity connectivityconfigurations.GroupConnectivity `tfschema:"group_connectivity"`
	GlobalMeshEnabled bool                                         `tfschema:"global_mesh_enabled"`
	NetworkGroupId    string                                       `tfschema:"network_group_id"`
	UseHubGateway     bool                                         `tfschema:"use_hub_gateway"`
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
	return connectivityconfigurations.ValidateConnectivityConfigurationID
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
			ValidateFunc: connectivityconfigurations.ValidateNetworkManagerID,
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
							string(connectivityconfigurations.GroupConnectivityNone),
							string(connectivityconfigurations.GroupConnectivityDirectlyConnected),
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
				string(connectivityconfigurations.ConnectivityTopologyHubAndSpoke),
				string(connectivityconfigurations.ConnectivityTopologyMesh),
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

			client := metadata.Client.Network.ConnectivityConfigurations
			networkManagerId, err := connectivityconfigurations.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := connectivityconfigurations.NewConnectivityConfigurationID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			conf := connectivityconfigurations.ConnectivityConfiguration{
				Properties: &connectivityconfigurations.ConnectivityConfigurationProperties{
					AppliesToGroups:       expandConnectivityGroupItemModel(model.AppliesToGroups),
					ConnectivityTopology:  model.ConnectivityTopology,
					DeleteExistingPeering: expandDeleteExistingPeering(model.DeleteExistingPeeringEnabled),
					IsGlobal:              expandConnectivityConfIsGlobal(model.GlobalMeshEnabled),
					Hubs:                  expandHubModel(model.Hub),
				},
			}

			if model.Description != "" {
				conf.Properties.Description = &model.Description
			}

			if _, err := client.CreateOrUpdate(ctx, id, conf); err != nil {
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
			client := metadata.Client.Network.ConnectivityConfigurations

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ManagerConnectivityConfigurationModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

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

			if _, err := client.CreateOrUpdate(ctx, *id, *existing.Model); err != nil {
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
			client := metadata.Client.Network.ConnectivityConfigurations

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerConnectivityConfigurationModel{
				Name:                         id.ConnectivityConfigurationName,
				NetworkManagerId:             connectivityconfigurations.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
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
			client := metadata.Client.Network.ConnectivityConfigurations

			id, err := connectivityconfigurations.ParseConnectivityConfigurationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			err = client.DeleteThenPoll(ctx, *id, connectivityconfigurations.DeleteOperationOptions{
				Force: utils.Bool(true),
			})
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandDeleteExistingPeering(input bool) *connectivityconfigurations.DeleteExistingPeering {
	output := connectivityconfigurations.DeleteExistingPeeringFalse
	if input {
		output = connectivityconfigurations.DeleteExistingPeeringTrue
	}
	return &output
}

func expandConnectivityConfIsGlobal(input bool) *connectivityconfigurations.IsGlobal {
	output := connectivityconfigurations.IsGlobalFalse
	if input {
		output = connectivityconfigurations.IsGlobalTrue
	}
	return &output
}

func expandConnectivityGroupItemModel(inputList []ConnectivityGroupItemModel) []connectivityconfigurations.ConnectivityGroupItem {
	var outputList []connectivityconfigurations.ConnectivityGroupItem
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.ConnectivityGroupItem{
			GroupConnectivity: input.GroupConnectivity,
			IsGlobal:          expandConnectivityConfIsGlobal(input.GlobalMeshEnabled),
			NetworkGroupId:    input.NetworkGroupId,
			UseHubGateway:     expandUseHubGateWay(input.UseHubGateway),
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func expandUseHubGateWay(input bool) *connectivityconfigurations.UseHubGateway {
	output := connectivityconfigurations.UseHubGatewayFalse
	if input {
		output = connectivityconfigurations.UseHubGatewayTrue
	}
	return &output
}

func expandHubModel(inputList []HubModel) *[]connectivityconfigurations.Hub {
	var outputList []connectivityconfigurations.Hub
	for _, v := range inputList {
		input := v
		output := connectivityconfigurations.Hub{
			ResourceId:   utils.String(input.ResourceId),
			ResourceType: utils.String(input.ResourceType),
		}

		outputList = append(outputList, output)
	}

	return &outputList
}

func flattenDeleteExistingPeering(input *connectivityconfigurations.DeleteExistingPeering) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.DeleteExistingPeeringTrue
}

func flattenConnectivityConfIsGlobal(input *connectivityconfigurations.IsGlobal) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.IsGlobalTrue
}

func flattenConnectivityGroupItemModel(inputList []connectivityconfigurations.ConnectivityGroupItem) []ConnectivityGroupItemModel {
	var outputList []ConnectivityGroupItemModel

	for _, input := range inputList {
		output := ConnectivityGroupItemModel{
			GroupConnectivity: input.GroupConnectivity,
			UseHubGateway:     flattenUseHubGateWay(input.UseHubGateway),
			GlobalMeshEnabled: flattenConnectivityConfIsGlobal(input.IsGlobal),
			NetworkGroupId:    input.NetworkGroupId,
		}

		outputList = append(outputList, output)
	}

	return outputList
}

func flattenUseHubGateWay(input *connectivityconfigurations.UseHubGateway) bool {
	if input == nil {
		return false
	}
	return *input == connectivityconfigurations.UseHubGatewayTrue
}

func flattenHubModel(inputList *[]connectivityconfigurations.Hub) []HubModel {
	var outputList []HubModel
	if inputList == nil {
		return outputList
	}

	for _, input := range *inputList {
		output := HubModel{}

		if input.ResourceId != nil {
			output.ResourceId = *input.ResourceId
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		outputList = append(outputList, output)
	}

	return outputList
}
