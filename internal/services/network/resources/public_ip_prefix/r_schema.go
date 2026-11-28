// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Public IP Prefix resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"custom_ip_prefix_id": {
			Description:  "The Custom IP Prefix ID associated with the Public IP Prefix.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: customipprefixes.ValidateCustomIPPrefixID,
		},

		"sku": {
			Description:  "The SKU of the Public IP Prefix.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipprefixes.PublicIPPrefixSkuNameStandard),
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForPublicIPPrefixSkuName(), false),
		},

		"sku_tier": {
			Description:  "The SKU Tier that should be used for the Public IP Prefix.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(publicipprefixes.PublicIPPrefixSkuTierRegional),
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForPublicIPPrefixSkuTier(), false),
		},

		"prefix_length": {
			Description:  "The number of bits of the prefix.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			Default:      28,
			ValidateFunc: validation.IntBetween(0, 127),
		},

		"ip_version": {
			Description:  "The IP Version to use, `IPv6` or `IPv4`.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(publicipprefixes.IPVersionIPvFour),
			ValidateFunc: validation.StringInSlice(publicipprefixes.PossibleValuesForIPVersion(), false),
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"ip_prefix": {
			Description: "The IP address prefix value that was allocated.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tags": commonschema.Tags(),
	}
}
