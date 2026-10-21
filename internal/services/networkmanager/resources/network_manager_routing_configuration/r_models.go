// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_configuration

func (Resource) ModelObject() any {
	return &ManagerRoutingConfigurationResourceModel{}
}

type ManagerRoutingConfigurationResourceModel struct {
	Description         string `tfschema:"description"`
	Name                string `tfschema:"name"`
	NetworkManagerId    string `tfschema:"network_manager_id"`
	RouteTableUsageMode string `tfschema:"route_table_usage_mode"`
}
