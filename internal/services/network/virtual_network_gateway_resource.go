// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

func resourceVirtualNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayCreateUpdate,
		Read:   resourceVirtualNetworkGatewayRead,
		Update: resourceVirtualNetworkGatewayCreateUpdate,
		Delete: resourceVirtualNetworkGatewayDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.VirtualNetworkGatewayID(id)
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
				string(network.VirtualNetworkGatewayTypeExpressRoute),
				string(network.VirtualNetworkGatewayTypeVpn),
			}, false),
		},

		"vpn_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(network.VpnTypeRouteBased),
			ValidateFunc: validation.StringInSlice([]string{
				string(network.VpnTypeRouteBased),
				string(network.VpnTypePolicyBased),
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
				string(network.VpnGatewayGenerationGeneration1),
				string(network.VpnGatewayGenerationGeneration2),
				string(network.VpnGatewayGenerationNone),
			}, false),
		},

		"ip_configuration": {
			Type:     pluginsdk.TypeList,
			Required: true,
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
							string(network.IPAllocationMethodStatic),
							string(network.IPAllocationMethodDynamic),
						}, false),
						Default: string(network.IPAllocationMethodDynamic),
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
						ValidateFunc: validate.PublicIpAddressID,
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
										string(network.VpnPolicyMemberAttributeTypeAADGroupID),
										string(network.VpnPolicyMemberAttributeTypeCertificateGroupID),
										string(network.VpnPolicyMemberAttributeTypeRadiusAzureGroupID),
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
										string(network.DhGroupDHGroup1),
										string(network.DhGroupDHGroup14),
										string(network.DhGroupDHGroup2),
										string(network.DhGroupDHGroup2048),
										string(network.DhGroupDHGroup24),
										string(network.DhGroupECP256),
										string(network.DhGroupECP384),
										string(network.DhGroupNone),
									}, false),
								},

								"ike_encryption": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.IkeEncryptionAES128),
										string(network.IkeEncryptionAES192),
										string(network.IkeEncryptionAES256),
										string(network.IkeEncryptionDES),
										string(network.IkeEncryptionDES3),
										string(network.IkeEncryptionGCMAES128),
										string(network.IkeEncryptionGCMAES256),
									}, false),
								},

								"ike_integrity": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.IkeIntegrityGCMAES128),
										string(network.IkeIntegrityGCMAES256),
										string(network.IkeIntegrityMD5),
										string(network.IkeIntegritySHA1),
										string(network.IkeIntegritySHA256),
										string(network.IkeIntegritySHA384),
									}, false),
								},

								"ipsec_encryption": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.IpsecEncryptionAES128),
										string(network.IpsecEncryptionAES192),
										string(network.IpsecEncryptionAES256),
										string(network.IpsecEncryptionDES),
										string(network.IpsecEncryptionDES3),
										string(network.IpsecEncryptionGCMAES128),
										string(network.IpsecEncryptionGCMAES192),
										string(network.IpsecEncryptionGCMAES256),
										string(network.IpsecEncryptionNone),
									}, false),
								},

								"ipsec_integrity": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.IpsecIntegrityGCMAES128),
										string(network.IpsecIntegrityGCMAES192),
										string(network.IpsecIntegrityGCMAES256),
										string(network.IpsecIntegrityMD5),
										string(network.IpsecIntegritySHA1),
										string(network.IpsecIntegritySHA256),
									}, false),
								},

								"pfs_group": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										string(network.PfsGroupECP256),
										string(network.PfsGroupECP384),
										string(network.PfsGroupNone),
										string(network.PfsGroupPFS1),
										string(network.PfsGroupPFS14),
										string(network.PfsGroupPFS2),
										string(network.PfsGroupPFS2048),
										string(network.PfsGroupPFS24),
										string(network.PfsGroupPFSMM),
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
									ValidateFunc: validation.IntBetween(1024, 2147483647),
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
								string(network.VpnAuthenticationTypeCertificate),
								string(network.VpnAuthenticationTypeAAD),
								string(network.VpnAuthenticationTypeRadius),
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
								string(network.VpnClientProtocolIkeV2),
								string(network.VpnClientProtocolOpenVPN),
								string(network.VpnClientProtocolSSTP),
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
			ValidateFunc: validate.LocalNetworkGatewayID,
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

		"tags": tags.Schema(),
	}
}

func resourceVirtualNetworkGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway creation.")

	id := parse.NewVirtualNetworkGatewayID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	var existingVNetGateway network.VirtualNetworkGateway
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway", id.ID())
		}
	} else {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return err
		}
		existingVNetGateway = existing
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties, err := getVirtualNetworkGatewayProperties(id, d, existingVNetGateway)
	if err != nil {
		return err
	}

	gateway := network.VirtualNetworkGateway{
		Name:                                  &id.Name,
		ExtendedLocation:                      expandEdgeZone(d.Get("edge_zone").(string)),
		Location:                              &location,
		Tags:                                  tags.Expand(t),
		VirtualNetworkGatewayPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, gateway)
	if err != nil {
		return fmt.Errorf("Creating/Updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayRead(d, meta)
}

func resourceVirtualNetworkGatewayRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	d.Set("edge_zone", flattenEdgeZone(resp.ExtendedLocation))

	if gw := resp.VirtualNetworkGatewayPropertiesFormat; gw != nil {
		d.Set("type", string(gw.GatewayType))
		d.Set("enable_bgp", gw.EnableBgp)
		d.Set("private_ip_address_enabled", gw.EnablePrivateIPAddress)
		d.Set("active_active", gw.ActiveActive)
		d.Set("bgp_route_translation_for_nat_enabled", gw.EnableBgpRouteTranslationForNat)
		d.Set("dns_forwarding_enabled", gw.EnableDNSForwarding)
		d.Set("ip_sec_replay_protection_enabled", !*gw.DisableIPSecReplayProtection)
		d.Set("remote_vnet_traffic_enabled", gw.AllowRemoteVnetTraffic)
		d.Set("virtual_wan_traffic_enabled", gw.AllowVirtualWanTraffic)
		d.Set("generation", string(gw.VpnGatewayGeneration))

		if string(gw.VpnType) != "" {
			d.Set("vpn_type", string(gw.VpnType))
		}

		if gw.GatewayDefaultSite != nil {
			d.Set("default_local_network_gateway_id", gw.GatewayDefaultSite.ID)
		}

		if gw.Sku != nil {
			d.Set("sku", string(gw.Sku.Name))
		}

		if err := d.Set("ip_configuration", flattenVirtualNetworkGatewayIPConfigurations(gw.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("policy_group", flattenVirtualNetworkGatewayPolicyGroups(gw.VirtualNetworkGatewayPolicyGroups)); err != nil {
			return fmt.Errorf("setting `policy_group`: %+v", err)
		}

		vpnClientConfig, err := flattenVirtualNetworkGatewayVpnClientConfig(gw.VpnClientConfiguration)
		if err != nil {
			return err
		}
		if err := d.Set("vpn_client_configuration", vpnClientConfig); err != nil {
			return fmt.Errorf("setting `vpn_client_configuration`: %+v", err)
		}

		bgpSettings, err := flattenVirtualNetworkGatewayBgpSettings(gw.BgpSettings)
		if err != nil {
			return err
		}
		if err := d.Set("bgp_settings", bgpSettings); err != nil {
			return fmt.Errorf("setting `bgp_settings`: %+v", err)
		}

		if err := d.Set("custom_route", flattenVirtualNetworkGatewayAddressSpace(gw.CustomRoutes)); err != nil {
			return fmt.Errorf("setting `custom_route`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualNetworkGatewayDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return nil
}

func getVirtualNetworkGatewayProperties(id parse.VirtualNetworkGatewayId, d *pluginsdk.ResourceData, existingVNetGateway network.VirtualNetworkGateway) (*network.VirtualNetworkGatewayPropertiesFormat, error) {
	gatewayType := network.VirtualNetworkGatewayType(d.Get("type").(string))
	vpnType := network.VpnType(d.Get("vpn_type").(string))
	enableBgp := d.Get("enable_bgp").(bool)
	enablePrivateIpAddress := d.Get("private_ip_address_enabled").(bool)
	activeActive := d.Get("active_active").(bool)
	generation := network.VpnGatewayGeneration(d.Get("generation").(string))
	customRoute := d.Get("custom_route").([]interface{})
	bgpRouteTranslationForNatEnabled := d.Get("bgp_route_translation_for_nat_enabled").(bool)
	ipSecReplayProtectionEnabled := d.Get("ip_sec_replay_protection_enabled").(bool)
	remoteVnetTrafficEnabled := d.Get("remote_vnet_traffic_enabled").(bool)
	virtualWanTrafficEnabled := d.Get("virtual_wan_traffic_enabled").(bool)

	props := &network.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:                     gatewayType,
		VpnType:                         vpnType,
		EnableBgp:                       &enableBgp,
		EnablePrivateIPAddress:          &enablePrivateIpAddress,
		ActiveActive:                    &activeActive,
		VpnGatewayGeneration:            generation,
		EnableBgpRouteTranslationForNat: utils.Bool(bgpRouteTranslationForNatEnabled),
		DisableIPSecReplayProtection:    utils.Bool(!ipSecReplayProtectionEnabled),
		AllowRemoteVnetTraffic:          utils.Bool(remoteVnetTrafficEnabled),
		AllowVirtualWanTraffic:          utils.Bool(virtualWanTrafficEnabled),
		Sku:                             expandVirtualNetworkGatewaySku(d),
		IPConfigurations:                expandVirtualNetworkGatewayIPConfigurations(d),
		CustomRoutes:                    expandVirtualNetworkGatewayAddressSpace(customRoute),
	}

	if v, ok := d.GetOk("dns_forwarding_enabled"); ok {
		props.EnableDNSForwarding = utils.Bool(v.(bool))
	}

	if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
		props.GatewayDefaultSite = &network.SubResource{
			ID: &gatewayDefaultSiteID,
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

	if existingVNetGateway.VirtualNetworkGatewayPropertiesFormat != nil && existingVNetGateway.VirtualNetworkGatewayPropertiesFormat.NatRules != nil {
		props.NatRules = existingVNetGateway.VirtualNetworkGatewayPropertiesFormat.NatRules
	}

	// Sku validation for policy-based VPN gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.VpnTypePolicyBased {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateVirtualNetworkGatewayPolicyBasedVpnSku()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of first geneneration
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.VpnTypeRouteBased && props.VpnGatewayGeneration == network.VpnGatewayGenerationGeneration1 {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of second geneneration
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.VpnTypeRouteBased && props.VpnGatewayGeneration == network.VpnGatewayGenerationGeneration2 {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2()); !ok {
			return nil, err
		}
	}

	// Sku validation for ExpressRoute gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeExpressRoute {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateVirtualNetworkGatewayExpressRouteSku()); !ok {
			return nil, err
		}
	}

	return props, nil
}

func expandVirtualNetworkGatewayBgpSettings(id parse.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*network.BgpSettings, error) {
	bgpSets := d.Get("bgp_settings").([]interface{})
	if len(bgpSets) == 0 {
		return nil, nil
	}

	bgp := bgpSets[0].(map[string]interface{})

	asn := int64(bgp["asn"].(int))
	peerWeight := int32(bgp["peer_weight"].(int))

	ipConfiguration := d.Get("ip_configuration").([]interface{})
	peeringAddresses, err := expandVirtualNetworkGatewayBgpPeeringAddresses(id, ipConfiguration, bgp["peering_addresses"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &network.BgpSettings{
		Asn:                 &asn,
		PeerWeight:          &peerWeight,
		BgpPeeringAddresses: peeringAddresses,
	}, nil
}

func expandVirtualNetworkGatewayBgpPeeringAddresses(id parse.VirtualNetworkGatewayId, ipConfig, input []interface{}) (*[]network.IPConfigurationBgpPeeringAddress, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]network.IPConfigurationBgpPeeringAddress, 0)

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

		ipConfigId := parse.NewVirtualNetworkGatewayIpConfigurationID(id.SubscriptionId, id.ResourceGroup, id.Name, ipConfigName)
		result = append(result, network.IPConfigurationBgpPeeringAddress{
			IpconfigurationID:    utils.String(ipConfigId.ID()),
			CustomBgpIPAddresses: utils.ExpandStringSlice(b["apipa_addresses"].([]interface{})),
		})
	}

	return &result, nil
}

func expandVirtualNetworkGatewayIPConfigurations(d *pluginsdk.ResourceData) *[]network.VirtualNetworkGatewayIPConfiguration {
	configs := d.Get("ip_configuration").([]interface{})
	ipConfigs := make([]network.VirtualNetworkGatewayIPConfiguration, 0, len(configs))

	for _, c := range configs {
		conf := c.(map[string]interface{})

		name := conf["name"].(string)
		privateIPAllocMethod := network.IPAllocationMethod(conf["private_ip_address_allocation"].(string))

		props := &network.VirtualNetworkGatewayIPConfigurationPropertiesFormat{
			PrivateIPAllocationMethod: privateIPAllocMethod,
		}

		if subnetID := conf["subnet_id"].(string); subnetID != "" {
			props.Subnet = &network.SubResource{
				ID: &subnetID,
			}
		}

		if publicIP := conf["public_ip_address_id"].(string); publicIP != "" {
			props.PublicIPAddress = &network.SubResource{
				ID: &publicIP,
			}
		}

		ipConfig := network.VirtualNetworkGatewayIPConfiguration{
			Name: &name,
			VirtualNetworkGatewayIPConfigurationPropertiesFormat: props,
		}

		ipConfigs = append(ipConfigs, ipConfig)
	}

	return &ipConfigs
}

func expandVirtualNetworkGatewayVpnClientConfig(d *pluginsdk.ResourceData, vnetGatewayId parse.VirtualNetworkGatewayId) *network.VpnClientConfiguration {
	configSets := d.Get("vpn_client_configuration").([]interface{})
	conf := configSets[0].(map[string]interface{})

	confAddresses := conf["address_space"].([]interface{})
	addresses := make([]string, 0, len(confAddresses))
	for _, addr := range confAddresses {
		addresses = append(addresses, addr.(string))
	}

	confAadTenant := conf["aad_tenant"].(string)
	confAadAudience := conf["aad_audience"].(string)
	confAadIssuer := conf["aad_issuer"].(string)

	var rootCerts []network.VpnClientRootCertificate
	for _, rootCertSet := range conf["root_certificate"].(*pluginsdk.Set).List() {
		rootCert := rootCertSet.(map[string]interface{})
		name := rootCert["name"].(string)
		publicCertData := rootCert["public_cert_data"].(string)
		r := network.VpnClientRootCertificate{
			Name: &name,
			VpnClientRootCertificatePropertiesFormat: &network.VpnClientRootCertificatePropertiesFormat{
				PublicCertData: &publicCertData,
			},
		}
		rootCerts = append(rootCerts, r)
	}

	var revokedCerts []network.VpnClientRevokedCertificate
	for _, revokedCertSet := range conf["revoked_certificate"].(*pluginsdk.Set).List() {
		revokedCert := revokedCertSet.(map[string]interface{})
		name := revokedCert["name"].(string)
		thumbprint := revokedCert["thumbprint"].(string)
		r := network.VpnClientRevokedCertificate{
			Name: &name,
			VpnClientRevokedCertificatePropertiesFormat: &network.VpnClientRevokedCertificatePropertiesFormat{
				Thumbprint: &thumbprint,
			},
		}
		revokedCerts = append(revokedCerts, r)
	}

	var vpnClientProtocols []network.VpnClientProtocol
	for _, vpnClientProtocol := range conf["vpn_client_protocols"].(*pluginsdk.Set).List() {
		p := network.VpnClientProtocol(vpnClientProtocol.(string))
		vpnClientProtocols = append(vpnClientProtocols, p)
	}

	confRadiusServerAddress := conf["radius_server_address"].(string)
	confRadiusServerSecret := conf["radius_server_secret"].(string)

	var vpnAuthTypes []network.VpnAuthenticationType
	for _, vpnAuthType := range conf["vpn_auth_types"].(*pluginsdk.Set).List() {
		a := network.VpnAuthenticationType(vpnAuthType.(string))
		vpnAuthTypes = append(vpnAuthTypes, a)
	}

	return &network.VpnClientConfiguration{
		VpnClientAddressPool: &network.AddressSpace{
			AddressPrefixes: &addresses,
		},
		AadTenant:                         &confAadTenant,
		AadAudience:                       &confAadAudience,
		AadIssuer:                         &confAadIssuer,
		VngClientConnectionConfigurations: expandVirtualNetworkGatewayClientConnections(conf["virtual_network_gateway_client_connection"].([]interface{}), vnetGatewayId),
		VpnClientIpsecPolicies:            expandVirtualNetworkGatewayIpsecPolicies(conf["ipsec_policy"].([]interface{})),
		VpnClientRootCertificates:         &rootCerts,
		VpnClientRevokedCertificates:      &revokedCerts,
		VpnClientProtocols:                &vpnClientProtocols,
		RadiusServers:                     expandVirtualNetworkGatewayRadiusServers(conf["radius_server"].([]interface{})),
		RadiusServerAddress:               &confRadiusServerAddress,
		RadiusServerSecret:                &confRadiusServerSecret,
		VpnAuthenticationTypes:            &vpnAuthTypes,
	}
}

func expandVirtualNetworkGatewaySku(d *pluginsdk.ResourceData) *network.VirtualNetworkGatewaySku {
	sku := d.Get("sku").(string)

	return &network.VirtualNetworkGatewaySku{
		Name: network.VirtualNetworkGatewaySkuName(sku),
		Tier: network.VirtualNetworkGatewaySkuTier(sku),
	}
}

func expandVirtualNetworkGatewayAddressSpace(input []interface{}) *network.AddressSpace {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &network.AddressSpace{
		AddressPrefixes: utils.ExpandStringSlice(v["address_prefixes"].(*pluginsdk.Set).List()),
	}
}

func expandVirtualNetworkGatewayIpsecPolicies(input []interface{}) *[]network.IpsecPolicy {
	results := make([]network.IpsecPolicy, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, network.IpsecPolicy{
			DhGroup:             network.DhGroup(v["dh_group"].(string)),
			IkeEncryption:       network.IkeEncryption(v["ike_encryption"].(string)),
			IkeIntegrity:        network.IkeIntegrity(v["ike_integrity"].(string)),
			IpsecEncryption:     network.IpsecEncryption(v["ipsec_encryption"].(string)),
			IpsecIntegrity:      network.IpsecIntegrity(v["ipsec_integrity"].(string)),
			PfsGroup:            network.PfsGroup(v["pfs_group"].(string)),
			SaLifeTimeSeconds:   utils.Int32(int32(v["sa_lifetime_in_seconds"].(int))),
			SaDataSizeKilobytes: utils.Int32(int32(v["sa_data_size_in_kilobytes"].(int))),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayRadiusServers(input []interface{}) *[]network.RadiusServer {
	results := make([]network.RadiusServer, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		v := item.(map[string]interface{})

		results = append(results, network.RadiusServer{
			RadiusServerAddress: utils.String(v["address"].(string)),
			RadiusServerScore:   utils.Int64(int64(v["score"].(int))),
			RadiusServerSecret:  utils.String(v["secret"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyGroups(input []interface{}) *[]network.VirtualNetworkGatewayPolicyGroup {
	results := make([]network.VirtualNetworkGatewayPolicyGroup, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyGroup := item.(map[string]interface{})

		results = append(results, network.VirtualNetworkGatewayPolicyGroup{
			Name: utils.String(policyGroup["name"].(string)),
			VirtualNetworkGatewayPolicyGroupProperties: &network.VirtualNetworkGatewayPolicyGroupProperties{
				IsDefault:     utils.Bool(policyGroup["is_default"].(bool)),
				PolicyMembers: expandVirtualNetworkGatewayPolicyMembers(policyGroup["policy_member"].([]interface{})),
				Priority:      utils.Int32(int32(policyGroup["priority"].(int))),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayPolicyMembers(input []interface{}) *[]network.VirtualNetworkGatewayPolicyGroupMember {
	results := make([]network.VirtualNetworkGatewayPolicyGroupMember, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		policyMember := item.(map[string]interface{})

		results = append(results, network.VirtualNetworkGatewayPolicyGroupMember{
			Name:           utils.String(policyMember["name"].(string)),
			AttributeType:  network.VpnPolicyMemberAttributeType(policyMember["type"].(string)),
			AttributeValue: utils.String(policyMember["value"].(string)),
		})
	}

	return &results
}

func expandVirtualNetworkGatewayClientConnections(input []interface{}, vnetGatewayId parse.VirtualNetworkGatewayId) *[]network.VngClientConnectionConfiguration {
	results := make([]network.VngClientConnectionConfiguration, 0)
	if len(input) == 0 {
		return &results
	}

	for _, item := range input {
		vngClientConnectionConfiguration := item.(map[string]interface{})

		results = append(results, network.VngClientConnectionConfiguration{
			Name: utils.String(vngClientConnectionConfiguration["name"].(string)),
			VngClientConnectionConfigurationProperties: &network.VngClientConnectionConfigurationProperties{
				VpnClientAddressPool:              expandVirtualNetworkGatewayAddressPool(vngClientConnectionConfiguration["address_prefixes"].([]interface{})),
				VirtualNetworkGatewayPolicyGroups: expandVirtualNetworkGatewayPolicyGroupNames(vngClientConnectionConfiguration["policy_group_names"].([]interface{}), vnetGatewayId),
			},
		})
	}

	return &results
}

func expandVirtualNetworkGatewayAddressPool(input []interface{}) *network.AddressSpace {
	if len(input) == 0 {
		return &network.AddressSpace{}
	}

	addressPrefixes := make([]string, 0)
	for _, v := range input {
		addressPrefixes = append(addressPrefixes, v.(string))
	}

	return &network.AddressSpace{
		AddressPrefixes: pointer.To(addressPrefixes),
	}
}

func expandVirtualNetworkGatewayPolicyGroupNames(input []interface{}, vnetGatewayId parse.VirtualNetworkGatewayId) *[]network.SubResource {
	results := make([]network.SubResource, 0)
	if len(input) == 0 {
		return &results
	}

	for _, v := range input {
		policyGroupId := parse.NewVirtualNetworkGatewayPolicyGroupID(vnetGatewayId.SubscriptionId, vnetGatewayId.ResourceGroup, vnetGatewayId.Name, v.(string))
		result := network.SubResource{
			ID: utils.String(policyGroupId.ID()),
		}
		results = append(results, result)
	}

	return &results
}

func flattenVirtualNetworkGatewayBgpSettings(settings *network.BgpSettings) ([]interface{}, error) {
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

func flattenVirtualNetworkGatewayBgpPeeringAddresses(input *[]network.IPConfigurationBgpPeeringAddress) (interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		var ipConfigName string
		if e.IpconfigurationID != nil {
			id, err := parse.VirtualNetworkGatewayIpConfigurationIDInsensitively(*e.IpconfigurationID)
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

func flattenVirtualNetworkGatewayIPConfigurations(ipConfigs *[]network.VirtualNetworkGatewayIPConfiguration) []interface{} {
	flat := make([]interface{}, 0)

	if ipConfigs != nil {
		for _, cfg := range *ipConfigs {
			props := cfg.VirtualNetworkGatewayIPConfigurationPropertiesFormat
			v := make(map[string]interface{})

			if name := cfg.Name; name != nil {
				v["name"] = *name
			}
			v["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)

			if subnet := props.Subnet; subnet != nil {
				if id := subnet.ID; id != nil {
					v["subnet_id"] = *id
				}
			}

			if pip := props.PublicIPAddress; pip != nil {
				if id := pip.ID; id != nil {
					v["public_ip_address_id"] = *id
				}
			}

			flat = append(flat, v)
		}
	}

	return flat
}

func flattenVirtualNetworkGatewayVpnClientConfig(cfg *network.VpnClientConfiguration) ([]interface{}, error) {
	if cfg == nil {
		return []interface{}{}, nil
	}
	flat := map[string]interface{}{
		"ipsec_policy":  flattenVirtualNetworkGatewayIPSecPolicies(cfg.VpnClientIpsecPolicies),
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
				"public_cert_data": *cert.VpnClientRootCertificatePropertiesFormat.PublicCertData,
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
				"thumbprint": *cert.VpnClientRevokedCertificatePropertiesFormat.Thumbprint,
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
		string(network.VirtualNetworkGatewaySkuTierBasic),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration1() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierBasic),
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuNameVpnGw1),
		string(network.VirtualNetworkGatewaySkuNameVpnGw2),
		string(network.VirtualNetworkGatewaySkuNameVpnGw3),
		string(network.VirtualNetworkGatewaySkuNameVpnGw1AZ),
		string(network.VirtualNetworkGatewaySkuNameVpnGw2AZ),
		string(network.VirtualNetworkGatewaySkuNameVpnGw3AZ),
	}, false)
}

func validateVirtualNetworkGatewayRouteBasedVpnSkuGeneration2() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuNameVpnGw2),
		string(network.VirtualNetworkGatewaySkuNameVpnGw3),
		string(network.VirtualNetworkGatewaySkuNameVpnGw4),
		string(network.VirtualNetworkGatewaySkuNameVpnGw5),
		string(network.VirtualNetworkGatewaySkuNameVpnGw2AZ),
		string(network.VirtualNetworkGatewaySkuNameVpnGw3AZ),
		string(network.VirtualNetworkGatewaySkuNameVpnGw4AZ),
		string(network.VirtualNetworkGatewaySkuNameVpnGw5AZ),
	}, false)
}

func validateVirtualNetworkGatewayExpressRouteSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuTierUltraPerformance),
		string(network.VirtualNetworkGatewaySkuNameErGw1AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw2AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw3AZ),
	}, false)
}

func flattenVirtualNetworkGatewayAddressSpace(input *network.AddressSpace) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"address_prefixes": utils.FlattenStringSlice(input.AddressPrefixes),
		},
	}
}

func flattenVirtualNetworkGatewayRadiusServers(input *[]network.RadiusServer) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"address": pointer.From(item.RadiusServerAddress),
			"secret":  pointer.From(item.RadiusServerSecret),
			"score":   pointer.From(item.RadiusServerScore),
		})
	}

	return results
}

func flattenVirtualNetworkGatewayIPSecPolicies(input *[]network.IpsecPolicy) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"dh_group":                  string(item.DhGroup),
			"ipsec_encryption":          string(item.IpsecEncryption),
			"ipsec_integrity":           string(item.IpsecIntegrity),
			"ike_encryption":            string(item.IkeEncryption),
			"ike_integrity":             string(item.IkeIntegrity),
			"pfs_group":                 string(item.PfsGroup),
			"sa_data_size_in_kilobytes": pointer.From(item.SaDataSizeKilobytes),
			"sa_lifetime_in_seconds":    pointer.From(item.SaLifeTimeSeconds),
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicyGroups(input *[]network.VirtualNetworkGatewayPolicyGroup) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"name":          pointer.From(item.Name),
			"is_default":    pointer.From(item.IsDefault),
			"policy_member": flattenVirtualNetworkGatewayPolicy(item.PolicyMembers),
			"priority":      pointer.From(item.Priority),
		})
	}
	return results
}

