// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package public_ip

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipaddresses"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func publicIpSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"allocation_method": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForIPAllocationMethod(), false),
		},

		// Optional
		"ddos_protection_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForDdosSettingsProtectionMode(), false),
			Default:      string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited),
		},

		"ddos_protection_plan_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"ip_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipaddresses.IPVersionIPvFour),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForIPVersion(), false),
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(publicipaddresses.PublicIPAddressSkuNameStandard),
			// https://azure.microsoft.com/en-us/updates/upgrade-to-standard-sku-public-ip-addresses-in-azure-by-30-september-2025-basic-sku-will-be-retired/
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressSkuName(), false),
		},

		"sku_tier": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipaddresses.PublicIPAddressSkuTierRegional),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressSkuTier(), false),
		},

		"idle_timeout_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"domain_name_label": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.PublicIpDomainNameLabel,
		},

		"domain_name_label_scope": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressDnsSettingsDomainNameLabelScope(), false),
		},

		"fqdn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"reverse_fqdn": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"public_ip_prefix_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
		},

		"ip_tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
