// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone_virtual_network_link

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func privateDnsZoneVirtualNetworkLinkSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Private DNS Zone Virtual Network Link.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"private_dns_zone_id": {
			Description:  "The ID of the Private DNS Zone where the resource exists.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
		},

		"virtual_network_id": {
			Description:  "The ID of the Virtual Network that should be linked to the DNS Zone.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualNetworkID,
		},

		"registration_enabled": {
			Description: "Whether auto-registration of Virtual Machine records in the Virtual Network in the Private DNS zone is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"resolution_policy": {
			Description: "The resolution policy of the Private DNS Zone Virtual Network Link.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// Note: O+C because when the `name` of `azurerm_private_dns_zone` is a private link endpoint, the service will set default value for this.
			Computed:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworklinks.PossibleValuesForResolutionPolicy(), false),
		},

		"tags": commonschema.Tags(),
	}
}
