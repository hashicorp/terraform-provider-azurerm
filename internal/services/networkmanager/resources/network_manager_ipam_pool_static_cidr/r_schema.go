// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package network_manager_ipam_pool_static_cidr

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/staticcidrs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkmanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name which should be used for this Network Manager IPAM Pool Static CIDR.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"ipam_pool_id": commonschema.ResourceIDReferenceRequiredForceNew(&staticcidrs.IPamPoolId{}, "Network Manager IP Address Management (IPAM) Pool"),

		"address_prefixes": {
			Description: "A list of IPv4 or IPv6 IP address prefixes which will be allocated to the Static CIDR.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
				// If `number_of_ip_addresses_to_allocate` is used instead of `address_prefixes` there is a perpetual diff
				// due to the API returning a CIDR range provisioned by the IP Address Management Pool.
				// Note: using `GetRawConfig` to avoid suppressing a diff if a user updates from `address_prefixes` to `number_of_ip_addresses_to_allocate`.
				rawNumberOfIpAddressesToAllocate := d.GetRawConfig().AsValueMap()["number_of_ip_addresses_to_allocate"]
				if !rawNumberOfIpAddressesToAllocate.IsNull() && rawNumberOfIpAddressesToAllocate.AsString() != "" {
					return true
				}

				return false
			},
			ExactlyOneOf: []string{"address_prefixes", "number_of_ip_addresses_to_allocate"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"number_of_ip_addresses_to_allocate": {
			Description:  "The number of IP addresses to allocate to the Static CIDR.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.NumberOfIpAddresses,
			DiffSuppressFunc: func(_, old, new string, d *pluginsdk.ResourceData) bool {
				rawAddressPrefixes := d.GetRawConfig().AsValueMap()["address_prefixes"]
				if !rawAddressPrefixes.IsNull() && len(rawAddressPrefixes.AsValueSlice()) > 0 {
					return true
				}

				return false
			},
			ExactlyOneOf: []string{"address_prefixes", "number_of_ip_addresses_to_allocate"},
		},
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
