// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_custom_domain_association

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/afddomains"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/routes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func cdnFrontdoorCustomDomainAssociationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"cdn_frontdoor_custom_domain_id": {
			Description:  "The ID of the Front Door Custom Domain that should be managed by the association resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: afddomains.ValidateCustomDomainID,
		},

		"cdn_frontdoor_route_ids": {
			Description: "One or more IDs of the Front Door Route to which the Front Door Custom Domain is associated with.",
			Type:        pluginsdk.TypeList,
			Required:    true,

			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: routes.ValidateRouteID,
			},
		},
	}
}
