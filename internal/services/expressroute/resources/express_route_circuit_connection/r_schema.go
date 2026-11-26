// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit_connection

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/expressroute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRouteCircuitConnectionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Express Route Circuit Connection.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ExpressRouteCircuitConnectionName,
		},

		"peering_id": {
			Description:  "The ID of the Express Route Circuit Private Peering that this Express Route Circuit Connection connects with.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
		},

		"peer_peering_id": {
			Description:  "The ID of the peered Express Route Circuit Private Peering.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateExpressRouteCircuitPeeringID,
		},

		"address_prefix_ipv4": {
			Description:  "The IPv4 address space from which to allocate customer address for global reach.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsCIDR,
		},

		"authorization_key": {
			Description:  "The authorization key which is associated with the Express Route Circuit Connection.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.IsUUID,
		},

		"address_prefix_ipv6": {
			Description:  "The IPv6 address space from which to allocate customer addresses for global reach.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsCIDR,
		},
	}
}
