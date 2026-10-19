// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_circuit

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressroutecircuits"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRouteCircuitSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(expressroutecircuits.PossibleValuesForExpressRouteCircuitSkuTier(), false),
					},

					"family": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(expressroutecircuits.PossibleValuesForExpressRouteCircuitSkuFamily(), false),
					},
				},
			},
		},

		"allow_classic_operations": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_provider_name": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			RequiredWith:     []string{"bandwidth_in_mbps", "peering_location"},
			ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"peering_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			RequiredWith:     []string{"bandwidth_in_mbps", "service_provider_name"},
			ConflictsWith:    []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"bandwidth_in_mbps": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			RequiredWith:  []string{"peering_location", "service_provider_name"},
			ConflictsWith: []string{"bandwidth_in_gbps", "express_route_port_id"},
		},

		"bandwidth_in_gbps": {
			Type:          pluginsdk.TypeFloat,
			Optional:      true,
			RequiredWith:  []string{"express_route_port_id"},
			ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
		},

		"express_route_port_id": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			RequiredWith:  []string{"bandwidth_in_gbps"},
			ConflictsWith: []string{"bandwidth_in_mbps", "peering_location", "service_provider_name"},
			ValidateFunc:  expressrouteports.ValidateExpressRoutePortID,
		},

		"rate_limiting_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"service_provider_provisioning_state": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"authorization_key": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			Sensitive: true,
		},

		"service_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": commonschema.Tags(),
	}
}
