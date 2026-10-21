// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func subnetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"virtual_network_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"address_prefixes": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "ip_address_pool"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				// If `ip_address_pool` is used instead of `address_prefixes` there is a perpetual diff
				// due to the API returning a CIDR range provisioned by the IP Address Management Pool.
				// Note: using `GetRawConfig` to avoid suppressing a diff if a user updates from `ip_address_pool` to `address_prefixes`.
				rawIpAddressPool := d.GetRawConfig().AsValueMap()["ip_address_pool"]
				if !rawIpAddressPool.IsNull() && len(rawIpAddressPool.AsValueSlice()) > 0 {
					return true
				}

				return false
			},
		},

		"service_endpoint": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"service": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.SubnetServiceEndpointName(),
					},
					"network_identifier": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
				},
			},
			// Temporary solution, see function notes for detail.
			DiffSuppressFunc: subnetServiceEndpointDiffSuppress,
		},

		"service_endpoint_policy_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
			},
		},

		"sharing_scope": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// todo "Tenant" is only supported until https://github.com/Azure/azure-rest-api-specs/issues/36446 is addressed
			ValidateFunc: validation.StringInSlice([]string{string(subnets.SharingScopeTenant)}, false),
		},

		"delegation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"service_delegation": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(subnetDelegationServiceNames, false),
								},

								"actions": {
									Type:     pluginsdk.TypeSet,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Microsoft.Network/networkinterfaces/*",
											"Microsoft.Network/publicIPAddresses/join/action",
											"Microsoft.Network/publicIPAddresses/read",
											"Microsoft.Network/virtualNetworks/read",
											"Microsoft.Network/virtualNetworks/subnets/action",
											"Microsoft.Network/virtualNetworks/subnets/join/action",
											"Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
											"Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
										}, false),
									},
								},
							},
						},
					},
				},
			},
		},

		"ip_address_pool": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"address_prefixes", "ip_address_pool"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: ipampools.ValidateIPamPoolID,
					},

					"number_of_ip_addresses": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[1-9]\d*$`),
							"`number_of_ip_addresses` must be a string that represents a positive number",
						),
					},

					"allocated_ip_address_prefixes": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"default_outbound_access_enabled": {
			Type:     pluginsdk.TypeBool,
			Default:  true,
			Optional: true,
		},

		"private_endpoint_network_policies": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
			ValidateFunc: validation.StringInSlice(subnets.PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies(), false),
		},

		"private_link_service_network_policies_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"network_security_group_id_wo": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			WriteOnly:    true,
			RequiredWith: []string{"network_security_group_id_wo_version"},
			ValidateFunc: networksecuritygroups.ValidateNetworkSecurityGroupID,
		},

		"network_security_group_id_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"network_security_group_id_wo"},
			ValidateFunc: validation.IntAtLeast(1),
		},

		"network_security_group_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"route_table_id_wo": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			WriteOnly:    true,
			RequiredWith: []string{"route_table_id_wo_version"},
			ValidateFunc: routetables.ValidateRouteTableID,
		},

		"route_table_id_wo_version": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			RequiredWith: []string{"route_table_id_wo"},
			ValidateFunc: validation.IntAtLeast(1),
		},

		"route_table_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
