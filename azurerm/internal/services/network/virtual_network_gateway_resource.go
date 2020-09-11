package network

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualNetworkGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkGatewayCreateUpdate,
		Read:   resourceArmVirtualNetworkGatewayRead,
		Update: resourceArmVirtualNetworkGatewayCreateUpdate,
		Delete: resourceArmVirtualNetworkGatewayDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		CustomizeDiff: resourceArmVirtualNetworkGatewayCustomizeDiff,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"type": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayTypeExpressRoute),
					string(network.VirtualNetworkGatewayTypeVpn),
				}, true),
			},

			"vpn_type": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Default:          string(network.RouteBased),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.RouteBased),
					string(network.PolicyBased),
				}, true),
			},

			"enable_bgp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"active_active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"sku": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				// This validator checks for all possible values for the SKU regardless of the attributes vpn_type and
				// type. For a validation which depends on the attributes vpn_type and type, refer to the special case
				// validators validateArmVirtualNetworkGatewayPolicyBasedVpnSku, validateArmVirtualNetworkGatewayRouteBasedVpnSku
				// and validateArmVirtualNetworkGatewayExpressRouteSku.
				ValidateFunc: validation.Any(
					validateArmVirtualNetworkGatewayPolicyBasedVpnSku(),
					validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration1(),
					validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration2(),
					validateArmVirtualNetworkGatewayExpressRouteSku(),
				),
			},

			"generation": {
				Type:     schema.TypeString,
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
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 2,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							// Azure Management API requires a name but does not generate a name if the field is missing
							// The name "vnetGatewayConfig" is used when creating a virtual network gateway via the
							// Azure portal.
							Default: "vnetGatewayConfig",
						},

						"private_ip_address_allocation": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Static),
								string(network.Dynamic),
							}, false),
							Default: string(network.Dynamic),
						},

						"subnet_id": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateFunc:     validateArmVirtualNetworkGatewaySubnetId,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"public_ip_address_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceIDOrEmpty,
						},
					},
				},
			},

			"vpn_client_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_space": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"aad_tenant": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},
						"aad_audience": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},
						"aad_issuer": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
						},

						"root_certificate": {
							Type:     schema.TypeSet,
							Optional: true,

							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"public_cert_data": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Set: hashVirtualNetworkGatewayRootCert,
						},

						"revoked_certificate": {
							Type:     schema.TypeSet,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.radius_server_address",
								"vpn_client_configuration.0.radius_server_secret",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"thumbprint": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Set: hashVirtualNetworkGatewayRevokedCert,
						},

						"radius_server_address": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.aad_tenant",
								"vpn_client_configuration.0.aad_audience",
								"vpn_client_configuration.0.aad_issuer",
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
							ValidateFunc: validate.IPv4Address,
						},

						"radius_server_secret": {
							Type:     schema.TypeString,
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
							Type:     schema.TypeSet,
							Optional: true,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.IkeV2),
									string(network.OpenVPN),
									string(network.SSTP),
								}, true),
							},
						},
					},
				},
			},

			"bgp_settings": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"peering_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},

						"peer_weight": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			"default_local_network_gateway_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmVirtualNetworkGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Network Gateway %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties, err := getArmVirtualNetworkGatewayProperties(d)
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
		return fmt.Errorf("Error Creating/Updating AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Virtual Network Gateway %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkGatewayRead(d, meta)
}

func resourceArmVirtualNetworkGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if gw := resp.VirtualNetworkGatewayPropertiesFormat; gw != nil {
		d.Set("type", string(gw.GatewayType))
		d.Set("enable_bgp", gw.EnableBgp)
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

		if err := d.Set("ip_configuration", flattenArmVirtualNetworkGatewayIPConfigurations(gw.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configuration`: %+v", err)
		}

		if err := d.Set("vpn_client_configuration", flattenArmVirtualNetworkGatewayVpnClientConfig(gw.VpnClientConfiguration)); err != nil {
			return fmt.Errorf("Error setting `vpn_client_configuration`: %+v", err)
		}

		if err := d.Set("bgp_settings", flattenArmVirtualNetworkGatewayBgpSettings(gw.BgpSettings)); err != nil {
			return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmVirtualNetworkGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getArmVirtualNetworkGatewayProperties(d *schema.ResourceData) (*network.VirtualNetworkGatewayPropertiesFormat, error) {
	gatewayType := network.VirtualNetworkGatewayType(d.Get("type").(string))
	vpnType := network.VpnType(d.Get("vpn_type").(string))
	enableBgp := d.Get("enable_bgp").(bool)
	activeActive := d.Get("active_active").(bool)
	generation := network.VpnGatewayGeneration(d.Get("generation").(string))

	props := &network.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:          gatewayType,
		VpnType:              vpnType,
		EnableBgp:            &enableBgp,
		ActiveActive:         &activeActive,
		VpnGatewayGeneration: generation,
		Sku:                  expandArmVirtualNetworkGatewaySku(d),
		IPConfigurations:     expandArmVirtualNetworkGatewayIPConfigurations(d),
	}

	if gatewayDefaultSiteID := d.Get("default_local_network_gateway_id").(string); gatewayDefaultSiteID != "" {
		props.GatewayDefaultSite = &network.SubResource{
			ID: &gatewayDefaultSiteID,
		}
	}

	if _, ok := d.GetOk("vpn_client_configuration"); ok {
		props.VpnClientConfiguration = expandArmVirtualNetworkGatewayVpnClientConfig(d)
	}

	if _, ok := d.GetOk("bgp_settings"); ok {
		props.BgpSettings = expandArmVirtualNetworkGatewayBgpSettings(d)
	}

	// Sku validation for policy-based VPN gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.PolicyBased {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayPolicyBasedVpnSku()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of first geneneration
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.RouteBased && props.VpnGatewayGeneration == network.VpnGatewayGenerationGeneration1 {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration1()); !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways of second geneneration
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.RouteBased && props.VpnGatewayGeneration == network.VpnGatewayGenerationGeneration2 {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration2()); !ok {
			return nil, err
		}
	}

	// Sku validation for ExpressRoute gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeExpressRoute {
		if ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayExpressRouteSku()); !ok {
			return nil, err
		}
	}

	return props, nil
}

func expandArmVirtualNetworkGatewayBgpSettings(d *schema.ResourceData) *network.BgpSettings {
	bgpSets := d.Get("bgp_settings").([]interface{})
	if len(bgpSets) == 0 {
		return nil
	}

	bgp := bgpSets[0].(map[string]interface{})

	asn := int64(bgp["asn"].(int))
	peeringAddress := bgp["peering_address"].(string)
	peerWeight := int32(bgp["peer_weight"].(int))

	return &network.BgpSettings{
		Asn:               &asn,
		BgpPeeringAddress: &peeringAddress,
		PeerWeight:        &peerWeight,
	}
}

func expandArmVirtualNetworkGatewayIPConfigurations(d *schema.ResourceData) *[]network.VirtualNetworkGatewayIPConfiguration {
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

func expandArmVirtualNetworkGatewayVpnClientConfig(d *schema.ResourceData) *network.VpnClientConfiguration {
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
	for _, rootCertSet := range conf["root_certificate"].(*schema.Set).List() {
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
	for _, revokedCertSet := range conf["revoked_certificate"].(*schema.Set).List() {
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
	for _, vpnClientProtocol := range conf["vpn_client_protocols"].(*schema.Set).List() {
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

func expandArmVirtualNetworkGatewaySku(d *schema.ResourceData) *network.VirtualNetworkGatewaySku {
	sku := d.Get("sku").(string)

	return &network.VirtualNetworkGatewaySku{
		Name: network.VirtualNetworkGatewaySkuName(sku),
		Tier: network.VirtualNetworkGatewaySkuTier(sku),
	}
}

func flattenArmVirtualNetworkGatewayBgpSettings(settings *network.BgpSettings) []interface{} {
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

		output = append(output, flat)
	}

	return output
}

func flattenArmVirtualNetworkGatewayIPConfigurations(ipConfigs *[]network.VirtualNetworkGatewayIPConfiguration) []interface{} {
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

func flattenArmVirtualNetworkGatewayVpnClientConfig(cfg *network.VpnClientConfiguration) []interface{} {
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
	flat["root_certificate"] = schema.NewSet(hashVirtualNetworkGatewayRootCert, rootCerts)

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
	flat["revoked_certificate"] = schema.NewSet(hashVirtualNetworkGatewayRevokedCert, revokedCerts)

	vpnClientProtocols := &schema.Set{F: schema.HashString}
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

	return hashcode.String(buf.String())
}

func hashVirtualNetworkGatewayRevokedCert(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["thumbprint"].(string)))

	return hashcode.String(buf.String())
}

func resourceGroupAndVirtualNetworkGatewayFromId(virtualNetworkGatewayId string) (string, string, error) {
	id, err := azure.ParseAzureResourceID(virtualNetworkGatewayId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["virtualNetworkGateways"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func validateArmVirtualNetworkGatewaySubnetId(i interface{}, k string) (warnings []string, errors []error) {
	value, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := azure.ParseAzureResourceID(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("expected %s to be an Azure resource id", k))
		return
	}

	subnet, ok := id.Path["subnets"]
	if !ok {
		errors = append(errors, fmt.Errorf("expected %s to reference a subnet resource", k))
		return
	}

	if strings.ToLower(subnet) != "gatewaysubnet" {
		errors = append(errors, fmt.Errorf("expected %s to reference a gateway subnet with name GatewaySubnet", k))
	}

	return warnings, errors
}

func validateArmVirtualNetworkGatewayPolicyBasedVpnSku() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierBasic),
	}, true)
}

func validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration1() schema.SchemaValidateFunc {
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

func validateArmVirtualNetworkGatewayRouteBasedVpnSkuGeneration2() schema.SchemaValidateFunc {
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

func validateArmVirtualNetworkGatewayExpressRouteSku() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuTierUltraPerformance),
		string(network.VirtualNetworkGatewaySkuNameErGw1AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw2AZ),
		string(network.VirtualNetworkGatewaySkuNameErGw3AZ),
	}, true)
}

func resourceArmVirtualNetworkGatewayCustomizeDiff(diff *schema.ResourceDiff, _ interface{}) error {
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
