// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule_collection

func (Resource) ModelObject() any {
	return &ManagerRoutingRuleCollectionResourceModel{}
}

type ManagerRoutingRuleCollectionResourceModel struct {
	BgpRoutePropagationEnabled bool     `tfschema:"bgp_route_propagation_enabled"`
	Description                string   `tfschema:"description"`
	Name                       string   `tfschema:"name"`
	NetworkGroupIds            []string `tfschema:"network_group_ids"`
	RoutingConfigurationId     string   `tfschema:"routing_configuration_id"`
}
