// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ddosprotectionplans"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceVirtualNetworkSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"address_space": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			ExactlyOneOf: []string{"address_space", "ip_address_pool"},
			MinItems:     1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				// If `ip_address_pool` is used instead of `address_space` there is a perpetual diff
				// due to the API returning a CIDR range provisioned by the IP Address Management Pool.
				// Note: using `GetRawConfig` to avoid suppressing a diff if a user updates from `ip_address_pool` to `address_space`.
				rawIpAddressPool := d.GetRawConfig().AsValueMap()["ip_address_pool"]
				if !rawIpAddressPool.IsNull() && len(rawIpAddressPool.AsValueSlice()) > 0 {
					return true
				}

				return false
			},
		},

		"bgp_community": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.VirtualNetworkBgpCommunity,
		},

		"ddos_protection_plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
					},

					"enable": {
						Type:     pluginsdk.TypeBool,
						Required: true,
					},
				},
			},
		},

		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"enforcement": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworks.PossibleValuesForVirtualNetworkEncryptionEnforcement(), false),
					},
				},
			},
		},

		"dns_servers": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"flow_timeout_in_minutes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"ip_address_pool": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     2,
			ExactlyOneOf: []string{"address_space", "ip_address_pool"},
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

		"guid": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet": {
			Type:       pluginsdk.TypeSet,
			Optional:   true,
			Computed:   true, // azignore:AZS007 - pre-existing violation
			ConfigMode: pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefixes": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MinItems: 1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"default_outbound_access_enabled": {
						Type:     pluginsdk.TypeBool,
						Default:  true,
						Optional: true,
					},

					"delegation": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						MaxItems:   1,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"service_delegation": {
									Type:       pluginsdk.TypeList,
									Required:   true,
									MaxItems:   1,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
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

					"route_table_id": commonschema.ResourceIDReferenceOptional(&routetables.RouteTableId{}),

					"security_group": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"service_endpoint": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
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
					},

					"service_endpoint_policy_ids": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
						},
					},

					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"private_endpoint_vnet_policies": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(virtualnetworks.PrivateEndpointVNetPoliciesDisabled),
			ValidateFunc: validation.StringInSlice(virtualnetworks.PossibleValuesForPrivateEndpointVNetPolicies(), false),
		},

		"tags": commonschema.Tags(),
	}
}
