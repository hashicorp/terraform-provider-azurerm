// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/localnetworkgateways"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/virtualnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayCreate,
		Read:   resourceVirtualNetworkGatewayRead,
		Update: resourceVirtualNetworkGatewayUpdate,
		Delete: resourceVirtualNetworkGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: resourceVirtualNetworkGatewaySchema(),
	}
}

func resourceVirtualNetworkGatewaySchema() map[string]*pluginsdk.Schema {
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(virtualnetworkgateways.VpnTypeRouteBased),
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgateways.VpnTypeRouteBased),
				string(virtualnetworkgateways.VpnTypePolicyBased),
			}, false),
		},

		"edge_zone": commonschema.EdgeZoneOptionalForceNew(),

		// TODO 4.0: change this from enable_* to *_enabled
		"enable_bgp": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
		},

		"private_ip_address_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
		},

		"active_active": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Computed: true,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualnetworkgateways.VpnGatewayGenerationGenerationOne),
				string(virtualnetworkgateways.VpnGatewayGenerationGenerationTwo),
				string(virtualnetworkgateways.VpnGatewayGenerationNone),
			}, false),
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
						Type:     pluginsdk.TypeString,
						Optional: true,
						ValidateFunc: validation.StringInSlice([]string{
							string(virtualnetworkgateways.IPAllocationMethodStatic),
							string(virtualnetworkgateways.IPAllocationMethodDynamic),
						}, false),
						Default: string(virtualnetworkgateways.IPAllocationMethodDynamic),
					},

					"subnet_id": {
						Type:             pluginsdk.TypeString,
						Required:         true,
						ValidateFunc:     validate.IsGatewaySubnet,
						DiffSuppressFunc: suppress.CaseDifference,
					},

					"public_ip_address_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
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
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.VpnPolicyMemberAttributeTypeAADGroupId),
										string(virtualnetworkgateways.VpnPolicyMemberAttributeTypeCertificateGroupId),
										string(virtualnetworkgateways.VpnPolicyMemberAttributeTypeRadiusAzureGroupId),
									}, false),
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
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.DhGroupDHGroupOne),
										string(virtualnetworkgateways.DhGroupDHGroupOneFour),
										string(virtualnetworkgateways.DhGroupDHGroupTwo),
										string(virtualnetworkgateways.DhGroupDHGroupTwoZeroFourEight),
										string(virtualnetworkgateways.DhGroupDHGroupTwoFour),
										string(virtualnetworkgateways.DhGroupECPTwoFiveSix),
										string(virtualnetworkgateways.DhGroupECPThreeEightFour),
										string(virtualnetworkgateways.DhGroupNone),
									}, false),
								},

								"ike_encryption": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.IkeEncryptionAESOneTwoEight),
										string(virtualnetworkgateways.IkeEncryptionAESOneNineTwo),
										string(virtualnetworkgateways.IkeEncryptionAESTwoFiveSix),
										string(virtualnetworkgateways.IkeEncryptionDES),
										string(virtualnetworkgateways.IkeEncryptionDESThree),
										string(virtualnetworkgateways.IkeEncryptionGCMAESOneTwoEight),
										string(virtualnetworkgateways.IkeEncryptionGCMAESTwoFiveSix),
									}, false),
								},

								"ike_integrity": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.IkeIntegrityGCMAESOneTwoEight),
										string(virtualnetworkgateways.IkeIntegrityGCMAESTwoFiveSix),
										string(virtualnetworkgateways.IkeIntegrityMDFive),
										string(virtualnetworkgateways.IkeIntegritySHAOne),
										string(virtualnetworkgateways.IkeIntegritySHATwoFiveSix),
										string(virtualnetworkgateways.IkeIntegritySHAThreeEightFour),
									}, false),
								},

								"ipsec_encryption": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.IPsecEncryptionAESOneTwoEight),
										string(virtualnetworkgateways.IPsecEncryptionAESOneNineTwo),
										string(virtualnetworkgateways.IPsecEncryptionAESTwoFiveSix),
										string(virtualnetworkgateways.IPsecEncryptionDES),
										string(virtualnetworkgateways.IPsecEncryptionDESThree),
										string(virtualnetworkgateways.IPsecEncryptionGCMAESOneTwoEight),
										string(virtualnetworkgateways.IPsecEncryptionGCMAESOneNineTwo),
										string(virtualnetworkgateways.IPsecEncryptionGCMAESTwoFiveSix),
										string(virtualnetworkgateways.IPsecEncryptionNone),
									}, false),
								},

								"ipsec_integrity": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.IPsecIntegrityGCMAESOneTwoEight),
										string(virtualnetworkgateways.IPsecIntegrityGCMAESOneNineTwo),
										string(virtualnetworkgateways.IPsecIntegrityGCMAESTwoFiveSix),
										string(virtualnetworkgateways.IPsecIntegrityMDFive),
										string(virtualnetworkgateways.IPsecIntegritySHAOne),
										string(virtualnetworkgateways.IPsecIntegritySHATwoFiveSix),
									}, false),
								},

								"pfs_group": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(virtualnetworkgateways.PfsGroupECPTwoFiveSix),
										string(virtualnetworkgateways.PfsGroupECPThreeEightFour),
										string(virtualnetworkgateways.PfsGroupNone),
										string(virtualnetworkgateways.PfsGroupPFSOne),
										string(virtualnetworkgateways.PfsGroupPFSOneFour),
										string(virtualnetworkgateways.PfsGroupPFSTwo),
										string(virtualnetworkgateways.PfsGroupPFSTwoZeroFourEight),
										string(virtualnetworkgateways.PfsGroupPFSTwoFour),
										string(virtualnetworkgateways.PfsGroupPFSMM),
									}, false),
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
						RequiredWith: []string{"vpn_client_configuration.0.radius_server_address"},
					},

					"vpn_auth_types": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true,
						MaxItems: 3,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgateways.VpnAuthenticationTypeCertificate),
								string(virtualnetworkgateways.VpnAuthenticationTypeAAD),
								string(virtualnetworkgateways.VpnAuthenticationTypeRadius),
							}, false),
						},
					},

					"vpn_client_protocols": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								string(virtualnetworkgateways.VpnClientProtocolIkeVTwo),
								string(virtualnetworkgateways.VpnClientProtocolOpenVPN),
								string(virtualnetworkgateways.VpnClientProtocolSSTP),
							}, false),
						},
					},
				},
			},
		},

		"bgp_settings": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
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
						Computed: true,
						Optional: true,
						MinItems: 1,
						MaxItems: 2,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"ip_configuration_name": {
									Type: pluginsdk.TypeString,
									// In case there is only one `ip_configuration` in root level. This property can be deduced from the that.
									Optional:     true,
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

func resourceVirtualNetworkGatewayCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway creation.")

	id := virtualnetworkgateways.NewVirtualNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_network_gateway", id.ID())
	}

	properties, err := getVirtualNetworkGatewayProperties(id, d)
	if err != nil {
		return err
	}

	gateway := virtualnetworkgateways.VirtualNetworkGateway{
		Name:             pointer.To(id.VirtualNetworkGatewayName),
		ExtendedLocation: expandEdgeZoneModel(d.Get("edge_zone").(string)),
		Location:         pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:             tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties:       *properties,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, gateway); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayRead(d, meta)
}

func resourceVirtualNetworkGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkGatewayName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("edge_zone", flattenEdgeZoneModel(model.ExtendedLocation))

		props := model.Properties

		d.Set("type", string(pointer.From(props.GatewayType)))
		d.Set("enable_bgp", props.EnableBgp)
		d.Set("private_ip_address_enabled", props.EnablePrivateIPAddress)
		d.Set("active_active", props.ActiveActive)
		d.Set("bgp_route_translation_for_nat_enabled", props.EnableBgpRouteTranslationForNat)
		d.Set("dns_forwarding_enabled", props.EnableDnsForwarding)
		d.Set("ip_sec_replay_protection_enabled", !*props.DisableIPSecReplayProtection)
		d.Set("remote_vnet_traffic_enabled", props.AllowRemoteVnetTraffic)
		d.Set("virtual_wan_traffic_enabled", props.AllowVirtualWanTraffic)
		d.Set("generation", string(pointer.From(props.VpnGatewayGeneration)))

		if props.VpnType != nil {
			d.Set("vpn_type", string(pointer.From(props.VpnType)))
		}

		if props.GatewayDefaultSite != nil {
			d.Set("default_local_network_gateway_id", props.GatewayDefaultSite.Id)
		}

		if props.Sku != nil {
			d.Set("sku", string(pointer.From(props.Sku.Name)))
		}

		if err := d.Set("ip_configuration", flattenVirtualNetworkGatewayIPConfigurations(props.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("policy_group", flattenVirtualNetworkGatewayPolicyGroups(props.VirtualNetworkGatewayPolicyGroups)); err != nil {
			return fmt.Errorf("setting `policy_group`: %+v", err)
		}

		vpnClientConfig, err := flattenVirtualNetworkGatewayVpnClientConfig(props.VpnClientConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("vpn_client_configuration", vpnClientConfig); err != nil {
			return fmt.Errorf("setting `vpn_client_configuration`: %+v", err)
		}

		bgpSettings, err := flattenVirtualNetworkGatewayBgpSettings(props.BgpSettings)
		if err != nil {
			return err
		}
		if err := d.Set("bgp_settings", bgpSettings); err != nil {
			return fmt.Errorf("setting `bgp_settings`: %+v", err)
		}

		if err := d.Set("custom_route", flattenVirtualNetworkGatewayAddressSpace(props.CustomRoutes)); err != nil {
			return fmt.Errorf("setting `custom_route`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceVirtualNetworkGatewayUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway update.")

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id)
	if err != nil {
		return err
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("enable_bgp") {
		payload.Properties.EnableBgp = pointer.To(d.Get("enable_bgp").(bool))
	}

	if d.HasChange("active_active") {
		payload.Properties.ActiveActive = pointer.To(d.Get("active_active").(bool))
	}

	if d.HasChange("sku") {
		payload.Properties.Sku = expandVirtualNetworkGatewaySku(d)
	}

	if d.HasChange("ip_configuration") {
		payload.Properties.IPConfigurations = expandVirtualNetworkGatewayIPConfigurations(d)
	}

	if d.HasChange("policy_group") {
		payload.Properties.VirtualNetworkGatewayPolicyGroups = expandVirtualNetworkGatewayPolicyGroups(d.Get("policy_group").([]interface{}))
	}

	if d.HasChange("vpn_client_configuration") {
		payload.Properties.VpnClientConfiguration = expandVirtualNetworkGatewayVpnClientConfig(d, *id)
	}

	if d.HasChange("bgp_settings") {
		bgpSettings, err := expandVirtualNetworkGatewayBgpSettings(*id, d)
		if err != nil {
			return err
		}
		payload.Properties.BgpSettings = bgpSettings
	}

	if d.HasChange("custom_route") {
		payload.Properties.CustomRoutes = expandVirtualNetworkGatewayAddressSpace(d.Get("custom_route").([]interface{}))
	}

	if d.HasChange("default_local_network_gateway_id") {
		payload.Properties.GatewayDefaultSite = nil
		if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
			payload.Properties.GatewayDefaultSite = &virtualnetworkgateways.SubResource{
				Id: &gatewayDefaultSiteID,
			}
		}
	}

	if d.HasChange("bgp_route_translation_for_nat_enabled") {
		payload.Properties.EnableBgpRouteTranslationForNat = pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool))
	}

	if d.HasChange("dns_forwarding_enabled") {
		payload.Properties.EnableDnsForwarding = pointer.To(d.Get("dns_forwarding_enabled").(bool))
	}

	if d.HasChange("ip_sec_replay_protection_enabled") {
		payload.Properties.DisableIPSecReplayProtection = pointer.To(!d.Get("ip_sec_replay_protection_enabled").(bool))
	}

	if d.HasChange("remote_vnet_traffic_enabled") {
		payload.Properties.AllowRemoteVnetTraffic = pointer.To(d.Get("remote_vnet_traffic_enabled").(bool))
	}

	if d.HasChange("virtual_wan_traffic_enabled") {
		payload.Properties.AllowVirtualWanTraffic = pointer.To(d.Get("virtual_wan_traffic_enabled").(bool))
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayRead(d, meta)
}

func resourceVirtualNetworkGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualNetworkGateways
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualnetworkgateways.ParseVirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func getVirtualNetworkGatewayProperties(id virtualnetworkgateways.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*virtualnetworkgateways.VirtualNetworkGatewayPropertiesFormat, error) {
	props := &virtualnetworkgateways.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:                     pointer.To(virtualnetworkgateways.VirtualNetworkGatewayType(d.Get("type").(string))),
		VpnType:                         pointer.To(virtualnetworkgateways.VpnType(d.Get("vpn_type").(string))),
		EnableBgp:                       pointer.To(d.Get("enable_bgp").(bool)),
		EnablePrivateIPAddress:          pointer.To(d.Get("private_ip_address_enabled").(bool)),
		ActiveActive:                    pointer.To(d.Get("active_active").(bool)),
		EnableBgpRouteTranslationForNat: pointer.To(d.Get("bgp_route_translation_for_nat_enabled").(bool)),
		DisableIPSecReplayProtection:    pointer.To(!d.Get("ip_sec_replay_protection_enabled").(bool)),
		AllowRemoteVnetTraffic:          pointer.To(d.Get("remote_vnet_traffic_enabled").(bool)),
		AllowVirtualWanTraffic:          pointer.To(d.Get("virtual_wan_traffic_enabled").(bool)),
		Sku:                             expandVirtualNetworkGatewaySku(d),
		IPConfigurations:                expandVirtualNetworkGatewayIPConfigurations(d),
		CustomRoutes:                    expandVirtualNetworkGatewayAddressSpace(d.Get("custom_route").([]interface{})),
	}

	if v, ok := d.GetOk("generation"); ok {
		props.VpnGatewayGeneration = pointer.To(virtualnetworkgateways.VpnGatewayGeneration(v.(string)))
	}

	if v, ok := d.GetOk("dns_forwarding_enabled"); ok {
		props.EnableDnsForwarding = pointer.To(v.(bool))
	}

	if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
		props.GatewayDefaultSite = &virtualnetworkgateways.SubResource{
			Id: &gatewayDefaultSiteID,
		}
	}

	if v, ok := d.GetOk("policy_group"); ok {
		props.VirtualNetworkGatewayPolicyGroups = expandVirtualNetworkGatewayPolicyGroups(v.([]interface{}))
	}

	if _, ok := d.GetOk("vpn_client_configuration"); ok {
		props.VpnClientConfiguration = expandVirtualNetworkGatewayVpnClientConfig(d, id)
	}

	if _, ok := d.GetOk("bgp_settings"); ok {
		bgpSettings, err := expandVirtualNetworkGatewayBgpSettings(id, d)
		if err != nil {
			return nil, err
		}
		props.BgpSettings = bgpSettings
	}

	gatewayType := pointer.From(props.GatewayType)
	vpnType := pointer.From(props.VpnType)
	vpnGatewayGeneration := pointer.From(props.VpnGatewayGeneration)
	skuName := string(pointer.From(props.Sku.Name))

	// Sku validation for policy-based VPN gateways
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypePolicyBased {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayPolicyBasedVpnSku()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of first geneneration
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypeRouteBased && vpnGatewayGeneration == virtualnetworkgateways.VpnGatewayGenerationGenerationOne {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of second geneneration
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeVpn && vpnType == virtualnetworkgateways.VpnTypeRouteBased && vpnGatewayGeneration == virtualnetworkgateways.VpnGatewayGenerationGenerationTwo {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2()); !ok {
			return nil, err
		}
	}

	// Sku validation for ExpressRoute gateways
	if gatewayType == virtualnetworkgateways.VirtualNetworkGatewayTypeExpressRoute {
		if ok, err := evaluateSchemaValidateFunc(skuName, "sku", validateVirtualNetworkGatewayExpressRouteSku()); !ok {
			return nil, err
		}
	}

	return props, nil
}

