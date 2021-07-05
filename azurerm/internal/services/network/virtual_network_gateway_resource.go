package network

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVirtualNetworkGateway() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayCreateUpdate,
		Read:   resourceVirtualNetworkGatewayRead,
		Update: resourceVirtualNetworkGatewayCreateUpdate,
		Delete: resourceVirtualNetworkGatewayDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(resourceVirtualNetworkGatewayCustomizeDiff),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"type": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayTypeExpressRoute),
					string(network.VirtualNetworkGatewayTypeVpn),
				}, true),
			},

			"vpn_type": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(network.VpnTypeRouteBased),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VpnTypeRouteBased),
					string(network.VpnTypePolicyBased),
				}, true),
			},

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
				Type:             pluginsdk.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
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
				MaxItems: 2,
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
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
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
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},
						"aad_audience": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},
						"aad_issuer": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},

						"root_certificate": {
							Type:     pluginsdk.TypeSet,
							Optional: true,

							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
							},
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
							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
							},
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

						"radius_server_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
							ValidateFunc: validation.IsIPv4Address,
						},

						"radius_server_secret": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
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
								}, true),
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
							AtLeastOneOf: []string{"bgp_settings.0.asn", "bgp_settings.0.peering_address",
								"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses"},
						},

						// TODO 3.0 - Remove this property
						"peering_address": {
							Type:       pluginsdk.TypeString,
							Optional:   true,
							Computed:   true,
							Deprecated: "Deprecated in favor of `bgp_settings.0.peering_addresses.0.default_addresses.0`",
							AtLeastOneOf: []string{"bgp_settings.0.asn", "bgp_settings.0.peering_address",
								"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses"},
						},

						"peer_weight": {
							Type:     pluginsdk.TypeInt,
							Optional: true,
							AtLeastOneOf: []string{"bgp_settings.0.asn", "bgp_settings.0.peering_address",
								"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses"},
						},

						//lintignore:XS003
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
							AtLeastOneOf: []string{"bgp_settings.0.asn", "bgp_settings.0.peering_address",
								"bgp_settings.0.peer_weight", "bgp_settings.0.peering_addresses"},
						},
					},
				},
			},

			//lintignore:XS003
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
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualNetworkGatewayCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	id := parse.NewVirtualNetworkGatewayID(subscriptionId, resGroup, name)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Network Gateway %q: %s", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			id, err := parse.VirtualNetworkGatewayID(*existing.ID)
			if err != nil {
				return err
			}
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties, err := getVirtualNetworkGatewayProperties(id, d)
	if err != nil {
		return err
	}

	gateway := network.VirtualNetworkGateway{
		Name:                                  &name,
		Location:                              &location,
		Tags:                                  tags.Expand(t),
		VirtualNetworkGatewayPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, gateway)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AzureRM Virtual Network Gateway %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of AzureRM Virtual Network Gateway %q: %+v", id, err)
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
		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway %q: %+v", id, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if gw := resp.VirtualNetworkGatewayPropertiesFormat; gw != nil {
		d.Set("type", string(gw.GatewayType))
		d.Set("enable_bgp", gw.EnableBgp)
		d.Set("private_ip_address_enabled", gw.EnablePrivateIPAddress)
		d.Set("active_active", gw.ActiveActive)
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
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("vpn_client_configuration", flattenVirtualNetworkGatewayVpnClientConfig(gw.VpnClientConfiguration)); err != nil {
			return fmt.Errorf("Error setting `vpn_client_configuration`: %+v", err)
		}

		bgpSettings, err := flattenVirtualNetworkGatewayBgpSettings(gw.BgpSettings)
		if err != nil {
			return err
		}
		if err := d.Set("bgp_settings", bgpSettings); err != nil {
			return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
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
		return fmt.Errorf("Error deleting Virtual Network Gateway %q: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway %q: %+v", id, err)
	}

	return nil
}

