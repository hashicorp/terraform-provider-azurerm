package network

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"

	"github.com/hashicorp/go-azure-helpers/response"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVPNGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVpnGatewayConnectionResourceCreateUpdate,
		Read:   resourceArmVpnGatewayConnectionResourceRead,
		Update: resourceArmVpnGatewayConnectionResourceCreateUpdate,
		Delete: resourceArmVpnGatewayConnectionResourceDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.VpnConnectionID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"vpn_gateway_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VpnGatewayID,
			},

			"remote_vpn_site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VpnSiteID,
			},

			"internet_security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// Service will create a route table for the user if this is not specified.
			"routing": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated_route_table": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.HubRouteTableID,
						},
						"propagated_route_tables": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.HubRouteTableID,
							},
						},
					},
				},
			},

			"vpn_link": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"vpn_site_link_id": {
							Type:     schema.TypeString,
							Required: true,
							// The vpn site link associated with one link connection can not be updated
							ForceNew:     true,
							ValidateFunc: validate.VpnSiteLinkID,
						},

						"route_weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							Default:      0,
						},

						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IKEv1),
								string(network.IKEv2),
							}, false),
							Default: string(network.IKEv2),
						},

						"bandwidth_mbps": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
							Default:      10,
						},

						"shared_key": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"bgp_enabled": {
							Type:     schema.TypeBool,
							ForceNew: true,
							Optional: true,
							Default:  false,
						},

						"ipsec_policy": {
							Type:     schema.TypeList,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sa_lifetime_sec": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(300, 172799),
									},
									"sa_data_size_kb": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1024, 2147483647),
									},
									"encryption_algorithm": {
										Type:     schema.TypeString,
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
									"integrity_algorithm": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.IpsecIntegrityMD5),
											string(network.IpsecIntegritySHA1),
											string(network.IpsecIntegritySHA256),
											string(network.IpsecIntegrityGCMAES128),
											string(network.IpsecIntegrityGCMAES192),
											string(network.IpsecIntegrityGCMAES256),
										}, false),
									},

									"ike_encryption_algorithm": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.DES),
											string(network.DES3),
											string(network.AES128),
											string(network.AES192),
											string(network.AES256),
											string(network.GCMAES128),
											string(network.GCMAES256),
										}, false),
									},

									"ike_integrity_algorithm": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.IkeIntegrityMD5),
											string(network.IkeIntegritySHA1),
											string(network.IkeIntegritySHA256),
											string(network.IkeIntegritySHA384),
											string(network.IkeIntegrityGCMAES128),
											string(network.IkeIntegrityGCMAES256),
										}, false),
									},

									"dh_group": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.None),
											string(network.DHGroup1),
											string(network.DHGroup2),
											string(network.DHGroup14),
											string(network.DHGroup24),
											string(network.DHGroup2048),
											string(network.ECP256),
											string(network.ECP384),
										}, false),
									},

									"pfs_group": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(network.PfsGroupNone),
											string(network.PfsGroupPFS1),
											string(network.PfsGroupPFS2),
											string(network.PfsGroupPFS14),
											string(network.PfsGroupPFS24),
											string(network.PfsGroupPFS2048),
											string(network.PfsGroupPFSMM),
											string(network.PfsGroupECP256),
											string(network.PfsGroupECP384),
										}, false),
									},
								},
							},
						},

						"ratelimit_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"local_azure_ip_address_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"policy_based_traffic_selector_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceArmVpnGatewayConnectionResourceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnConnectionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	gatewayId, err := parse.VpnGatewayID(d.Get("vpn_gateway_id").(string))
	if err != nil {
		return err
	}

	if d.IsNewResource() {
		resp, err := client.Get(ctx, gatewayId.ResourceGroup, gatewayId.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", name, gatewayId.ResourceGroup, gatewayId.Name, err)
			}
		}

		if resp.ID != nil && *resp.ID != "" {
			return tf.ImportAsExistsError("azurerm_vpn_gateway_connection", *resp.ID)
		}
	}

	locks.ByName(gatewayId.Name, VPNGatewayResourceName)
	defer locks.UnlockByName(gatewayId.Name, VPNGatewayResourceName)

	param := network.VpnConnection{
		Name: &name,
		VpnConnectionProperties: &network.VpnConnectionProperties{
			EnableInternetSecurity: utils.Bool(d.Get("internet_security_enabled").(bool)),
			RemoteVpnSite: &network.SubResource{
				ID: utils.String(d.Get("remote_vpn_site_id").(string)),
			},
			VpnLinkConnections:   expandArmVpnGatewayConnectionVpnSiteLinkConnections(d.Get("vpn_link").([]interface{})),
			RoutingConfiguration: expandArmVpnGatewayConnectionRoutingConfiguration(d.Get("routing").([]interface{})),
		},
	}

	future, err := client.CreateOrUpdate(ctx, gatewayId.ResourceGroup, gatewayId.Name, name, param)
	if err != nil {
		return fmt.Errorf("creating Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", name, gatewayId.ResourceGroup, gatewayId.Name, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", name, gatewayId.ResourceGroup, gatewayId.Name, err)
	}

	resp, err := client.Get(ctx, gatewayId.ResourceGroup, gatewayId.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway: %q): %+v", name, gatewayId.ResourceGroup, gatewayId.Name, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway: %q) ID", name, gatewayId.ResourceGroup, gatewayId.Name)
	}

	id, err := parse.VpnConnectionID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(""))

	return resourceArmVpnGatewayConnectionResourceRead(d, meta)
}

func resourceArmVpnGatewayConnectionResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VpnGatewayName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Vpn Gateway Connection Resource %q was not found in VPN Gateway %q in Resource Group %q - removing from state!", id.Name, id.VpnGatewayName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", id.Name, id.ResourceGroup, id.VpnGatewayName, err)
	}

	d.Set("name", id.Name)

	gatewayId := parse.NewVpnGatewayID(id.SubscriptionId, id.ResourceGroup, id.VpnGatewayName)
	d.Set("vpn_gateway_id", gatewayId.ID(""))

	if prop := resp.VpnConnectionProperties; prop != nil {
		vpnSiteId := ""
		if site := prop.RemoteVpnSite; site != nil {
			if id := site.ID; id != nil {
				theVpnSiteId, err := parse.VpnSiteID(*id)
				if err != nil {
					return err
				}
				vpnSiteId = theVpnSiteId.ID("")
			}
		}
		d.Set("remote_vpn_site_id", vpnSiteId)

		enableInternetSecurity := false
		if prop.EnableInternetSecurity != nil {
			enableInternetSecurity = *prop.EnableInternetSecurity
		}
		d.Set("internet_security_enabled", enableInternetSecurity)

		if err := d.Set("routing", flattenArmVpnGatewayConnectionRoutingConfiguration(prop.RoutingConfiguration)); err != nil {
			return fmt.Errorf(`setting "routing": %v`, err)
		}

		if err := d.Set("vpn_link", flattenArmVpnGatewayConnectionVpnSiteLinkConnections(prop.VpnLinkConnections)); err != nil {
			return fmt.Errorf(`setting "vpn_link": %v`, err)
		}
	}

	return nil
}

func resourceArmVpnGatewayConnectionResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VpnConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VpnGatewayName, VPNGatewayResourceName)
	defer locks.UnlockByName(id.VpnGatewayName, VPNGatewayResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VpnGatewayName, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", id.Name, id.ResourceGroup, id.VpnGatewayName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for the deletion of VPN Gateway Connection %q (Resource Group %q / VPN Gateway %q): %+v", id.Name, id.ResourceGroup, id.VpnGatewayName, err)
		}
	}

	return nil
}

func expandArmVpnGatewayConnectionVpnSiteLinkConnections(input []interface{}) *[]network.VpnSiteLinkConnection {
	if len(input) == 0 {
		return nil
	}

	result := make([]network.VpnSiteLinkConnection, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		v := network.VpnSiteLinkConnection{
			Name: utils.String(e["name"].(string)),
			VpnSiteLinkConnectionProperties: &network.VpnSiteLinkConnectionProperties{
				VpnSiteLink:                    &network.SubResource{ID: utils.String(e["vpn_site_link_id"].(string))},
				RoutingWeight:                  utils.Int32(int32(e["route_weight"].(int))),
				VpnConnectionProtocolType:      network.VirtualNetworkGatewayConnectionProtocol(e["protocol"].(string)),
				ConnectionBandwidth:            utils.Int32(int32(e["bandwidth_mbps"].(int))),
				EnableBgp:                      utils.Bool(e["bgp_enabled"].(bool)),
				IpsecPolicies:                  expandArmVpnGatewayConnectionIpSecPolicies(e["ipsec_policy"].([]interface{})),
				EnableRateLimiting:             utils.Bool(e["ratelimit_enabled"].(bool)),
				UseLocalAzureIPAddress:         utils.Bool(e["local_azure_ip_address_enabled"].(bool)),
				UsePolicyBasedTrafficSelectors: utils.Bool(e["policy_based_traffic_selector_enabled"].(bool)),
			},
		}

		if sharedKey := e["shared_key"]; sharedKey != "" {
			sharedKey := sharedKey.(string)
			v.VpnSiteLinkConnectionProperties.SharedKey = &sharedKey
		}
		result = append(result, v)
	}

	return &result
}

