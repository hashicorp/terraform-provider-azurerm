// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_hub_route_table_route

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualHubRouteTableRouteSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"route_table_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: virtualwans.ValidateHubRouteTableID,
		},

		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"destinations": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"destinations_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"CIDR",
				"ResourceId",
				"Service",
			}, false),
		},

		"next_hop": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"next_hop_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "ResourceId",
			ValidateFunc: validation.StringInSlice([]string{
				"ResourceId",
			}, false),
		},
	}
}