func getVirtualNetworkGatewayProperties(id parse.VirtualNetworkGatewayId, d *pluginsdk.ResourceData) (*network.VirtualNetworkGatewayPropertiesFormat, error) {
	gatewayType := network.VirtualNetworkGatewayType(d.Get("type").(string))
	vpnType := network.VpnType(d.Get("vpn_type").(string))
	enableBgp := d.Get("enable_bgp").(bool)
	enablePrivateIpAddress := d.Get("private_ip_address_enabled").(bool)
	activeActive := d.Get("active_active").(bool)
	generation := network.VpnGatewayGeneration(d.Get("generation").(string))
	customRoute := d.Get("custom_route").([]interface{})

	props := &network.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:            gatewayType,
		VpnType:                vpnType,
		EnableBgp:              &enableBgp,
		EnablePrivateIPAddress: &enablePrivateIpAddress,
		ActiveActive:           &activeActive,
		VpnGatewayGeneration:   generation,
		Sku:                    expandVirtualNetworkGatewaySku(d),
		IPConfigurations:       expandVirtualNetworkGatewayIPConfigurations(d),
		CustomRoutes:           expandVirtualNetworkGatewayAddressSpace(customRoute),
	}

	if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
		props.GatewayDefaultSite = &network.SubResource{
			ID: &gatewayDefaultSiteID,
		}
	}

	if _, ok := d.GetOk("vpn_client_configuration"); ok {
		props.VpnClientConfiguration = expandVirtualNetworkGatewayVpnClientConfig(d)
	}

	if _, ok := d.GetOk("bgp_settings"); ok {
		bgpSettings, err := expandVirtualNetworkGatewayBgpSettings(id, d)
		if err != nil {
			return nil, err
		}
		props.BgpSettings = bgpSettings
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
	peeringAddress := bgp["peering_address"].(string)
	peerWeight := int32(bgp["peer_weight"].(int))

	ipConfiguration := d.Get("ip_configuration").([]interface{})
	peeringAddresses, err := expandVirtualNetworkGatewayBgpPeeringAddresses(id, ipConfiguration, bgp["peering_addresses"].([]interface{}))
	if err != nil {
		return nil, err
	}

	return &network.BgpSettings{
		Asn:                 &asn,
		BgpPeeringAddress:   &peeringAddress,
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

func expandVirtualNetworkGatewayVpnClientConfig(d *pluginsdk.ResourceData) *network.VpnClientConfiguration {
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

	return &network.VpnClientConfiguration{
		VpnClientAddressPool: &network.AddressSpace{
			AddressPrefixes: &addresses,
		},
		AadTenant:                    &confAadTenant,
		AadAudience:                  &confAadAudience,
		AadIssuer:                    &confAadIssuer,
		VpnClientRootCertificates:    &rootCerts,
		VpnClientRevokedCertificates: &revokedCerts,
		VpnClientProtocols:           &vpnClientProtocols,
		RadiusServerAddress:          &confRadiusServerAddress,
		RadiusServerSecret:           &confRadiusServerSecret,
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

func flattenVirtualNetworkGatewayBgpSettings(settings *network.BgpSettings) ([]interface{}, error) {
	output := make([]interface{}, 0)

	if settings != nil {
		flat := make(map[string]interface{})

		if asn := settings.Asn; asn != nil {
			flat["asn"] = int(*asn)
		}
		if address := settings.BgpPeeringAddress; address != nil {
			flat["peering_address"] = *address
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
			id, err := parse.VirtualNetworkGatewayIpConfigurationID(*e.IpconfigurationID)
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

func flattenVirtualNetworkGatewayVpnClientConfig(cfg *network.VpnClientConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}
	flat := make(map[string]interface{})

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

	if v := cfg.RadiusServerAddress; v != nil {
		flat["radius_server_address"] = *v
	}

	if v := cfg.RadiusServerSecret; v != nil {
		flat["radius_server_secret"] = *v
	}

	return []interface{}{flat}
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
	}, true)
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
	}, true)
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
	}, true)
}

func validateVirtualNetworkGatewayExpressRouteSku() pluginsdk.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuTierUltraPerformance),
		string(network.VirtualNetworkGatewaySkuNameErGw1AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw2AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw3AZ),
	}, true)
}

func resourceVirtualNetworkGatewayCustomizeDiff(ctx context.Context, diff *pluginsdk.ResourceDiff, _ interface{}) error {
	if vpnClient, ok := diff.GetOk("vpn_client_configuration"); ok {
		if vpnClientConfig, ok := vpnClient.([]interface{})[0].(map[string]interface{}); ok {
			hasAadTenant := vpnClientConfig["aad_tenant"] != ""
			hasAadAudience := vpnClientConfig["aad_audience"] != ""
			hasAadIssuer := vpnClientConfig["aad_issuer"] != ""

			if hasAadTenant && (!hasAadAudience || !hasAadIssuer) {
				return fmt.Errorf("if aad_tenant is set aad_audience and aad_issuer must also be set")
			}
			if hasAadAudience && (!hasAadTenant || !hasAadIssuer) {
				return fmt.Errorf("if aad_audience is set aad_tenant and aad_issuer must also be set")
			}
			if hasAadIssuer && (!hasAadTenant || !hasAadAudience) {
				return fmt.Errorf("if aad_issuer is set aad_tenant and aad_audience must also be set")
			}

			hasRadiusAddress := vpnClientConfig["radius_server_address"] != ""
			hasRadiusSecret := vpnClientConfig["radius_server_secret"] != ""

			if hasRadiusAddress && !hasRadiusSecret {
				return fmt.Errorf("if radius_server_address is set radius_server_secret must also be set")
			}
			if !hasRadiusAddress && hasRadiusSecret {
				return fmt.Errorf("if radius_server_secret is set radius_server_address must also be set")
			}
		}
	}
	return nil
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