func flattenArmVpnGatewayConnectionVpnSiteLinkConnections(input *[]network.VpnSiteLinkConnection) interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		name := ""
		if e.Name != nil {
			name = *e.Name
		}

		vpnSiteLinkId := ""
		if e.VpnSiteLink != nil && e.VpnSiteLink.ID != nil {
			vpnSiteLinkId = *e.VpnSiteLink.ID
		}

		routeWeight := 0
		if e.RoutingWeight != nil {
			routeWeight = int(*e.RoutingWeight)
		}

		bandwidth := 0
		if e.ConnectionBandwidth != nil {
			bandwidth = int(*e.ConnectionBandwidth)
		}

		sharedKey := ""
		if e.SharedKey != nil {
			sharedKey = *e.SharedKey
		}

		bgpEnabled := false
		if e.EnableBgp != nil {
			bgpEnabled = *e.EnableBgp
		}

		usePolicyBased := false
		if e.UsePolicyBasedTrafficSelectors != nil {
			usePolicyBased = *e.UsePolicyBasedTrafficSelectors
		}

		rateLimitEnabled := false
		if e.EnableRateLimiting != nil {
			rateLimitEnabled = *e.EnableRateLimiting
		}

		useLocalAzureIpAddress := false
		if e.UseLocalAzureIPAddress != nil {
			useLocalAzureIpAddress = *e.UseLocalAzureIPAddress
		}

		v := map[string]interface{}{
			"name":                                  name,
			"vpn_site_link_id":                      vpnSiteLinkId,
			"route_weight":                          routeWeight,
			"protocol":                              string(e.VpnConnectionProtocolType),
			"bandwidth_mbps":                        bandwidth,
			"shared_key":                            sharedKey,
			"bgp_enabled":                           bgpEnabled,
			"ipsec_policy":                          flattenArmVpnGatewayConnectionIpSecPolicies(e.IpsecPolicies),
			"ratelimit_enabled":                     rateLimitEnabled,
			"local_azure_ip_address_enabled":        useLocalAzureIpAddress,
			"policy_based_traffic_selector_enabled": usePolicyBased,
		}

		output = append(output, v)
	}

	return output
}

func expandArmVpnGatewayConnectionIpSecPolicies(input []interface{}) *[]network.IpsecPolicy {
	if len(input) == 0 {
		return nil
	}

	result := make([]network.IpsecPolicy, 0)

	for _, e := range input {
		e := e.(map[string]interface{})
		v := network.IpsecPolicy{
			SaLifeTimeSeconds:   utils.Int32(int32(e["sa_lifetime_sec"].(int))),
			SaDataSizeKilobytes: utils.Int32(int32(e["sa_data_size_kb"].(int))),
			IpsecEncryption:     network.IpsecEncryption(e["encryption_algorithm"].(string)),
			IpsecIntegrity:      network.IpsecIntegrity(e["integrity_algorithm"].(string)),
			IkeEncryption:       network.IkeEncryption(e["ike_encryption_algorithm"].(string)),
			IkeIntegrity:        network.IkeIntegrity(e["ike_integrity_algorithm"].(string)),
			DhGroup:             network.DhGroup(e["dh_group"].(string)),
			PfsGroup:            network.PfsGroup(e["pfs_group"].(string)),
		}
		result = append(result, v)
	}

	return &result
}

func flattenArmVpnGatewayConnectionIpSecPolicies(input *[]network.IpsecPolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)

	for _, e := range *input {
		saLifetimeSec := 0
		if e.SaLifeTimeSeconds != nil {
			saLifetimeSec = int(*e.SaLifeTimeSeconds)
		}

		saDataSizeKb := 0
		if e.SaDataSizeKilobytes != nil {
			saDataSizeKb = int(*e.SaDataSizeKilobytes)
		}

		v := map[string]interface{}{
			"sa_lifetime_sec":          saLifetimeSec,
			"sa_data_size_kb":          saDataSizeKb,
			"encryption_algorithm":     string(e.IpsecEncryption),
			"integrity_algorithm":      string(e.IpsecIntegrity),
			"ike_encryption_algorithm": string(e.IkeEncryption),
			"ike_integrity_algorithm":  string(e.IkeIntegrity),
			"dh_group":                 string(e.DhGroup),
			"pfs_group":                string(e.PfsGroup),
		}

		output = append(output, v)
	}

	return output
}

func expandArmVpnGatewayConnectionRoutingConfiguration(input []interface{}) *network.RoutingConfiguration {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	raw := input[0].(map[string]interface{})
	output := &network.RoutingConfiguration{
		AssociatedRouteTable:  &network.SubResource{ID: utils.String(raw["associated_route_table"].(string))},
		PropagatedRouteTables: &network.PropagatedRouteTable{Ids: expandNetworkSubResourceID(raw["propagated_route_tables"].([]interface{}))},
	}

	return output
}

func flattenArmVpnGatewayConnectionRoutingConfiguration(input *network.RoutingConfiguration) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	associateRouteTable := ""
	if input.AssociatedRouteTable != nil && input.AssociatedRouteTable.ID != nil {
		associateRouteTable = *input.AssociatedRouteTable.ID
	}

	propagatedRouteTables := []interface{}{}
	if input.PropagatedRouteTables != nil && input.PropagatedRouteTables.Ids != nil {
		for _, subresource := range *input.PropagatedRouteTables.Ids {
			id := ""
			if subresource.ID != nil {
				id = *subresource.ID
			}
			propagatedRouteTables = append(propagatedRouteTables, id)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"associated_route_table":  associateRouteTable,
			"propagated_route_tables": propagatedRouteTables,
		},
	}
}
