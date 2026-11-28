// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func localNetworkGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Local Network Gateway.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"gateway_address": {
			Description:  "The gateway IP address to connect with.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"gateway_address", "gateway_fqdn"},
		},

		"gateway_fqdn": {
			Description:  "The gateway FQDN to connect with.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ExactlyOneOf: []string{"gateway_address", "gateway_fqdn"},
		},

		"address_space": {
			Description: "The list of string CIDRs representing the address spaces the gateway exposes.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"bgp_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"asn": {
						Description: "The BGP speaker's ASN.",
						Type:        pluginsdk.TypeInt,
						Required:    true,
					},

					"bgp_peering_address": {
						Description: "The BGP peering address and BGP identifier of this BGP speaker.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"peer_weight": {
						Description: "The weight added to routes learned from this BGP speaker.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
