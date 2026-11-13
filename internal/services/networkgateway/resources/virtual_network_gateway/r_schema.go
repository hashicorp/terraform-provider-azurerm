// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualNetworkGatewaySchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute),
				string(virtualnetworkgateways.VirtualNetworkGatewayTypeVpn),
			}, false),
		},

		"vpn_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualnetworkgateways.VpnTypeRouteBased),
			ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnType(), false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"bgp_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"private_ip_address_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"active_active": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			// This validator checks for all possible values for the SKU regardless of the attributes vpn_type and
			// type. For a validation which depends on the attributes vpn_type and type, refer to the special case
			// validators validateVirtualNetworkGatewayPolicyBasedVpnSku, validateVirtualNetworkGatewayRouteBasedVpnSku
			// and validateVirtualNetworkGatewayExpressRouteSku.
			ValidateFunc: validation.Any(
				validateVirtualNetworkGatewayPolicyBasedVpnSku(),
				validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1(),
				validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2(),
				validateVirtualNetworkGatewayExpressRouteSku(),
			),
		},

		"generation": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnGatewayGeneration(), false),
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
			// Each type gateway requires exact number of `ip_configuration`, and overwriting an existing one is not allowed.
			ForceNew: true,
			MaxItems: 3,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Azure Management API requires a name but does not generate a name if the field is missing
						// The name "vnetGatewayConfig" is used when creating a virtual network gateway via the
						// Azure portal.
						Default: "vnetGatewayConfig",
					},

					"private_ip_address_allocation": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPAllocationMethod(), false),
						Default:      string(virtualnetworkgateways.IPAllocationMethodDynamic),
					},

					"subnet_id": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validate.IsGatewaySubnet,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"public_ip_address_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: commonids.ValidatePublicIPAddressID,
					},
				},
			},
		},

		"policy_group": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PolicyGroupName,
					},

					"policy_member": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnPolicyMemberAttributeType(), false),
								},

								"value": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"is_default": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"priority": {
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"vpn_client_configuration": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"aad_tenant": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						RequiredWith: []string{
							"vpn_client_configuration.0.aad_audience",
							"vpn_client_configuration.0.aad_issuer",
						},
					},
					"aad_audience": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						RequiredWith: []string{
							"vpn_client_configuration.0.aad_issuer",
							"vpn_client_configuration.0.aad_tenant",
						},
					},
					"aad_issuer": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						RequiredWith: []string{
							"vpn_client_configuration.0.aad_audience",
							"vpn_client_configuration.0.aad_tenant",
						},
					},

					"virtual_network_gateway_client_connection": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"policy_group_names": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.PolicyGroupName,
									},
								},

								"address_prefixes": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"ipsec_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dh_group": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForDhGroup(), false),
								},

								"ike_encryption": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIkeEncryption(), false),
								},

								"ike_integrity": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIkeIntegrity(), false),
								},

								"ipsec_encryption": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPsecEncryption(), false),
								},

								"ipsec_integrity": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPsecIntegrity(), false),
								},

								"pfs_group": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForPfsGroup(), false),
								},

								"sa_lifetime_in_seconds": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(300, 172799),
								},

								"sa_data_size_in_kilobytes": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1024, math.MaxInt32),
								},
							},
						},
					},

					"root_certificate": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"public_cert_data": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
						Set: hashVirtualNetworkGatewayRootCert,
					},

					"revoked_certificate": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"thumbprint": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
						Set: hashVirtualNetworkGatewayRevokedCert,
					},

					"radius_server": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"secret": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringLenBetween(1, 128),
									Sensitive:    true,
									// not returned by API - This prevents a diff, however, the state value will be nil so cannot be exported
									// TODO - Convert this to an Write Only property?
									DiffSuppressFunc: func(k, oldValue, newValue string, d *pluginsdk.ResourceData) bool {
										return len(newValue) == 0
									},
								},

								"score": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 30),
								},
							},
						},
					},

					"radius_server_address": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
						RequiredWith: []string{"vpn_client_configuration.0.radius_server_secret"},
					},

					"radius_server_secret": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						RequiredWith: []string{"vpn_client_configuration.0.radius_server_address"},
					},

					"vpn_auth_types": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						MaxItems: 3,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnAuthenticationType(), false),
						},
					},

					"vpn_client_protocols": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnClientProtocol(), false),
						},
					},
				},
			},
		},

		"bgp_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"asn": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
					},

					"peer_weight": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
					},

					// lintignore:XS003
					"peering_addresses": {
						Type:     pluginsdk.TypeList,
						Computed: true, // azignore:AZS007 - pre-existing violation
						Optional: true,
						MinItems: 1,
						MaxItems: 2,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_configuration_name": {
									Type:     pluginsdk.TypeString,
									Optional: true,
									// Note: O+C In case there is only one `ip_configuration` in root level. This property can be deduced from the that.
									Computed:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"apipa_addresses": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									MinItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.IPAddressInAzureReservedAPIPARange,
									},
								},
								"default_addresses": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"tunnel_ip_addresses": {
									Type:     pluginsdk.TypeList,
									Computed: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
						AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
					},
				},
			},
		},

		// lintignore:XS003
		"custom_route": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefixes": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"default_local_network_gateway_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: localnetworkgateways.ValidateLocalNetworkGatewayID,
		},

		"bgp_route_translation_for_nat_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"dns_forwarding_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"ip_sec_replay_protection_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"maximum_scale_unit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(1, 40),
			RequiredWith: []string{"maximum_scale_unit", "minimum_scale_unit"},
		},

		"minimum_scale_unit": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(1, 40),
			RequiredWith: []string{"maximum_scale_unit", "minimum_scale_unit"},
		},

		"remote_vnet_traffic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"virtual_wan_traffic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": commonschema.Tags(),
	}
}
