// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_routing_rule

func (Resource) ModelObject() any {
	return &ManagerRoutingRuleResourceModel{}
}

type ManagerRoutingRuleResourceModel struct {
	Description      string                          `tfschema:"description"`
	Destination      []ManagerRoutingRuleDestination `tfschema:"destination"`
	Name             string                          `tfschema:"name"`
	NextHop          []ManagerRoutingRuleNextHop     `tfschema:"next_hop"`
	RuleCollectionId string                          `tfschema:"rule_collection_id"`
}
