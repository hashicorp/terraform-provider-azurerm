// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func resourceVirtualNetworkSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Virtual Network.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"address_space": {
			Description: "The address space that is used the Virtual Network. You can supply more than one address space.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				// If `ip_address_pool` is used instead of `address_space` there is a perpetual diff
				// due to the API returning a CIDR range provisioned by the IP Address Management Pool.
				// Note: using `GetRawConfig` to avoid suppressing a diff if a user updates from `ip_address_pool` to `address_space`.
				rawIpAddressPool := d.GetRawConfig().AsValueMap()["ip_address_pool"]
				if !rawIpAddressPool.IsNull() && (!rawIpAddressPool.IsKnown() || len(rawIpAddressPool.AsValueSlice()) > 0) {
					return true
				}

				return false
			},
			ExactlyOneOf: []string{"address_space", "ip_address_pool"},
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"advertised_gateway_prefixes": {
			Description: "A set of CIDR prefixes that Azure gateways advertise to on-premises networks.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsCIDR,
			},
		},

		"bgp_community": {
			Description:  "The BGP community attribute in format `<as-number>:<community-value>`.",
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
						Description:  "The ID of DDoS Protection Plan.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: ddosprotectionplans.ValidateDdosProtectionPlanID,
					},

					"enable": {
						Description: "Enable/disable DDoS Protection Plan on Virtual Network.",
						Type:        pluginsdk.TypeBool,
						Required:    true,
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
						Description:  "Whether the encrypted Virtual Network allows VM that does not support encryption.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworks.PossibleValuesForVirtualNetworkEncryptionEnforcement(), false),
					},
				},
			},
		},

		"dns_servers": {
			Description: "List of IP addresses of DNS servers.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"flow_timeout_in_minutes": {
			Description:  "The flow timeout in minutes for the Virtual Network, which is used to enable connection tracking for intra-VM flows.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(4, 30),
		},

		"ip_address_pool": {
			Description:  "Only one association of each IP type (IPv4 or IPv6) is allowed.",
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     2,
			ExactlyOneOf: []string{"address_space", "ip_address_pool"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Description:  "The ID of the Network Manager IP Address Management (IPAM) Pool.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: ipampools.ValidateIPamPoolID,
					},

					"number_of_ip_addresses": {
						Description: "The number of IP addresses to allocated to the Virtual Network.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^[1-9]\d*$`),
							"`number_of_ip_addresses` must be a string that represents a positive number",
						),
					},

					"allocated_ip_address_prefixes": {
						Description: "The list of IP address prefixes allocated to the Virtual Network.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"guid": {
			Description: "The GUID of the Virtual Network.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet": {
			Description: "Can be specified multiple times to define multiple subnets.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ConfigMode:  pluginsdk.SchemaConfigModeAttr,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the subnet.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"address_prefixes": {
						Description: "The address prefixes to use for the subnet.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						MinItems:    1,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"default_outbound_access_enabled": {
						Description: "Enable default outbound access to the internet for the subnet.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"delegation": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						MaxItems:   1,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "A name for this delegation.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"service_delegation": {
									Type:       pluginsdk.TypeList,
									Required:   true,
									MaxItems:   1,
									ConfigMode: pluginsdk.SchemaConfigModeAttr,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"name": {
												Description:  "The name of service to delegate to.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(helpers.SubnetDelegationServiceNames, false),
											},

											"actions": {
												Description: "A list of Actions which should be delegated. This list is specific to the service to delegate to.",
												Type:        pluginsdk.TypeSet,
												Optional:    true,
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
						Description:  "Enable or Disable network policies for the Private Endpoint on the subnet.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
						ValidateFunc: validation.StringInSlice(subnets.PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies(), false),
					},

					"private_link_service_network_policies_enabled": {
						Description: "Enable or Disable network policies for the Private Link Service on the subnet.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},

					"route_table_id": commonschema.ResourceIDReferenceOptional(&routetables.RouteTableId{}, "Route Table"),

					"security_group": {
						Description: "The Network Security Group to associate with the subnet. (Referenced by `id`, i.e. `azurerm_network_security_group.example.id`).",
						Type:        pluginsdk.TypeString,
						Optional:    true,
					},

					"service_endpoint": {
						Type:       pluginsdk.TypeList,
						Optional:   true,
						ConfigMode: pluginsdk.SchemaConfigModeAttr,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service": {
									Description:  "The name of the Service endpoint to associate with the subnet.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validate.SubnetServiceEndpointName(),
								},
								"network_identifier": {
									Description:  "The ARM resource ID of the network identifier to associate with the service endpoint.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: azure.ValidateResourceID,
								},
							},
						},
					},

					"service_endpoint_policy_ids": {
						Description: "The list of IDs of Service Endpoint Policies to associate with the subnet.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
						},
					},

					"id": {
						Description: "The ID of this subnet.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},
				},
			},
		},

		"private_endpoint_vnet_policies": {
			Description:  "The Private Endpoint VNet Policies for the Virtual Network.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(virtualnetworks.PrivateEndpointVNetPoliciesDisabled),
			ValidateFunc: validation.StringInSlice(virtualnetworks.PossibleValuesForPrivateEndpointVNetPolicies(), false),
		},

		"tags": commonschema.Tags(),
	}
}
