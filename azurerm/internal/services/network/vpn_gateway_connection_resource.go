package network

import (
	"fmt"
	"log"
	"time"

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
			_, err := parse.VPNGatewayConnectionID(id)
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
				ValidateFunc: validate.VPNGatewayID,
			},

			"remote_vpn_site_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.VpnSiteID,
			},

			// TODO: make it settable once the routetable PR is ready
			// https://github.com/terraform-providers/terraform-provider-azurerm/pull/8939
			"routing_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"associated_route_table": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"propagated_route_tables": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},

			"vpn_link_connection": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"vpn_site_link_id": {
							Type: schema.TypeString,
							Required: true,
							ForceNew: true,
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
							Optional: true,
							Default:  false,
						},

						"use_policy_based_traffic_selector": {
							Type:     schema.TypeBool,
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
										Type:     schema.TypeInt,
										Required: true,
									},
									"sa_data_size_kb": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"ipsec_encryption_algorithm": {
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
									"ipsec_integrity_algorithm": {
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
										Type:     schema.TypeInt,
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

						"use_local_azure_ip_address": {
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	gatewayId, err := parse.VPNGatewayID(d.Get("vpn_gateway_id").(string))
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
			return tf.ImportAsExistsError("azurerm_vpn_gateway_connection_resource", *resp.ID)
		}
	}

	// TODO: d.Get() && Create
	param := network.VpnConnection{
		Name: &name,
		VpnConnectionProperties: &network.VpnConnectionProperties{
			RemoteVpnSite:      &network.SubResource{ID: utils.String(d.Get("remote_vpn_site_id").(string))},
			VpnLinkConnections: expandArmVpnGatewayConnectionVpnSiteLinkConnections(d.Get("vpn_link_connection").([]interface{})),
			// TODO:
			//RoutingConfiguration:           &network.RoutingConfiguration{...},
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

	id, err := parse.VPNGatewayConnectionID(*resp.ID)
	if err != nil {
		return err
	}
	d.SetId(id.ID(subscriptionId))

	return resourceArmVpnGatewayConnectionResourceRead(d, meta)
}

func resourceArmVpnGatewayConnectionResourceRead(d *schema.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Network.VpnConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VPNGatewayConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Gateway, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Vpn Gateway Connection Resource %q was not found in VPN Gateway %q in Resource Group %q - removing from state!", id.Name, id.Gateway, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", id.Name, id.ResourceGroup, id.Gateway, err)
	}

	d.Set("name", id.Name)

	gatewayId := parse.NewVPNGatewayID(id.ResourceGroup, id.Gateway)
	d.Set("vpn_gateway_id", gatewayId.ID(subscriptionId))

	if prop := resp.VpnConnectionProperties; prop != nil {
		vpnSiteId := ""
		if site := prop.RemoteVpnSite; site != nil {
			if id := site.ID; id != nil {
				theVpnSiteId, err := parse.VpnSiteID(*id)
				if err != nil {
					return err
				}
				vpnSiteId = theVpnSiteId.ID(subscriptionId)
			}
		}

		d.Set("remote_vpn_site_id", vpnSiteId)

		if err := d.Set("routing_configuration", flattenArmVpnGatewayConnectionRoutingConfiguration(prop.RoutingConfiguration)); err != nil {
			return fmt.Errorf(`setting "routing_configuration": %v`, err)
		}

		if err := d.Set("vpn_link_connection", flattenArmVpnGatewayConnectionVpnSiteLinkConnections(prop.VpnLinkConnections)); err != nil {
			return fmt.Errorf(`setting "vpn_link_connection": %v`, err)
		}
	}

	return nil
}

func resourceArmVpnGatewayConnectionResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VpnConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VPNGatewayConnectionID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.Gateway, id.Name); err != nil {
		return fmt.Errorf("deleting Vpn Gateway Connection Resource %q (Resource Group %q / VPN Gateway %q): %+v", id.Name, id.ResourceGroup, id.Gateway, err)
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
				RoutingWeight:                  utils.Int32(int32(e["route_weight"].(int))),
				VpnConnectionProtocolType:      network.VirtualNetworkGatewayConnectionProtocol(e["protocol"].(string)),
				ConnectionBandwidth:            utils.Int32(int32(e["bandwidth_mbps"].(int))),
				EnableBgp:                      utils.Bool(e["bgp_enabled"].(bool)),
				UsePolicyBasedTrafficSelectors: utils.Bool(e["use_policy_based_traffic_selector"].(bool)),
				IpsecPolicies:                  expandArmVpnGatewayConnectionIpSecPolicies(e["ipsec_policy"].([]interface{})),
				EnableRateLimiting:             utils.Bool(e["ratelimit_enabled"].(bool)),
				UseLocalAzureIPAddress:         utils.Bool(e["use_local_azure_ip_address"].(bool)),
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
			"name":                              name,
			"route_weight":                      routeWeight,
			"protocol":                          string(e.VpnConnectionProtocolType),
			"bandwidth_mbps":                    bandwidth,
			"shared_key":                        sharedKey,
			"bgp_enabled":                       bgpEnabled,
			"use_policy_based_traffic_selector": usePolicyBased,
			"ipsec_policy":                      flattenArmVpnGatewayConnectionIpSecPolicies(e.IpsecPolicies),
			"ratelimit_enabled":                 rateLimitEnabled,
			"use_local_azure_ip_address":        useLocalAzureIpAddress,
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
			IpsecEncryption:     network.IpsecEncryption(e["ipsec_encryption_algorithm"].(string)),
			IpsecIntegrity:      network.IpsecIntegrity(e["ipsec_integrity_algorithm"].(string)),
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
			"sa_lifetime_sec":            saLifetimeSec,
			"sa_data_size_kb":            saDataSizeKb,
			"ipsec_encryption_algorithm": string(e.IpsecEncryption),
			"ipsec_integrity_algorithm":  string(e.IpsecIntegrity),
			"ike_encryption_algorithm":   string(e.IkeEncryption),
			"ike_integrity_algorithm":    string(e.IkeIntegrity),
			"dh_group":                   string(e.DhGroup),
			"pfs_group":                  string(e.PfsGroup),
		}

		output = append(output, v)
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
