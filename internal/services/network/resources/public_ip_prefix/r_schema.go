// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip_prefix

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/customipprefixes"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func publicIpPrefixSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"custom_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: customipprefixes.ValidateCustomIPPrefixID,
		},

		"sku": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipprefixes.PublicIPPrefixSkuNameStandard),
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForPublicIPPrefixSkuName(), false),
		},

		"sku_tier": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipprefixes.PublicIPPrefixSkuTierRegional),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForPublicIPPrefixSkuTier(), false),
		},

		"prefix_length": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      28,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(0, 127),
		},

		"ip_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipprefixes.IPVersionIPvFour),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForIPVersion(), false),
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"ip_prefix": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": commonschema.Tags(),
	}
}
