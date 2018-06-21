package azurerm

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayTypeExpressRoute),
					string(network.VirtualNetworkGatewayTypeVpn),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"vpn_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  string(network.RouteBased),
				ValidateFunc: validation.StringInSlice([]string{
					string(network.RouteBased),
					string(network.PolicyBased),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewaySkuTierBasic),
					string(network.VirtualNetworkGatewaySkuTierStandard),
					string(network.VirtualNetworkGatewaySkuTierHighPerformance),
					string(network.VirtualNetworkGatewaySkuTierUltraPerformance),
					string(network.VirtualNetworkGatewaySkuNameVpnGw1),
					string(network.VirtualNetworkGatewaySkuNameVpnGw2),
					string(network.VirtualNetworkGatewaySkuNameVpnGw3),
				}, true),
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
						"public_ip_address_id": {
							Type:     schema.TypeString,
							Optional: true,
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
						"vpn_client_protocols": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									string(network.IkeV2),
									string(network.SSTP),
								}, true),
							},
						},
						"root_certificate": {
							Type:     schema.TypeSet,
							Optional: true,

							ConflictsWith: []string{
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
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
							},
							ValidateFunc: validate.Ip4Address,
						},
						"radius_server_secret": {
							Type:     schema.TypeString,
							Optional: true,
							ConflictsWith: []string{
								"vpn_client_configuration.0.root_certificate",
								"vpn_client_configuration.0.revoked_certificate",
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
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkGatewayCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties, err := getArmVirtualNetworkGatewayProperties(d)
	if err != nil {
		return err
	}

	gateway := network.VirtualNetworkGateway{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		VirtualNetworkGatewayPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, gateway)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AzureRM Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
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
	client := meta.(*ArmClient).vnetGatewayClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if resp.VirtualNetworkGatewayPropertiesFormat != nil {
		gw := *resp.VirtualNetworkGatewayPropertiesFormat

		d.Set("type", string(gw.GatewayType))
		d.Set("enable_bgp", gw.EnableBgp)
		d.Set("active_active", gw.ActiveActive)

		if string(gw.VpnType) != "" {
			d.Set("vpn_type", string(gw.VpnType))
		}

		if gw.GatewayDefaultSite != nil {
			d.Set("default_local_network_gateway_id", gw.GatewayDefaultSite.ID)
		}

		if gw.Sku != nil {
			d.Set("sku", string(gw.Sku.Name))
		}

		d.Set("ip_configuration", flattenArmVirtualNetworkGatewayIPConfigurations(gw.IPConfigurations))

		if gw.VpnClientConfiguration != nil {
			vpnConfigFlat := flattenArmVirtualNetworkGatewayVpnClientConfig(gw.VpnClientConfiguration)
			if err := d.Set("vpn_client_configuration", vpnConfigFlat); err != nil {
				return fmt.Errorf("Error setting `vpn_client_configuration`: %+v", err)
			}
		}

		if gw.BgpSettings != nil {
			bgpSettingsFlat := flattenArmVirtualNetworkGatewayBgpSettings(gw.BgpSettings)
			if err := d.Set("bgp_settings", bgpSettingsFlat); err != nil {
				return fmt.Errorf("Error setting `bgp_settings`: %+v", err)
			}
		}
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualNetworkGatewayDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getArmVirtualNetworkGatewayProperties(d *schema.ResourceData) (*network.VirtualNetworkGatewayPropertiesFormat, error) {
	gatewayType := network.VirtualNetworkGatewayType(d.Get("type").(string))
	vpnType := network.VpnType(d.Get("vpn_type").(string))
	enableBgp := d.Get("enable_bgp").(bool)
	activeActive := d.Get("active_active").(bool)

	props := &network.VirtualNetworkGatewayPropertiesFormat{
		GatewayType:      gatewayType,
		VpnType:          vpnType,
		EnableBgp:        &enableBgp,
		ActiveActive:     &activeActive,
		Sku:              expandArmVirtualNetworkGatewaySku(d),
		IPConfigurations: expandArmVirtualNetworkGatewayIPConfigurations(d),
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
		ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayPolicyBasedVpnSku())

		if !ok {
			return nil, err
		}
	}

	// Sku validation for route-based VPN gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeVpn && props.VpnType == network.RouteBased {
		ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayRouteBasedVpnSku())

		if !ok {
			return nil, err
		}
	}

	// Sku validation for ExpressRoute gateways
	if props.GatewayType == network.VirtualNetworkGatewayTypeExpressRoute {
		ok, err := evaluateSchemaValidateFunc(string(props.Sku.Name), "sku", validateArmVirtualNetworkGatewayExpressRouteSku())

		if !ok {
			return nil, err
		}
	}

	return props, nil
}

