// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package express_route_port

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/expressrouteports"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/expressroute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func expressRoutePortResourceSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ExpressRoutePortName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"peering_location": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"bandwidth_in_gbps": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"encapsulation": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(expressrouteports.PossibleValuesForExpressRoutePortsEncapsulation(), false),
		},

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"billing_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(expressrouteports.ExpressRoutePortsBillingTypeMeteredData),
			ValidateFunc: validation.StringInSlice(expressrouteports.PossibleValuesForExpressRoutePortsBillingType(), false),
		},

		"link1": expressRoutePortSchema,

		"link2": expressRoutePortSchema,

		"ethertype": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"mtu": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}

var expressRoutePortSchema = &pluginsdk.Schema{
	Type:     pluginsdk.TypeList,
	Optional: true,
	// Note: O+C because the service will always create a pair of links automatically. Users can't add or remove link, but only manipulate existing ones.
	// This is because the link is actually a map to the physical pair of ports on the MS edge device.
	Computed: true,
	MinItems: 1,
	MaxItems: 1,
	Elem: &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"admin_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
			"macsec_cipher": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(expressrouteports.ExpressRouteLinkMacSecCipherGcmAesOneTwoEight),
				ValidateFunc: validation.StringInSlice(expressrouteports.PossibleValuesForExpressRouteLinkMacSecCipher(), false),
			},
			"macsec_ckn_keyvault_secret_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"macsec_cak_keyvault_secret_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"macsec_sci_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},
			"id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"router_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"interface_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"patch_panel_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"rack_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"connector_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	},
}
