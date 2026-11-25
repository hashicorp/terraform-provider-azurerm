// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package lb

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/loadbalancers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/publicipprefixes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func lbSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Load Balancer.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"sku": {
			Description:  "The SKU of the Azure Load Balancer.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.LoadBalancerSkuNameStandard),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForLoadBalancerSkuName(), false),
		},

		"sku_tier": {
			Description:  "`sku_tier` - The SKU tier of this Load Balancer.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(loadbalancers.LoadBalancerSkuTierRegional),
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForLoadBalancerSkuTier(), false),
		},

		"frontend_ip_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the frontend IP configuration.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"subnet_id": {
						Description:  "The ID of the Subnet which should be associated with the IP Configuration.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidateSubnetID,
					},

					"private_ip_address": {
						Description: "Private IP Address to assign to the Load Balancer. The last one and first four IPs in any range are reserved and cannot be manually assigned.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Note: O+C because Azure assigns a private IP from the subnet when not specified
						Computed: true,
						ValidateFunc: validation.Any(
							validation.IsIPAddress,
							validation.StringIsEmpty,
						),
					},

					"private_ip_address_version": {
						Description: "The version of IP that the Private IP Address is.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Not using O+C here causes drift
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: validation.StringInSlice(loadbalancers.PossibleValuesForIPVersion(), false),
					},

					"public_ip_address_id": {
						Description:  "The ID of a Public IP Address which should be associated with the Load Balancer.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},

					"public_ip_prefix_id": {
						Description:  "The ID of a Public IP Prefix which should be associated with the Load Balancer. Public IP Prefix can only be used with outbound rules.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: publicipprefixes.ValidatePublicIPPrefixID,
					},

					"private_ip_address_allocation": {
						Description:      "The allocation method for the Private IP Address used by this Load Balancer.",
						Type:             pluginsdk.TypeString,
						Optional:         true,
						Computed:         true, // azignore:AZS007 - pre-existing violation
						ValidateFunc:     validation.StringInSlice(loadbalancers.PossibleValuesForIPAllocationMethod(), true),
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"gateway_load_balancer_frontend_ip_configuration_id": {
						Description:  "The Frontend IP Configuration ID of a Gateway SKU Load Balancer.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true, // azignore:AZS007 - pre-existing violation
						ValidateFunc: loadbalancers.ValidateFrontendIPConfigurationID,
					},

					"load_balancer_rules": {
						Description: "The list of IDs of load balancing rules that use this frontend IP.",
						Type:        pluginsdk.TypeSet,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"inbound_nat_rules": {
						Description: "The list of IDs of inbound rules that use this frontend IP.",
						Type:        pluginsdk.TypeSet,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"outbound_rules": {
						Description: "The list of IDs outbound rules that use this frontend IP.",
						Type:        pluginsdk.TypeSet,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},

					"zones": commonschema.ZonesMultipleOptional(),

					"id": {
						Description: "The ID of the Frontend IP Configuration.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"private_ip_address": {
			Description: "The first private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
		"private_ip_addresses": {
			Description: "The list of private IP address assigned to the load balancer in `frontend_ip_configuration` blocks, if any.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"tags": commonschema.Tags(),
	}
}
