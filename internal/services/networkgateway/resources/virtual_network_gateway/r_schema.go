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
			Description:  "The name of the Virtual Network Gateway.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"type": {
			Description: "The type of the Virtual Network Gateway.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute),
				string(virtualnetworkgateways.VirtualNetworkGatewayTypeVpn),
			}, false),
		},

		"vpn_type": {
			Description:  "The routing type of the Virtual Network Gateway.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(virtualnetworkgateways.VpnTypeRouteBased),
			ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnType(), false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		"bgp_enabled": {
			Description: "If `true`, BGP (Border Gateway Protocol) will be enabled for this Virtual Network Gateway.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"private_ip_address_enabled": {
			Description: "Whether private IP should be enabled on this gateway for connections.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
		},

		"active_active": {
			Description: "If `true`, an active-active Virtual Network Gateway will be created. An active-active gateway requires a `HighPerformance` or an `UltraPerformance` SKU. If `false`, an active-standby gateway will be created.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
		},

		"sku": {
			Description: "Configuration of the size and capacity of the Virtual Network Gateway. A `PolicyBased` gateway only supports the `Basic` SKU. Further, the `UltraPerformance` and `ErGwScale` SKU is only supported by an `ExpressRoute` gateway.",
			Type:        pluginsdk.TypeString,
			Required:    true,
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
			Description:  "The Generation of the Virtual Network Gateway.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnGatewayGeneration(), false),
		},

		"ip_configuration": {
			Description: "An active-standby gateway requires exactly one `ip_configuration` block, an active-active gateway requires exactly two `ip_configuration` blocks whereas an active-active zone redundant gateway with P2S configuration requires exactly three `ip_configuration` blocks.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			// Each type gateway requires exact number of `ip_configuration`, and overwriting an existing one is not allowed.
			ForceNew: true,
			MaxItems: 3,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description: "A user-defined name of the IP configuration.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						// Azure Management API requires a name but does not generate a name if the field is missing
						// The name "vnetGatewayConfig" is used when creating a virtual network gateway via the
						// Azure portal.
						Default: "vnetGatewayConfig",
					},

					"private_ip_address_allocation": {
						Description:  "How the private IP address of the gateway's virtual interface is assigned.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Default:      string(virtualnetworkgateways.IPAllocationMethodDynamic),
						ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPAllocationMethod(), false),
					},

					"subnet_id": {
						Description:      "The ID of the gateway subnet of a Virtual Network in which the Virtual Network Gateway will be created. It is mandatory that the associated subnet is named `GatewaySubnet`. Therefore, each Virtual Network can contain at most a single Virtual Network Gateway.",
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validate.IsGatewaySubnet,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"public_ip_address_id": {
						Description:  "The ID of the public IP address to associate with the Virtual Network Gateway.",
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
						Description:  "The name of the Virtual Network Gateway Policy Group.",
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
									Description:  "The name of the Virtual Network Gateway Policy Group Member.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"type": {
									Description:  "The VPN Policy Member attribute type.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnPolicyMemberAttributeType(), false),
								},

								"value": {
									Description:  "The value of attribute that is used for this Virtual Network Gateway Policy Group Member.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},

					"is_default": {
						Description: "Whether this is a Default Virtual Network Gateway Policy Group.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"priority": {
						Description:  "The priority for the Virtual Network Gateway Policy Group.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						Default:      0,
						ValidateFunc: validation.IntAtLeast(0),
					},
				},
			},
		},

		"vpn_client_configuration": {
			Description: "In this block the Virtual Network Gateway can be configured to accept IPSec point-to-site connections.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_space": {
						Description: "The address space out of which IP addresses for VPN clients will be taken. You can provide more than one address space, e.g. in CIDR notation.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"aad_tenant": {
						Description: "AzureAD Tenant URL.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						RequiredWith: []string{
							"vpn_client_configuration.0.aad_audience",
							"vpn_client_configuration.0.aad_issuer",
						},
					},
					"aad_audience": {
						Description: "The client ID of the Azure VPN application. See [Create an Active Directory (AD) tenant for P2S OpenVPN protocol connections](https://learn.microsoft.com/azure/vpn-gateway/openvpn-azure-ad-tenant-multi-app) for values.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
						RequiredWith: []string{
							"vpn_client_configuration.0.aad_issuer",
							"vpn_client_configuration.0.aad_tenant",
						},
					},
					"aad_issuer": {
						Description: "The STS URL for your tenant.",
						Type:        pluginsdk.TypeString,
						Optional:    true,
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
									Description:  "The name of the Virtual Network Gateway Client Connection.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"policy_group_names": {
									Description: "A list of names of Virtual Network Gateway Policy Groups.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.PolicyGroupName,
									},
								},

								"address_prefixes": {
									Description: "A list of address prefixes for P2S VPN Client.",
									Type:        pluginsdk.TypeList,
									Required:    true,
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
									Description:  "The DH Group, used in IKE Phase 1.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForDhGroup(), false),
								},

								"ike_encryption": {
									Description:  "The IKE encryption algorithm, used for IKE Phase 2.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIkeEncryption(), false),
								},

								"ike_integrity": {
									Description:  "The IKE encryption integrity algorithm, used for IKE Phase 2.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIkeIntegrity(), false),
								},

								"ipsec_encryption": {
									Description:  "The IPSec encryption algorithm, used for IKE phase 1.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPsecEncryption(), false),
								},

								"ipsec_integrity": {
									Description:  "The IPSec integrity algorithm, used for IKE phase 1.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForIPsecIntegrity(), false),
								},

								"pfs_group": {
									Description:  "The Pfs Group, used in IKE Phase 2.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForPfsGroup(), false),
								},

								"sa_lifetime_in_seconds": {
									Description:  "The IPSec Security Association lifetime in seconds for a Site-to-Site VPN tunnel.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(300, 172799),
								},

								"sa_data_size_in_kilobytes": {
									Description:  "The IPSec Security Association payload size in KB for a Site-to-Site VPN tunnel.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1024, math.MaxInt32),
								},
							},
						},
					},

					"root_certificate": {
						Description: "These root certificates are used to sign the client certificate used by the VPN clients to connect to the gateway.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Set:         hashVirtualNetworkGatewayRootCert,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "A user-defined name of the root certificate.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"public_cert_data": {
									Description: "The public certificate of the root certificate authority. The certificate must be provided in Base-64 encoded X.509 format (PEM). In particular, this argument *must not* include the `-----BEGIN CERTIFICATE-----` or `-----END CERTIFICATE-----` markers, nor any newlines.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},

					"revoked_certificate": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Set:      hashVirtualNetworkGatewayRevokedCert,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "The name of the certificate resource.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"thumbprint": {
									Description: "The public data of the certificate.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
							},
						},
					},

					"radius_server": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"address": {
									Description:  "The address of the Radius Server.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsIPv4Address,
								},

								"secret": {
									Description:  "The secret that is used to communicate with the Radius Server.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									Sensitive:    true,
									ValidateFunc: validation.StringLenBetween(1, 128),
									// not returned by API - This prevents a diff, however, the state value will be nil so cannot be exported
									// TODO - Convert this to an Write Only property?
									DiffSuppressFunc: func(k, oldValue, newValue string, d *pluginsdk.ResourceData) bool {
										return len(newValue) == 0
									},
								},

								"score": {
									Description:  "The score of the Radius Server determines the priority of the server.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 30),
								},
							},
						},
					},

					"radius_server_address": {
						Description:  "The address of the Radius server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsIPv4Address,
						RequiredWith: []string{"vpn_client_configuration.0.radius_server_secret"},
					},

					"radius_server_secret": {
						Description:  "The secret used by the Radius server.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						RequiredWith: []string{"vpn_client_configuration.0.radius_server_address"},
					},

					"vpn_auth_types": {
						Description: "List of the VPN authentication types for the Virtual Network Gateway.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						MaxItems:    3,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnAuthenticationType(), false),
						},
					},

					"vpn_client_protocols": {
						Description: "List of the protocols supported by the VPN client. Values `SSTP` and `IkeV2` are incompatible with the use of `aad_tenant`, `aad_audience` and `aad_issuer`.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(virtualnetworkgateways.PossibleValuesForVpnClientProtocol(), false),
						},
					},
				},
			},
		},

		"bgp_settings": {
			Description: "In this block the BGP specific settings can be defined.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"asn": {
						Description: "The Autonomous System Number (ASN) to use as part of the BGP.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
					},

					"peer_weight": {
						Description: "The weight added to routes which have been learned through BGP peering.",
						Type:        pluginsdk.TypeInt,
						Optional:    true,
						AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
					},

					// lintignore:XS003
					"peering_addresses": {
						Description: "Only one `peering_addresses` block can be specified except when `active_active` of this Virtual Network Gateway is `true`.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Computed:    true, // azignore:AZS007 - pre-existing violation
						MinItems:    1,
						MaxItems:    2, AtLeastOneOf: []string{
							"bgp_settings.0.asn",
							"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses",
						},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_configuration_name": {
									Description: "The name of the IP configuration of this Virtual Network Gateway. In case there are multiple `ip_configuration` blocks defined, this property is **required** to specify.",
									Type:        pluginsdk.TypeString,
									Optional:    true,
									// Note: O+C In case there is only one `ip_configuration` in root level. This property can be deduced from the that.
									Computed:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"apipa_addresses": {
									Description: "A list of Azure custom APIPA addresses assigned to the BGP peer of the Virtual Network Gateway.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validate.IPAddressInAzureReservedAPIPARange,
									},
								},
								"default_addresses": {
									Description: "A list of peering address assigned to the BGP peer of the Virtual Network Gateway.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
								"tunnel_ip_addresses": {
									Description: "A list of tunnel IP addresses assigned to the BGP peer of the Virtual Network Gateway.",
									Type:        pluginsdk.TypeList,
									Computed:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		// lintignore:XS003
		"custom_route": {
			Description: "A custom routes address space for a Virtual Network Gateway and a VpnClient.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"address_prefixes": {
						Description: "A list of address blocks reserved for this Virtual Network in CIDR notation.",
						Type:        pluginsdk.TypeSet,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"default_local_network_gateway_id": {
			Description:  "The ID of the Local Network Gateway through which outbound Internet traffic from the Virtual Network in which the gateway is created will be routed (*forced tunnelling*). Refer to the [Azure documentation on forced tunnelling](https://learn.microsoft.com/azure/vpn-gateway/vpn-gateway-forced-tunneling-rm). If not specified, forced tunnelling is disabled.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: localnetworkgateways.ValidateLocalNetworkGatewayID,
		},

		"bgp_route_translation_for_nat_enabled": {
			Description: "Whether BGP Route Translation for NAT is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"dns_forwarding_enabled": {
			Description: "Whether DNS forwarding is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"ip_sec_replay_protection_enabled": {
			Description: "Whether IP Sec Replay Protection is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"maximum_scale_unit": {
			Description:  "The maximum scale unit for the Virtual Network Gateway.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(1, 40),
			RequiredWith: []string{"maximum_scale_unit", "minimum_scale_unit"},
		},

		"minimum_scale_unit": {
			Description:  "The minimum scale unit for the Virtual Network Gateway.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.IntBetween(1, 40),
			RequiredWith: []string{"maximum_scale_unit", "minimum_scale_unit"},
		},

		"remote_vnet_traffic_enabled": {
			Description: "Whether remote VNet traffic that is used to configure this gateway to accept traffic from other Azure Virtual Networks is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"virtual_wan_traffic_enabled": {
			Description: "Whether remote VNet traffic that is used to configure this gateway to accept traffic from remote Virtual WAN networks is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"tags": commonschema.Tags(),
	}
}
