// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_connectivity_configuration

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/connectivityconfigurations"
)

type ManagerConnectivityConfigurationModel struct {
	Name                                string                                          `tfschema:"name"`
	NetworkManagerId                    string                                          `tfschema:"network_manager_id"`
	AppliesToGroups                     []ConnectivityGroupItemModel                    `tfschema:"applies_to_group"`
	ConnectedGroupAddressOverlapEnabled bool                                            `tfschema:"connected_group_address_overlap_enabled"`
	ConnectedGroupPrivateEndpointsScale string                                          `tfschema:"connected_group_private_endpoints_scale"`
	ConnectivityTopology                connectivityconfigurations.ConnectivityTopology `tfschema:"connectivity_topology"`
	DeleteExistingPeeringEnabled        bool                                            `tfschema:"delete_existing_peering_enabled"`
	Description                         string                                          `tfschema:"description"`
	GlobalMeshEnabled                   bool                                            `tfschema:"global_mesh_enabled"`
	Hub                                 []HubModel                                      `tfschema:"hub"`
	PeeringEnforcementEnabled           bool                                            `tfschema:"peering_enforcement_enabled"`
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

func (r Resource) ModelObject() any {
	return &ManagerConnectivityConfigurationModel{}
}