func expandVirtualNetworkGatewayBgpSettings(id virtualnetworkgateways.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*virtualnetworkgateways.BgpSettings, error) {
	bgpSets := d.Get("bgp_settings").([]interface{})
	if len(bgpSets) == 0 {
		return nil, nil
	}

	bgp := bgpSets[0].(map[string]interface{})

	peeringAddresses, err := expandVirtualNetworkGatewayBgpPeeringAddresses(id, d.Get("ip_configuration").([]interface{}), bgp["peering_addresses"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &virtualnetworkgateways.BgpSettings{
		Asn:                 pointer.To(int64(bgp["asn"].(int))),
		PeerWeight:          pointer.To(int64(bgp["peer_weight"].(int))),
		BgpPeeringAddresses: peeringAddresses,
	}, nil
}

func expandVirtualNetworkGatewayBgpPeeringAddresses(id virtualnetworkgateways.VirtualNetworkGatewayId, ipConfig, input []interface{}) (*[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]virtualnetworkgateways.IPConfigurationBgpPeeringAddress, 0)

	var existIpConfigName string
	if len(ipConfig) == 1 {
		existIpConfigName = ipConfig[0].(map[string]interface{})["name"].(string)
	}

	for _, e := range input {
		if e == nil {
			continue
		}
		b := e.(map[string]interface{})

		ipConfigName := b["ip_configuration_name"].(string)

		if ipConfigName == "" {
			// existIpConfigName is empty means there are more than one `ip_configuration` blocks, in which case users have to specify the
			// `ip_configuration_name` used in current `peering_addresses` setting.
			if existIpConfigName == "" {
				return nil, fmt.Errorf("`ip_configuration_name` has to be set in current `peering_addresses` block in case there are multiple `ip_configuration` blocks")
			}

			ipConfigName = existIpConfigName
		}

		// If there is an existing ip configuration name defined in the only `ip_configuration` block, and users explicitly set the `ip_configuration_name` in current
		// `peering_addresses` block, they should be the same name.
		if existIpConfigName != "" && ipConfigName != "" {
			if ipConfigName != existIpConfigName {
				return nil, fmt.Errorf("`ip_configuration.0.name` is not the same as `bgp_settings.0.peering_addresses.*.ip_configuration_name`")
			}
		}

		ipConfigId := parse.NewVirtualNetworkGatewayIpConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkGatewayName, ipConfigName)
		result = append(result, virtualnetworkgateways.IPConfigurationBgpPeeringAddress{
			IPconfigurationId:    pointer.To(ipConfigId.ID()),
			CustomBgpIPAddresses: utils.ExpandStringSlice(b["apipa_addresses"].([]interface{})),
		})
	}

	return &result, nil
}