func expandArmVirtualNetworkGatewayBgpSettings(d *schema.ResourceData) *network.BgpSettings {
	bgpSets := d.Get("bgp_settings").([]interface{})
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
	flat := make(map[string]interface{})

	flat["asn"] = int(*settings.Asn)
	flat["peering_address"] = *settings.BgpPeeringAddress
	flat["peer_weight"] = int(*settings.PeerWeight)

	return []interface{}{flat}
}

func flattenArmVirtualNetworkGatewayIPConfigurations(ipConfigs *[]network.VirtualNetworkGatewayIPConfiguration) []interface{} {
	flat := make([]interface{}, 0, len(*ipConfigs))

	for _, cfg := range *ipConfigs {
		props := cfg.VirtualNetworkGatewayIPConfigurationPropertiesFormat
		v := make(map[string]interface{})

		v["name"] = *cfg.Name
		v["private_ip_address_allocation"] = string(props.PrivateIPAllocationMethod)
		v["subnet_id"] = *props.Subnet.ID
		v["public_ip_address_id"] = *props.PublicIPAddress.ID

		flat = append(flat, v)
	}

	return flat
}

func flattenArmVirtualNetworkGatewayVpnClientConfig(cfg *network.VpnClientConfiguration) []interface{} {
	flat := make(map[string]interface{})

	addressSpace := make([]interface{}, 0)
	if pool := cfg.VpnClientAddressPool; pool != nil {
		if prefixes := pool.AddressPrefixes; prefixes != nil {
			for _, addr := range *prefixes {
				addressSpace = append(addressSpace, addr)
			}
		}
	}
	flat["address_space"] = addressSpace

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
	id, err := parseAzureResourceID(virtualNetworkGatewayId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["virtualNetworkGateways"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func validateArmVirtualNetworkGatewaySubnetId(i interface{}, k string) (s []string, es []error) {
	value, ok := i.(string)
	if !ok {
		es = append(es, fmt.Errorf("expected type of %s to be string", k))
		return
	}

	id, err := parseAzureResourceID(value)
	if err != nil {
		es = append(es, fmt.Errorf("expected %s to be an Azure resource id", k))
		return
	}

	subnet, ok := id.Path["subnets"]
	if !ok {
		es = append(es, fmt.Errorf("expected %s to reference a subnet resource", k))
		return
	}

	if strings.ToLower(subnet) != "gatewaysubnet" {
		es = append(es, fmt.Errorf("expected %s to reference a gateway subnet with name GatewaySubnet", k))
	}

	return
}

func validateArmVirtualNetworkGatewayPolicyBasedVpnSku() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierBasic),
	}, true)
}

func validateArmVirtualNetworkGatewayRouteBasedVpnSku() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierBasic),
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuNameVpnGw1),
		string(network.VirtualNetworkGatewaySkuNameVpnGw2),
		string(network.VirtualNetworkGatewaySkuNameVpnGw3),
	}, true)
}

func validateArmVirtualNetworkGatewayExpressRouteSku() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		string(network.VirtualNetworkGatewaySkuTierStandard),
		string(network.VirtualNetworkGatewaySkuTierHighPerformance),
		string(network.VirtualNetworkGatewaySkuTierUltraPerformance),
	}, true)
}

func resourceArmVirtualNetworkGatewayCustomizeDiff(diff *schema.ResourceDiff, v interface{}) error {

	if vpnClient, ok := diff.GetOk("vpn_client_configuration"); ok {
		if vpnClientConfig, ok := vpnClient.([]interface{})[0].(map[string]interface{}); ok {
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
