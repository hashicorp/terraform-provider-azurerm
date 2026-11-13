// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_strata_cloud_manager

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
)

type NextGenerationFirewallVNetStrataCloudManagerModel struct {
	Name                         string                       `tfschema:"name"`
	ResourceGroupName            string                       `tfschema:"resource_group_name"`
	Location                     string                       `tfschema:"location"`
	NetworkProfile               []schema.NetworkProfileVnet  `tfschema:"network_profile"`
	StrataCloudManagerTenantName string                       `tfschema:"strata_cloud_manager_tenant_name"`
	DNSSettings                  []schema.DNSSettings         `tfschema:"dns_settings"`
	FrontEnd                     []schema.DestinationNAT      `tfschema:"destination_nat"`
	MarketplaceOfferId           string                       `tfschema:"marketplace_offer_id"`
	PlanId                       string                       `tfschema:"plan_id"`
	Identity                     []identity.ModelUserAssigned `tfschema:"identity"`
	Tags                         map[string]any               `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &NextGenerationFirewallVNetStrataCloudManagerModel{}
}