func expandVirtualNetworkGatewayIPConfigurations(d *pluginsdk.ResourceData) *[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration, 0, len(configs))

	for _, c := range configs {
		conf := c.(map[string]interface{})

		name := conf["name"].(string)

		props := &virtualnetworkgateways.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: pointer.To(virtualnetworkgateways.IPAllocationMethod(conf["private_ip_address_allocation"].(string))),
		}

		if subnetID := conf["subnet_id"].(string); subnetID != "" {
			props.Subnet = &virtualnetworkgateways.SubResource{
				Id: &subnetID,
			}
		}

		if publicIP := conf["public_ip_address_id"].(string); publicIP != "" {
			props.PublicIPAddress = &virtualnetworkgateways.SubResource{
				Id: &publicIP,
			}
		}

		ipConfig := virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration{
			Name:       &name,
			Properties: props,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return &ipConfigs
}

func expandVirtualNetworkGatewayVpnClientConfig(d *pluginsdk.ResourceData, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *virtualnetworkgateways.VpnClientConfiguration {
	configSets := d.Get("vpn_client_configuration").([]interface{})
	if len(configSets) == 0 {
		// return nil will delete the existing vpn client configuration
		return nil
	}

	conf := configSets[0].(map[string]interface{})

	confAddresses := conf["address_space"].([]interface{})
	addresses := make([]string, 0, len(confAddresses))
	for _, addr := range confAddresses {
		addresses = append(addresses, addr.(string))
	}

	rootCertsConf := conf["root_certificate"].(*pluginsdk.Set).List()
	rootCerts := make([]virtualnetworkgateways.VpnClientRootCertificate, 0, len(rootCertsConf))
	for _, rootCertSet := range rootCertsConf {
		rootCert := rootCertSet.(map[string]interface{})
		r := virtualnetworkgateways.VpnClientRootCertificate{
			Name: pointer.To(rootCert["name"].(string)),
			Properties: virtualnetworkgateways.VpnClientRootCertificatePropertiesFormat{
				PublicCertData: rootCert["public_cert_data"].(string),
			},
		}
		rootCerts = append(rootCerts, r)
	}

	revokedCertsConf := conf["revoked_certificate"].(*pluginsdk.Set).List()
	revokedCerts := make([]virtualnetworkgateways.VpnClientRevokedCertificate, 0, len(revokedCertsConf))
	for _, revokedCertSet := range revokedCertsConf {
		revokedCert := revokedCertSet.(map[string]interface{})
		r := virtualnetworkgateways.VpnClientRevokedCertificate{
			Name: pointer.To(revokedCert["name"].(string)),
			Properties: &virtualnetworkgateways.VpnClientRevokedCertificatePropertiesFormat{
				Thumbprint: pointer.To(revokedCert["thumbprint"].(string)),
			},
		}
		revokedCerts = append(revokedCerts, r)
	}

	vpnClientProtocolsConf := conf["vpn_client_protocols"].(*pluginsdk.Set).List()
	vpnClientProtocols := make([]virtualnetworkgateways.VpnClientProtocol, 0, len(vpnClientProtocolsConf))
	for _, vpnClientProtocol := range vpnClientProtocolsConf {
		p := virtualnetworkgateways.VpnClientProtocol(vpnClientProtocol.(string))
		vpnClientProtocols = append(vpnClientProtocols, p)
	}

	vpnAuthTypesConf := conf["vpn_auth_types"].(*pluginsdk.Set).List()
	vpnAuthTypes := make([]virtualnetworkgateways.VpnAuthenticationType, 0, len(vpnAuthTypesConf))
	for _, vpnAuthType := range vpnAuthTypesConf {
		a := virtualnetworkgateways.VpnAuthenticationType(vpnAuthType.(string))
		vpnAuthTypes = append(vpnAuthTypes, a)
	}

	return &virtualnetworkgateways.VpnClientConfiguration{
		VpnClientAddressPool: &virtualnetworkgateways.AddressSpace{
			AddressPrefixes: &addresses,
		},
		AadTenant:                         pointer.To(conf["aad_tenant"].(string)),
		AadAudience:                       pointer.To(conf["aad_audience"].(string)),
		AadIssuer:                         pointer.To(conf["aad_issuer"].(string)),
		VngClientConnectionConfigurations: expandVirtualNetworkGatewayClientConnections(conf["virtual_network_gateway_client_connection"].([]interface{}), vnetGatewayId),
		VpnClientIPsecPolicies:            expandVirtualNetworkGatewayIpsecPolicies(conf["ipsec_policy"].([]interface{})),
		VpnClientRootCertificates:         &rootCerts,
		VpnClientRevokedCertificates:      &revokedCerts,
		VpnClientProtocols:                &vpnClientProtocols,
		RadiusServers:                     expandVirtualNetworkGatewayRadiusServers(conf["radius_server"].([]interface{})),
		RadiusServerAddress:               pointer.To(conf["radius_server_address"].(string)),
		RadiusServerSecret:                pointer.To(conf["radius_server_secret"].(string)),
		VpnAuthenticationTypes:            &vpnAuthTypes,
	}
}

func expandVirtualNetworkGatewaySku(d *pluginsdk.ResourceData) *virtualnetworkgateways.VirtualNetworkGatewaySku {
	sku := d.Get("sku").(string)

	return &virtualnetworkgateways.VirtualNetworkGatewaySku{
		Name: pointer.To(virtualnetworkgateways.VirtualNetworkGatewaySkuName(sku)),
		Tier: pointer.To(virtualnetworkgateways.VirtualNetworkGatewaySkuTier(sku)),
	}
}

func expandVirtualNetworkGatewayAddressSpace(input []interface{}) *virtualnetworkgateways.AddressSpace {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &virtualnetworkgateways.AddressSpace{
		AddressPrefixes: utils.ExpandStringSlice(v["address_prefixes"].(*pluginsdk.Set).List()),
	}
}

func expandVirtualNetworkGatewayIpsecPolicies(input []interface{}) *[]virtualnetworkgateways.IPsecPolicy {
	results := make([]virtualnetworkgateways.IPsecPolicy, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, virtualnetworkgateways.IPsecPolicy{
			DhGroup:             virtualnetworkgateways.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       virtualnetworkgateways.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        virtualnetworkgateways.IkeIntegrity(v["ike_integrity"].(string)),
			IPsecEncryption:     virtualnetworkgateways.IPsecEncryption(v["ipsec_encryption"].(string)),
			IPsecIntegrity:      virtualnetworkgateways.IPsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            virtualnetworkgateways.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   int64(v["sa_lifetime_in_seconds"].(int)),
			SaDataSizeKilobytes: int64(v["sa_data_size_in_kilobytes"].(int)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayRadiusServers(input []interface{}) *[]virtualnetworkgateways.RadiusServer {
	results := make([]virtualnetworkgateways.RadiusServer, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, virtualnetworkgateways.RadiusServer{
			RadiusServerAddress: v["address"].(string),
			RadiusServerScore:   pointer.To(int64(v["score"].(int))),
			RadiusServerSecret:  pointer.To(v["secret"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyGroups(input []interface{}) *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup {
	results := make([]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyGroup := item.(map[string]interface{})

		results = append(results, virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup{
			Name: pointer.To(policyGroup["name"].(string)),
			Properties: &virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupProperties{
				IsDefault:     policyGroup["is_default"].(bool),
				PolicyMembers: *expandVirtualNetworkGatewayPolicyMembers(policyGroup["policy_member"].([]interface{})),
				Priority:      int64(policyGroup["priority"].(int)),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyMembers(input []interface{}) *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember {
	results := make([]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyMember := item.(map[string]interface{})

		results = append(results, virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember{
			Name:           pointer.To(policyMember["name"].(string)),
			AttributeType:  pointer.To(virtualnetworkgateways.VpnPolicyMemberAttributeType(policyMember["type"].(string))),
			AttributeValue: pointer.To(policyMember["value"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayClientConnections(input []interface{}, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *[]virtualnetworkgateways.VngClientConnectionConfiguration {
	results := make([]virtualnetworkgateways.VngClientConnectionConfiguration, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		vngClientConnectionConfiguration := item.(map[string]interface{})

		results = append(results, virtualnetworkgateways.VngClientConnectionConfiguration{
			Name: pointer.To(vngClientConnectionConfiguration["name"].(string)),
			Properties: &virtualnetworkgateways.VngClientConnectionConfigurationProperties{
				VpnClientAddressPool:              *expandVirtualNetworkGatewayAddressPool(vngClientConnectionConfiguration["address_prefixes"].([]interface{})),
				VirtualNetworkGatewayPolicyGroups: *expandVirtualNetworkGatewayPolicyGroupNames(vngClientConnectionConfiguration["policy_group_names"].([]interface{}), vnetGatewayId),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayAddressPool(input []interface{}) *virtualnetworkgateways.AddressSpace {
	if len(input) == 0 {
		return &virtualnetworkgateways.AddressSpace{}
	}

	addressPrefixes := make([]string, 0)
	for _, v := range input {
		addressPrefixes = append(addressPrefixes, v.(string))
	}

	return &virtualnetworkgateways.AddressSpace{
		AddressPrefixes: pointer.To(addressPrefixes),
	}
}

func expandVirtualNetworkGatewayPolicyGroupNames(input []interface{}, vnetGatewayId virtualnetworkgateways.VirtualNetworkGatewayId) *[]virtualnetworkgateways.SubResource {
	results := make([]virtualnetworkgateways.SubResource, 0)
	if len(input) == 0 {
		return &results
	}

	for _, v := range input {
		policyGroupId := parse.NewVirtualNetworkGatewayPolicyGroupID(vnetGatewayId.SubscriptionId, vnetGatewayId.ResourceGroupName, vnetGatewayId.VirtualNetworkGatewayName, v.(string))
		result := virtualnetworkgateways.SubResource{
			Id: pointer.To(policyGroupId.ID()),
		}
		results = append(results, result)
	}

	return &results
}

func flattenVirtualNetworkGatewayBgpSettings(settings *virtualnetworkgateways.BgpSettings) ([]interface{}, error) {
	output := make([]interface{}, 0)

	if settings != nil {
		flat := make(map[string]interface{})

		if asn := settings.Asn; asn != nil {
			flat["asn"] = int(*asn)
		}
		if weight := settings.PeerWeight; weight != nil {
			flat["peer_weight"] = int(*weight)
		}

		var err error
		flat["peering_addresses"], err = flattenVirtualNetworkGatewayBgpPeeringAddresses(settings.BgpPeeringAddresses)
		if err != nil {
			return nil, err
		}

		output = append(output, flat)
	}

	return output, nil
}

func flattenVirtualNetworkGatewayBgpPeeringAddresses(input *[]virtualnetworkgateways.IPConfigurationBgpPeeringAddress) (interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var ipConfigName string
		if e.IPconfigurationId != nil {
			id, err := parse.VirtualNetworkGatewayIpConfigurationIDInsensitively(*e.IPconfigurationId)
			if err != nil {
				return nil, err
			}
			ipConfigName = id.IpConfigurationName
		}

		output = append(output, map[string]interface{}{
			"ip_configuration_name": ipConfigName,
			"apipa_addresses":       utils.FlattenStringSlice(e.CustomBgpIPAddresses),
			"default_addresses":     utils.FlattenStringSlice(e.DefaultBgpIPAddresses),
			"tunnel_ip_addresses":   utils.FlattenStringSlice(e.TunnelIPAddresses),
		})
	}

	return output, nil
}

func flattenVirtualNetworkGatewayIPConfigurations(ipConfigs *[]virtualnetworkgateways.VirtualNetworkGatewayIPConfiguration) []interface{} {
	flat := make([]interface{}, 0)

	if ipConfigs != nil {
		for _, cfg := range *ipConfigs {
			props := cfg.Properties
			v := make(map[string]interface{})

			if name := cfg.Name; name != nil {
				v["name"] = *name
			}
			v["private_ip_address_allocation"] = string(pointer.From(props.PrivateIPAllocationMethod))

			if subnet := props.Subnet; subnet != nil {
				if id := subnet.Id; id != nil {
					v["subnet_id"] = *id
				}
			}

			if pip := props.PublicIPAddress; pip != nil {
				if id := pip.Id; id != nil {
					v["public_ip_address_id"] = *id
				}
			}

			flat = append(flat, v)
		}
	}

	return flat
}

func flattenVirtualNetworkGatewayVpnClientConfig(cfg *virtualnetworkgateways.VpnClientConfiguration) ([]interface{}, error) {
	if cfg == nil {
		return []interface{}{}, nil
	}
	flat := map[string]interface{}{
		"ipsec_policy":  flattenVirtualNetworkGatewayIPSecPolicies(cfg.VpnClientIPsecPolicies),
		"radius_server": flattenVirtualNetworkGatewayRadiusServers(cfg.RadiusServers),
	}

	connection, err := flattenVirtualNetworkGatewayClientConnections(cfg.VngClientConnectionConfigurations)
	if err != nil {
		return nil, err
	}
	flat["virtual_network_gateway_client_connection"] = connection

	if pool := cfg.VpnClientAddressPool; pool != nil {
		flat["address_space"] = utils.FlattenStringSlice(pool.AddressPrefixes)
	} else {
		flat["address_space"] = []interface{}{}
	}

	if v := cfg.AadTenant; v != nil {
		flat["aad_tenant"] = *v
	}

	if v := cfg.AadAudience; v != nil {
		flat["aad_audience"] = *v
	}

	if v := cfg.AadIssuer; v != nil {
		flat["aad_issuer"] = *v
	}

	rootCerts := make([]interface{}, 0)
	if certs := cfg.VpnClientRootCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]interface{}{
				"name":             *cert.Name,
				"public_cert_data": cert.Properties.PublicCertData,
			}
			rootCerts = append(rootCerts, v)
		}
	}
	flat["root_certificate"] = pluginsdk.NewSet(hashVirtualNetworkGatewayRootCert, rootCerts)

	revokedCerts := make([]interface{}, 0)
	if certs := cfg.VpnClientRevokedCertificates; certs != nil {
		for _, cert := range *certs {
			v := map[string]interface{}{
				"name":       *cert.Name,
				"thumbprint": *cert.Properties.Thumbprint,
			}
			revokedCerts = append(revokedCerts, v)
		}
	}
	flat["revoked_certificate"] = pluginsdk.NewSet(hashVirtualNetworkGatewayRevokedCert, revokedCerts)

	vpnClientProtocols := &pluginsdk.Set{F: pluginsdk.HashString}
	if vpnProtocols := cfg.VpnClientProtocols; vpnProtocols != nil {
		for _, protocol := range *vpnProtocols {
			vpnClientProtocols.Add(string(protocol))
		}
	}
	flat["vpn_client_protocols"] = vpnClientProtocols

	vpnAuthTypes := &pluginsdk.Set{F: pluginsdk.HashString}
	if authTypes := cfg.VpnAuthenticationTypes; authTypes != nil {
		for _, authType := range *authTypes {
			vpnAuthTypes.Add(string(authType))
		}
	}
	flat["vpn_auth_types"] = vpnAuthTypes

	if v := cfg.RadiusServerAddress; v != nil {
		flat["radius_server_address"] = *v
	}

	if v := cfg.RadiusServerSecret; v != nil {
		flat["radius_server_secret"] = *v
	}

	return []interface{}{flat}, nil
}

func hashVirtualNetworkGatewayRootCert(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["public_cert_data"].(string)))

	return pluginsdk.HashString(buf.String())
}

func hashVirtualNetworkGatewayRevokedCert(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["thumbprint"].(string)))

	return pluginsdk.HashString(buf.String())
}

func validateVirtualNetworkGatewayPolicyBasedVpnSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierBasic),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierBasic),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierStandard),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierHighPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwOne),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwo),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThree),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwOneAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThreeAZ),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwo),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThree),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFour),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFive),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwThreeAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFourAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameVpnGwFiveAZ),
	}, false)
}

func validateVirtualNetworkGatewayExpressRouteSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierStandard),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierHighPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuTierUltraPerformance),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwOneAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwTwoAZ),
		string(virtualnetworkgateways.VirtualNetworkGatewaySkuNameErGwThreeAZ),
	}, false)
}

