// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2026-03-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRouteCircuitSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the ExpressRoute circuit.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"sku": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"tier": {
						Description:  "The service tier.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(expressroutecircuits.PossibleValuesForExpressRouteCircuitSkuTier(), false),
					},

					"family": {
						Description:  "The billing mode for bandwidth.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(expressroutecircuits.PossibleValuesForExpressRouteCircuitSkuFamily(), false),
					},
				},
			},
		},

		"allow_classic_operations": {
			Description: "Allow the circuit to interact with classic (RDFE) resources.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"service_provider_name": {
			Description:      "The name of the ExpressRoute Service Provider.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			RequiredWith:     []string{"bandwidth_in_mbps", "peering_location"},
			ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"peering_location": {
			Description:      "The name of the peering location and **not** the Azure resource location.",
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			RequiredWith:     []string{"bandwidth_in_mbps", "service_provider_name"},
			ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"bandwidth_in_mbps": {
			Description:   "The bandwidth in Mbps of the circuit being created on the Service Provider.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			RequiredWith:  []string{"peering_location", "service_provider_name"},
			ConflictsWith: []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"bandwidth_in_gbps": {
			Description:   "The bandwidth in Gbps of the circuit being created on the Express Route Port.",
			Type:          pluginsdk.TypeFloat,
			Optional:      true,
			RequiredWith:  []string{"express_route_port_id"},
			ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
		},

		"express_route_port_id": {
			Description:   "The ID of the Express Route Port this Express Route Circuit is based on.",
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ValidateFunc:  expressrouteports.ValidateExpressRoutePortID,
			RequiredWith:  []string{"bandwidth_in_gbps"},
			ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
		},

		"rate_limiting_enabled": {
			Description: "Enable [rate limiting](https://learn.microsoft.com/azure/expressroute/rate-limit) for the circuit. Only works with ExpressRoute Ports.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"service_provider_provisioning_state": {
			Description: "The ExpressRoute circuit provisioning state from your chosen service provider.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"authorization_key": {
			Description: "The authorization key. This can be used to set up an ExpressRoute Circuit with an ExpressRoute Port from another subscription.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Sensitive:   true,
		},

		"service_key": {
			Description: "The string needed by the service provider to provision the ExpressRoute circuit.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"tags": commonschema.Tags(),
	}
}
