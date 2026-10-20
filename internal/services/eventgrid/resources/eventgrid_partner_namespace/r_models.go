// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package eventgrid_partner_namespace

type EventGridPartnerNamespaceResourceModel struct {
	PartnerNamespaceName                string                      `tfschema:"name"`
	ResourceGroup                       string                      `tfschema:"resource_group_name"`
	Location                            string                      `tfschema:"location"`
	InboundIPRules                      []PartnerInboundIpRuleModel `tfschema:"inbound_ip_rule"`
	LocalAuthEnabled                    bool                        `tfschema:"local_authentication_enabled"`
	PartnerRegistrationFullyQualifiedID string                      `tfschema:"partner_registration_id"`
	PartnerTopicRoutingMode             string                      `tfschema:"partner_topic_routing_mode"`
	PublicNetworkAccess                 string                      `tfschema:"public_network_access"`
	Endpoint                            string                      `tfschema:"endpoint"`
	Tags                                map[string]string           `tfschema:"tags"`
}

type PartnerInboundIpRuleModel struct {
	IpMask string `tfschema:"ip_mask"`
	Action string `tfschema:"action"`
}

func (r Resource) ModelObject() any {
	return &EventGridPartnerNamespaceResourceModel{}
}
