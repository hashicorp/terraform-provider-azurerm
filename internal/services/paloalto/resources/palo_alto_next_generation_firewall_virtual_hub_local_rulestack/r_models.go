// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_hub_local_rulestack

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

type NextGenerationFirewallVHubLocalRuleStackModel struct {
	Name              string `tfschema:"name"`
	ResourceGroupName string `tfschema:"resource_group_name"`
	RuleStackId       string `tfschema:"rulestack_id"`

	NetworkProfile     []schema.NetworkProfileVHub `tfschema:"network_profile"`
	DNSSettings        []schema.DNSSettings        `tfschema:"dns_settings"`
	FrontEnd           []schema.DestinationNAT     `tfschema:"destination_nat"`
	MarketplaceOfferId string                      `tfschema:"marketplace_offer_id"`
	PlanId             string                      `tfschema:"plan_id"`
	Tags               map[string]any              `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &NextGenerationFirewallVHubLocalRuleStackModel{}
}
