// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package palo_alto_next_generation_firewall_virtual_network_local_rulestack

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2025-10-08/localrulestackresources"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.NextGenerationFirewallName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"rulestack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: localrulestackresources.ValidateLocalRulestackID,
		},

		"network_profile": schema.VnetNetworkProfileSchema(),

		// Optional
		"dns_settings": schema.DNSSettingsSchema(),

		"destination_nat": schema.DestinationNATSchema(),

		"marketplace_offer_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      "pan_swfw_cloud_ngfw",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "panw-cngfw-payg",
			ValidateFunc: validation.StringLenBetween(1, 50),
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
