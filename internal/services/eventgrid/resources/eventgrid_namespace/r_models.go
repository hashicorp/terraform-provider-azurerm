// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_namespace

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
)

type EventGridNamespaceResourceModel struct {
	Name                     string                                     `tfschema:"name"`
	Location                 string                                     `tfschema:"location"`
	ResourceGroup            string                                     `tfschema:"resource_group_name"`
	Capacity                 int64                                      `tfschema:"capacity"`
	InboundIpRules           []InboundIpRuleModel                       `tfschema:"inbound_ip_rule"`
	Identity                 []identity.ModelSystemAssignedUserAssigned `tfschema:"identity"`
	PublicNetworkAccess      string                                     `tfschema:"public_network_access"`
	Sku                      string                                     `tfschema:"sku"`
	TopicSpacesConfiguration []TopicSpacesConfigurationModel            `tfschema:"topic_spaces_configuration"`
	Tags                     map[string]string                          `tfschema:"tags"`
}

type InboundIpRuleModel struct {
	IpMask string `tfschema:"ip_mask"`
	Action string `tfschema:"action"`
}

type TopicSpacesConfigurationModel struct {
	AlternativeAuthenticationNameSources       []string                 `tfschema:"alternative_authentication_name_source"`
	MaximumClientSessionsPerAuthenticationName int64                    `tfschema:"maximum_client_sessions_per_authentication_name"`
	MaximumSessionExpiryInHours                int64                    `tfschema:"maximum_session_expiry_in_hours"`
	RouteTopicResourceId                       string                   `tfschema:"route_topic_id"`
	DynamicRoutingEnrichment                   []RoutingEnrichmentModel `tfschema:"dynamic_routing_enrichment"`
	StaticRoutingEnrichment                    []RoutingEnrichmentModel `tfschema:"static_routing_enrichment"`
}

type RoutingEnrichmentModel struct {
	Key   string `tfschema:"key"`
	Value string `tfschema:"value"`
}

func (r Resource) ModelObject() any {
	return &EventGridNamespaceResourceModel{}
}