func flattenVirtualNetworkGatewayPolicy(input *[]network.VirtualNetworkGatewayPolicyGroupMember) []interface{} {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, map[string]interface{}{
			"name":  pointer.From(item.Name),
			"type":  string(item.AttributeType),
			"value": pointer.From(item.AttributeValue),
		})
	}
	return results
}

func flattenVirtualNetworkGatewayClientConnections(input *[]network.VngClientConnectionConfiguration) ([]interface{}, error) {
	results := make([]interface{}, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, item := range *input {
		result := map[string]interface{}{
			"name":             pointer.From(item.Name),
			"address_prefixes": pointer.From(item.VpnClientAddressPool.AddressPrefixes),
		}

		policyGroupNames, err := flattenVirtualNetworkGatewayPolicyGroupNames(item.VirtualNetworkGatewayPolicyGroups)
		if err != nil {
			return nil, err
		}
		result["policy_group_names"] = policyGroupNames

		results = append(results, result)
	}

	return results, nil
}

func flattenVirtualNetworkGatewayPolicyGroupNames(input *[]network.SubResource) ([]string, error) {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, item := range *input {
		policyGroupId, err := parse.VirtualNetworkGatewayPolicyGroupID(*item.ID)
		if err != nil {
			return nil, err
		}
		results = append(results, policyGroupId.Name)
	}

	return results, nil
}
