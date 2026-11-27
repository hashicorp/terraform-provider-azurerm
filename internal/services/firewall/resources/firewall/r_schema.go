// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package firewall

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/azurefirewalls"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/firewallpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/firewall/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func firewallSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FirewallName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		// lintignore:S013
		"sku_name": {
			Description:  "SKU name of the Firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallSkuName(), false),
		},

		// lintignore:S013
		"sku_tier": {
			Description:  "SKU tier of the Firewall.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallSkuTier(), false),
		},

		"firewall_policy_id": {
			Description:  "The ID of the Firewall Policy applied to this Firewall.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: firewallpolicies.ValidateFirewallPolicyID,
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subnet_id": {
						Description:  "Reference to the subnet associated with the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validate.FirewallSubnetName,
					},
					"public_ip_address_id": {
						Description:  "The ID of the Public IP Address associated with the firewall.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},
					"private_ip_address": {
						Description: "The Private IP address of the Azure Firewall.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"management_ip_configuration": {
			Description: "Adding or removing this block or changing the `subnet_id` in an existing block forces a new resource to be created.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			ForceNew:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"subnet_id": {
						Description:  "Reference to the subnet associated with the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validate.FirewallManagementSubnetName,
					},
					"public_ip_address_id": {
						Description:  "The ID of the Public IP Address associated with the firewall.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},
					"private_ip_address": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"threat_intel_mode": {
			Description:  "The operation mode for threat intelligence-based filtering.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice(azurefirewalls.PossibleValuesForAzureFirewallThreatIntelMode(), false),
		},

		"dns_servers": {
			Description: "A list of DNS servers that the Azure Firewall will direct DNS traffic for name resolution.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPAddress,
			},
		},

		"dns_proxy_enabled": {
			Description: "Whether DNS proxy is enabled. It will forward DNS requests to the DNS servers when set to `true`. It will be set to `true` if `dns_servers` provided with a not empty list.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"private_ip_ranges": {
			Description: "A list of SNAT private CIDR IP ranges, or the special string `IANAPrivateRanges`, which indicates Azure Firewall does not SNAT when the destination IP address is a private range per IANA RFC 1918.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.Any(
					validation.IsCIDR,
					validation.StringInSlice([]string{"IANAPrivateRanges"}, false),
				),
			},
		},

		"virtual_hub": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"virtual_hub_id": {
						Description:  "The ID of the Virtual Hub where the Firewall resides in.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: virtualwans.ValidateVirtualHubID,
					},
					"public_ip_count": {
						Description:  "The number of public IPs to assign to the Firewall.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      1,
						ValidateFunc: validation.IntAtLeast(1),
					},
					"public_ip_addresses": {
						Description: "The list of public IP addresses associated with the Firewall.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem:        &pluginsdk.Schema{Type: pluginsdk.TypeString},
					},
					"private_ip_address": {
						Description: "The private IP address associated with the Firewall.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"zones": commonschema.ZonesMultipleOptionalForceNew(),

		"tags": commonschema.Tags(),
	}
}
