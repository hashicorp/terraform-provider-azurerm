package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualNetworkGatewayConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayConnectionCreateUpdate,
		Read:   resourceVirtualNetworkGatewayConnectionRead,
		Update: resourceVirtualNetworkGatewayConnectionCreateUpdate,
		Delete: resourceVirtualNetworkGatewayConnectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.NetworkGatewayConnectionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayConnectionTypeExpressRoute),
					string(network.VirtualNetworkGatewayConnectionTypeIPsec),
					string(network.VirtualNetworkGatewayConnectionTypeVnet2Vnet),
				}, false),
			},

			"virtual_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"authorization_key": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"dpd_timeout_seconds": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"express_route_circuit_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"egress_nat_rule_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.VirtualNetworkGatewayNatRuleID,
				},
			},

			"ingress_nat_rule_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.VirtualNetworkGatewayNatRuleID,
				},
			},

			"peer_virtual_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			"local_azure_ip_address_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"local_network_gateway_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: azure.ValidateResourceIDOrEmpty,
			},

			// TODO 4.0: change this from enable_* to *_enabled
			"enable_bgp": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"use_policy_based_traffic_selectors": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"routing_weight": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 32000),
			},

			"shared_key": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"express_route_gateway_bypass": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"connection_protocol": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayConnectionProtocolIKEv1),
					string(network.VirtualNetworkGatewayConnectionProtocolIKEv2),
				}, false),
			},

			"connection_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.VirtualNetworkGatewayConnectionModeInitiatorOnly),
					string(network.VirtualNetworkGatewayConnectionModeResponderOnly),
					string(network.VirtualNetworkGatewayConnectionModeDefault),
				}, false),
				Default: string(network.VirtualNetworkGatewayConnectionModeDefault),
			},

			"traffic_selector_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"remote_address_cidrs": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"custom_bgp_addresses": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"primary": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.Any(validation.IsIPv4Address),
						},
						"secondary": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.Any(validation.IsIPv4Address),
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

						"sa_datasize": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1024),
						},

						"sa_lifetime": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(300),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceVirtualNetworkGatewayConnectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	vnetGatewayClient := meta.(*clients.Client).Network.VnetGatewayClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway Connection creation.")

	id := parse.NewNetworkGatewayConnectionID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ConnectionName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway_connection", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	var virtualNetworkGateway network.VirtualNetworkGateway
	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := parse.VirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return err
		}

		virtualNetworkGateway, err = vnetGatewayClient.Get(ctx, id.ResourceGroup, gwid.Name)
		if err != nil {
			return err
		}
	}

	properties, err := getVirtualNetworkGatewayConnectionProperties(d, virtualNetworkGateway)
	if err != nil {
		return err
	}

	connection := network.VirtualNetworkGatewayConnection{
		Name:     &id.ConnectionName,
		Location: &location,
		Tags:     tags.Expand(t),
		VirtualNetworkGatewayConnectionPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ConnectionName, connection)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for completion of %s: %+v", id, err)
	}

	if properties.SharedKey != nil && !d.IsNewResource() {
		future, err := client.SetSharedKey(ctx, id.ResourceGroup, id.ConnectionName, network.ConnectionSharedKey{
			Value: properties.SharedKey,
		})
		if err != nil {
			return fmt.Errorf("updating Shared Key for %s: %+v", id, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for updating Shared Key for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkGatewayConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ConnectionName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	conn := *resp.VirtualNetworkGatewayConnectionPropertiesFormat

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if string(conn.ConnectionType) != "" {
		d.Set("type", string(conn.ConnectionType))
	}

	if conn.VirtualNetworkGateway1 != nil {
		d.Set("virtual_network_gateway_id", conn.VirtualNetworkGateway1.ID)
	}

	if conn.AuthorizationKey != nil {
		d.Set("authorization_key", conn.AuthorizationKey)
	}

	if conn.DpdTimeoutSeconds != nil {
		d.Set("dpd_timeout_seconds", conn.DpdTimeoutSeconds)
	}

	if conn.Peer != nil {
		d.Set("express_route_circuit_id", conn.Peer.ID)
	}

	if conn.VirtualNetworkGateway2 != nil {
		d.Set("peer_virtual_network_gateway_id", conn.VirtualNetworkGateway2.ID)
	}

	if conn.UseLocalAzureIPAddress != nil {
		d.Set("local_azure_ip_address_enabled", conn.UseLocalAzureIPAddress)
	}

	if conn.LocalNetworkGateway2 != nil {
		d.Set("local_network_gateway_id", conn.LocalNetworkGateway2.ID)
	}

	if conn.EnableBgp != nil {
		d.Set("enable_bgp", conn.EnableBgp)
	}

	if conn.UsePolicyBasedTrafficSelectors != nil {
		d.Set("use_policy_based_traffic_selectors", conn.UsePolicyBasedTrafficSelectors)
	}

	if conn.RoutingWeight != nil {
		d.Set("routing_weight", conn.RoutingWeight)
	}

	if conn.SharedKey != nil {
		d.Set("shared_key", conn.SharedKey)
	}

	if conn.GatewayCustomBgpIPAddresses != nil {
		adresses := flattenGatewayCustomBgpIPAddresses(conn.GatewayCustomBgpIPAddresses)
		if err := d.Set("custom_bgp_addresses", adresses); err != nil {
			return fmt.Errorf("setting `custom_bgp_addresses`: %+v", err)
		}
	}

	d.Set("connection_protocol", string(conn.ConnectionProtocol))

	d.Set("connection_mode", string(conn.ConnectionMode))

	if conn.ExpressRouteGatewayBypass != nil {
		d.Set("express_route_gateway_bypass", conn.ExpressRouteGatewayBypass)
	}

	if conn.IpsecPolicies != nil {
		ipsecPolicies := flattenVirtualNetworkGatewayConnectionIpsecPolicies(conn.IpsecPolicies)

		if err := d.Set("ipsec_policy", ipsecPolicies); err != nil {
			return fmt.Errorf("setting `ipsec_policy`: %+v", err)
		}
	}

	trafficSelectorPolicies := flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(conn.TrafficSelectorPolicies)
	if err := d.Set("traffic_selector_policy", trafficSelectorPolicies); err != nil {
		return fmt.Errorf("setting `traffic_selector_policy`: %+v", err)
	}

	if err := d.Set("egress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(conn.EgressNatRules)); err != nil {
		return fmt.Errorf("setting `egress_nat_rule_ids`: %+v", err)
	}

	if err := d.Set("ingress_nat_rule_ids", flattenVirtualNetworkGatewayConnectionNatRuleIds(conn.IngressNatRules)); err != nil {
		return fmt.Errorf("setting `ingress_nat_rule_ids`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualNetworkGatewayConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.NetworkGatewayConnectionID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ConnectionName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", id, err)
	}

	return nil
}

func getVirtualNetworkGatewayConnectionProperties(d *pluginsdk.ResourceData, virtualNetworkGateway network.VirtualNetworkGateway) (*network.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := network.VirtualNetworkGatewayConnectionType(d.Get("type").(string))
	connectionMode := network.VirtualNetworkGatewayConnectionMode(d.Get("connection_mode").(string))

	props := &network.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType:                 connectionType,
		ConnectionMode:                 connectionMode,
		EnableBgp:                      utils.Bool(d.Get("enable_bgp").(bool)),
		ExpressRouteGatewayBypass:      utils.Bool(d.Get("express_route_gateway_bypass").(bool)),
		UsePolicyBasedTrafficSelectors: utils.Bool(d.Get("use_policy_based_traffic_selectors").(bool)),
	}

	if virtualNetworkGateway.Name != nil && virtualNetworkGateway.ID != nil {
		props.VirtualNetworkGateway1 = &network.VirtualNetworkGateway{
			ID:   virtualNetworkGateway.ID,
			Name: virtualNetworkGateway.Name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		authorizationKey := v.(string)
		props.AuthorizationKey = &authorizationKey
	}

	if v, ok := d.GetOk("dpd_timeout_seconds"); ok {
		props.DpdTimeoutSeconds = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("express_route_circuit_id"); ok {
		expressRouteCircuitId := v.(string)
		props.Peer = &network.SubResource{
			ID: &expressRouteCircuitId,
		}
	}

	if v, ok := d.GetOk("egress_nat_rule_ids"); ok {
		props.EgressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("ingress_nat_rule_ids"); ok {
		props.IngressNatRules = expandVirtualNetworkGatewayConnectionNatRuleIds(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("peer_virtual_network_gateway_id"); ok {
		gwid, err := parse.VirtualNetworkGatewayID(v.(string))
		if err != nil {
			return nil, err
		}
		props.VirtualNetworkGateway2 = &network.VirtualNetworkGateway{
			ID:   utils.String(gwid.ID()),
			Name: &gwid.Name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("local_azure_ip_address_enabled"); ok {
		props.UseLocalAzureIPAddress = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("local_network_gateway_id"); ok {
		localNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return nil, fmt.Errorf("Getting LocalNetworkGateway Name and Group:: %+v", err)
		}

		props.LocalNetworkGateway2 = &network.LocalNetworkGateway{
			ID:   &localNetworkGatewayId,
			Name: &name,
			LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &network.AddressSpace{},
			},
		}
	}

	if v, ok := d.GetOk("routing_weight"); ok {
		routingWeight := int32(v.(int))
		props.RoutingWeight = &routingWeight
	}

	if v, ok := d.GetOk("shared_key"); ok {
		props.SharedKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("connection_protocol"); ok {
		connectionProtocol := v.(string)
		props.ConnectionProtocol = network.VirtualNetworkGatewayConnectionProtocol(connectionProtocol)
	}

	if v, ok := d.GetOk("traffic_selector_policy"); ok {
		props.TrafficSelectorPolicies = expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(v.([]interface{}))
	}

	if v, ok := d.GetOk("ipsec_policy"); ok {
		props.IpsecPolicies = expandVirtualNetworkGatewayConnectionIpsecPolicies(v.([]interface{}))
	}

	if utils.NormaliseNilableBool(props.EnableBgp) {
		if _, ok := d.GetOk("custom_bgp_addresses"); ok {
			gatewayCustomBgpIPAddresses, err := expandGatewayCustomBgpIPAddresses(d, virtualNetworkGateway.VirtualNetworkGatewayPropertiesFormat.BgpSettings.BgpPeeringAddresses)
			if err != nil {
				return nil, err
			}

			props.GatewayCustomBgpIPAddresses = gatewayCustomBgpIPAddresses
		}
	}

	if props.ConnectionType == network.VirtualNetworkGatewayConnectionTypeExpressRoute {
		if props.Peer == nil || props.Peer.ID == nil {
			return nil, fmt.Errorf("`express_route_circuit_id` must be specified when `type` is set to `ExpressRoute`")
		}
	}

	if props.ConnectionType == network.VirtualNetworkGatewayConnectionTypeIPsec {
		if props.LocalNetworkGateway2 == nil || props.LocalNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` must be specified when `type` is set to `IPsec`")
		}
	}

	if props.ConnectionType == network.VirtualNetworkGatewayConnectionTypeVnet2Vnet {
		if props.VirtualNetworkGateway2 == nil || props.VirtualNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`peer_virtual_network_gateway_id` must be specified when `type` is set to `Vnet2Vnet`")
		}
	}

	if props.GatewayCustomBgpIPAddresses != nil && props.ConnectionType != network.VirtualNetworkGatewayConnectionTypeIPsec {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `type` is set to `IPsec`")
	}

	if props.GatewayCustomBgpIPAddresses != nil && virtualNetworkGateway.VirtualNetworkGatewayPropertiesFormat.ActiveActive == utils.Bool(false) {
		return nil, fmt.Errorf("`custom_bgp_addresses` can only be used when `azurerm_virtual_network_gateway` `active_active` is set enabled`")
	}

	return props, nil
}

func expandVirtualNetworkGatewayConnectionIpsecPolicies(schemaIpsecPolicies []interface{}) *[]network.IpsecPolicy {
	ipsecPolicies := make([]network.IpsecPolicy, 0, len(schemaIpsecPolicies))

	for _, d := range schemaIpsecPolicies {
		schemaIpsecPolicy := d.(map[string]interface{})
		ipsecPolicy := &network.IpsecPolicy{}

		if dhGroup, ok := schemaIpsecPolicy["dh_group"].(string); ok && dhGroup != "" {
			ipsecPolicy.DhGroup = network.DhGroup(dhGroup)
		}

		if ikeEncryption, ok := schemaIpsecPolicy["ike_encryption"].(string); ok && ikeEncryption != "" {
			ipsecPolicy.IkeEncryption = network.IkeEncryption(ikeEncryption)
		}

		if ikeIntegrity, ok := schemaIpsecPolicy["ike_integrity"].(string); ok && ikeIntegrity != "" {
			ipsecPolicy.IkeIntegrity = network.IkeIntegrity(ikeIntegrity)
		}

		if ipsecEncryption, ok := schemaIpsecPolicy["ipsec_encryption"].(string); ok && ipsecEncryption != "" {
			ipsecPolicy.IpsecEncryption = network.IpsecEncryption(ipsecEncryption)
		}

		if ipsecIntegrity, ok := schemaIpsecPolicy["ipsec_integrity"].(string); ok && ipsecIntegrity != "" {
			ipsecPolicy.IpsecIntegrity = network.IpsecIntegrity(ipsecIntegrity)
		}

		if pfsGroup, ok := schemaIpsecPolicy["pfs_group"].(string); ok && pfsGroup != "" {
			ipsecPolicy.PfsGroup = network.PfsGroup(pfsGroup)
		}

		if v, ok := schemaIpsecPolicy["sa_datasize"].(int); ok {
			saDatasize := int32(v)
			ipsecPolicy.SaDataSizeKilobytes = &saDatasize
		}

		if v, ok := schemaIpsecPolicy["sa_lifetime"].(int); ok {
			saLifetime := int32(v)
			ipsecPolicy.SaLifeTimeSeconds = &saLifetime
		}

		ipsecPolicies = append(ipsecPolicies, *ipsecPolicy)
	}

	return &ipsecPolicies
}

func expandVirtualNetworkGatewayConnectionTrafficSelectorPolicies(schemaTrafficSelectorPolicies []interface{}) *[]network.TrafficSelectorPolicy {
	trafficSelectorPolicies := make([]network.TrafficSelectorPolicy, 0, len(schemaTrafficSelectorPolicies))

	for _, d := range schemaTrafficSelectorPolicies {
		schemaTrafficSelectorPolicy := d.(map[string]interface{})
		trafficSelectorPolicy := &network.TrafficSelectorPolicy{}
		if localAddressRanges, ok := schemaTrafficSelectorPolicy["local_address_cidrs"].([]interface{}); ok {
			trafficSelectorPolicy.LocalAddressRanges = utils.ExpandStringSlice(localAddressRanges)
		}
		if remoteAddressRanges, ok := schemaTrafficSelectorPolicy["remote_address_cidrs"].([]interface{}); ok {
			trafficSelectorPolicy.RemoteAddressRanges = utils.ExpandStringSlice(remoteAddressRanges)
		}

		trafficSelectorPolicies = append(trafficSelectorPolicies, *trafficSelectorPolicy)
	}

	return &trafficSelectorPolicies
}

func expandGatewayCustomBgpIPAddresses(d *pluginsdk.ResourceData, bgpPeeringAddresses *[]network.IPConfigurationBgpPeeringAddress) (*[]network.GatewayCustomBgpIPAddressIPConfiguration, error) {
	customBgpIpAddresses := make([]network.GatewayCustomBgpIPAddressIPConfiguration, 0)

	bgpAddresses := d.Get("custom_bgp_addresses").([]interface{})
	if len(bgpAddresses) == 0 {
		return &customBgpIpAddresses, nil
	}

	bgAs := bgpAddresses[0].(map[string]interface{})
	primaryAddress := bgAs["primary"].(string)
	secondaryAddress := bgAs["secondary"].(string)

	var primaryIpConfiguration *string
	var secondaryIpConfiguration *string

	for _, address := range *bgpPeeringAddresses {
		for _, ip := range *address.CustomBgpIPAddresses {
			if ip == primaryAddress {
				primaryIpConfiguration = address.IpconfigurationID
			} else if ip == secondaryAddress {
				secondaryIpConfiguration = address.IpconfigurationID
			}
		}
	}

	if len(*primaryIpConfiguration) == 0 || len(*secondaryIpConfiguration) == 0 {
		return &customBgpIpAddresses, fmt.Errorf("primary or secondary address not found at `virtual_network_gateway` configuration `bgp_settings` `peering_addresses`")
	}

	customBgpIpAddresses = append(customBgpIpAddresses, network.GatewayCustomBgpIPAddressIPConfiguration{
		IPConfigurationID:  primaryIpConfiguration,
		CustomBgpIPAddress: utils.String(primaryAddress),
	}, network.GatewayCustomBgpIPAddressIPConfiguration{
		IPConfigurationID:  secondaryIpConfiguration,
		CustomBgpIPAddress: utils.String(secondaryAddress),
	})

	return &customBgpIpAddresses, nil
}

func flattenVirtualNetworkGatewayConnectionIpsecPolicies(ipsecPolicies *[]network.IpsecPolicy) []interface{} {
	schemaIpsecPolicies := make([]interface{}, 0)

	if ipsecPolicies != nil {
		for _, ipsecPolicy := range *ipsecPolicies {
			schemaIpsecPolicy := make(map[string]interface{})

			schemaIpsecPolicy["dh_group"] = string(ipsecPolicy.DhGroup)
			schemaIpsecPolicy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
			schemaIpsecPolicy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
			schemaIpsecPolicy["ipsec_encryption"] = string(ipsecPolicy.IpsecEncryption)
			schemaIpsecPolicy["ipsec_integrity"] = string(ipsecPolicy.IpsecIntegrity)
			schemaIpsecPolicy["pfs_group"] = string(ipsecPolicy.PfsGroup)

			if ipsecPolicy.SaDataSizeKilobytes != nil {
				schemaIpsecPolicy["sa_datasize"] = int(*ipsecPolicy.SaDataSizeKilobytes)
			}

			if ipsecPolicy.SaLifeTimeSeconds != nil {
				schemaIpsecPolicy["sa_lifetime"] = int(*ipsecPolicy.SaLifeTimeSeconds)
			}

			schemaIpsecPolicies = append(schemaIpsecPolicies, schemaIpsecPolicy)
		}
	}

	return schemaIpsecPolicies
}

func flattenGatewayCustomBgpIPAddresses(gatewayCustomBgpIPAddresses *[]network.GatewayCustomBgpIPAddressIPConfiguration) interface{} {
	customBgpIpAdresses := make([]interface{}, 0)

	if len(*gatewayCustomBgpIPAddresses) == 2 {
		addresses := *gatewayCustomBgpIPAddresses
		customBgpIpAdresses = append(customBgpIpAdresses, map[string]interface{}{
			"primary":   addresses[0].CustomBgpIPAddress,
			"secondary": addresses[1].CustomBgpIPAddress,
		})
	}

	return customBgpIpAdresses
}

func flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(trafficSelectorPolicies *[]network.TrafficSelectorPolicy) []interface{} {
	schemaTrafficSelectorPolicies := make([]interface{}, 0)

	if trafficSelectorPolicies != nil {
		for _, trafficSelectorPolicy := range *trafficSelectorPolicies {
			schemaTrafficSelectorPolicies = append(schemaTrafficSelectorPolicies, map[string]interface{}{
				"local_address_cidrs":  utils.FlattenStringSlice(trafficSelectorPolicy.LocalAddressRanges),
				"remote_address_cidrs": utils.FlattenStringSlice(trafficSelectorPolicy.RemoteAddressRanges),
			})
		}
	}

	return schemaTrafficSelectorPolicies
}

func expandVirtualNetworkGatewayConnectionNatRuleIds(input []interface{}) *[]network.SubResource {
	results := make([]network.SubResource, 0)

	for _, item := range input {
		results = append(results, network.SubResource{
			ID: utils.String(item.(string)),
		})
	}

	return &results
}

func flattenVirtualNetworkGatewayConnectionNatRuleIds(input *[]network.SubResource) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		var id string
		if item.ID != nil {
			id = *item.ID
		}

		results = append(results, id)
	}

	return results
}
