// Copyright IBM Corp. 2014, 2026
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
			Description:  "The name of the Public IP.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"allocation_method": {
			Description:  "The allocation method for this IP address.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForIPAllocationMethod(), false),
		},

		// Optional
		"ddos_protection_mode": {
			Description:  "The DDoS protection mode of the public IP.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(publicipaddresses.DdosSettingsProtectionModeVirtualNetworkInherited),
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForDdosSettingsProtectionMode(), false),
		},

		"ddos_protection_plan_id": {
			Description:  "The ID of DDoS protection plan associated with the public IP.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"ip_version": {
			Description:  "The IP Version to use, IPv6 or IPv4.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(publicipaddresses.IPVersionIPvFour),
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForIPVersion(), false),
		},

		"sku": {
			Description: "The SKU of the Public IP.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			ForceNew:    true,
			Default:     string(publicipaddresses.PublicIPAddressSkuNameStandard),
			// https://azure.microsoft.com/en-us/updates/upgrade-to-standard-sku-public-ip-addresses-in-azure-by-30-september-2025-basic-sku-will-be-retired/
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressSkuName(), false),
		},

		"sku_tier": {
			Description:  "The SKU Tier that should be used for the Public IP.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(publicipaddresses.PublicIPAddressSkuTierRegional),
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressSkuTier(), false),
		},

		"idle_timeout_in_minutes": {
			Description:  "The timeout for the TCP idle connection.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"domain_name_label": {
			Description:  "Label for the Domain Name. Will be used to make up the FQDN. If a domain name label is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.PublicIpDomainNameLabel,
		},

		"domain_name_label_scope": {
			Description:  "Scope for the domain name label. If a domain name label scope is specified, an A DNS record is created for the public IP in the Microsoft Azure DNS system with a hashed value includes in FQDN.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(publicipaddresses.PossibleValuesForPublicIPAddressDnsSettingsDomainNameLabelScope(), false),
		},

		"fqdn": {
			Description: "Fully qualified domain name of the A DNS record associated with the public IP. `domain_name_label` must be specified to get the `fqdn`. This is the concatenation of the `domain_name_label` and the regionalized DNS zone.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"reverse_fqdn": {
			Description: "A fully qualified domain name that resolves to this public IP address. If the reverseFqdn is specified, then a PTR DNS record is created pointing from the IP address in the in-addr.arpa domain to the reverse FQDN.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"ip_address": {
			Description: "The IP address value that was allocated.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"public_ip_prefix_id": {
			Description:  "If specified then public IP address allocated will be provided from the public IP prefix resource.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
		},

		"ip_tags": {
			Description: "A mapping of IP tags to assign to the public IP.",
			Type:        pluginsdk.TypeMap,
			Optional:    true,
			ForceNew:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