func flattenVirtualNetworkGatewayAddressSpace(input *virtualnetworkgateways.AddressSpace) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"address_prefixes": utils.FlattenStringSlice(input.AddressPrefixes),
		},
	}
}

func flattenVirtualNetworkGatewayRadiusServers(input *[]virtualnetworkgateways.RadiusServer) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"address": item.RadiusServerAddress,
			"secret":  pointer.From(item.RadiusServerSecret),
			"score":   pointer.From(item.RadiusServerScore),
		})
	}

	return results
}

func flattenVirtualNetworkGatewayIPSecPolicies(input *[]virtualnetworkgateways.IPsecPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"dh_group":                  string(item.DhGroup),
			"ipsec_encryption":          string(item.IPsecEncryption),
			"ipsec_integrity":           string(item.IPsecIntegrity),
			"ike_encryption":            string(item.IkeEncryption),
			"ike_integrity":             string(item.IkeIntegrity),
			"pfs_group":                 string(item.PfsGroup),
			"sa_data_size_in_kilobytes": item.SaDataSizeKilobytes,
			"sa_lifetime_in_seconds":    item.SaLifeTimeSeconds,
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicyGroups(input *[]virtualnetworkgateways.VirtualNetworkGatewayPolicyGroup) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		var isDefault bool
		var policyMember interface{}
		var priority int64
		if props := item.Properties; props != nil {
			isDefault = props.IsDefault
			policyMember = flattenVirtualNetworkGatewayPolicy(props.PolicyMembers)
			priority = props.Priority
		}
		results = append(results, map[string]interface{}{
			"name":          pointer.From(item.Name),
			"is_default":    isDefault,
			"policy_member": policyMember,
			"priority":      priority,
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicy(input []virtualnetworkgateways.VirtualNetworkGatewayPolicyGroupMember) []interface{} {
	results := make([]interface{}, 0)
	if len(input) == 0 {
		return results
	}

	for _, item := range input {
		results = append(results, map[string]interface{}{
			"name":  pointer.From(item.Name),
			"type":  string(pointer.From(item.AttributeType)),
			"value": pointer.From(item.AttributeValue),
		})
	}
	return results
}

func flattenVirtualNetworkGatewayClientConnections(input *[]virtualnetworkgateways.VngClientConnectionConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, item := range *input {
		var err error
		addressPrefixes := make([]string, 0)
		policyGroupNames := make([]string, 0)
		if props := item.Properties; props != nil {
			addressPrefixes = pointer.From(props.VpnClientAddressPool.AddressPrefixes)
			policyGroupNames, err = flattenVirtualNetworkGatewayPolicyGroupNames(props.VirtualNetworkGatewayPolicyGroups)
			if err != nil {
				return nil, err
			}
		}
		result := map[string]interface{}{
			"name":             pointer.From(item.Name),
			"address_prefixes": addressPrefixes,
		}

		result["policy_group_names"] = policyGroupNames

		results = append(results, result)
	}

	return results, nil
}

func flattenVirtualNetworkGatewayPolicyGroupNames(input []virtualnetworkgateways.SubResource) ([]string, error) {
	results := make([]string, 0)
	if len(input) == 0 {
		return results, nil
	}

	for _, item := range input {
		policyGroupId, err := parse.VirtualNetworkGatewayPolicyGroupID(*item.Id)
		if err != nil {
			return nil, err
		}
		results = append(results, policyGroupId.Name)
	}

	return results, nil
}
