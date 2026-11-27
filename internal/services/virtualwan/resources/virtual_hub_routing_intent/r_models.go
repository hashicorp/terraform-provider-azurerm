// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_routing_intent

type VirtualHubRoutingIntentModel struct {
	Name            string          `tfschema:"name"`
	VirtualHubId    string          `tfschema:"virtual_hub_id"`
	RoutingPolicies []RoutingPolicy `tfschema:"routing_policy"`
}

func (r Resource) ModelObject() any {
	return &VirtualHubRoutingIntentModel{}
}
