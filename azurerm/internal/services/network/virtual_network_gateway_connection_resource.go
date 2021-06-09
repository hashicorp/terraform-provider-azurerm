package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceVirtualNetworkGatewayConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkGatewayConnectionCreateUpdate,
		Read:   resourceVirtualNetworkGatewayConnectionRead,
		Update: resourceVirtualNetworkGatewayConnectionCreateUpdate,
		Delete: resourceVirtualNetworkGatewayConnectionDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
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

			"traffic_selector_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
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

			"ipsec_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"dh_group": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.DhGroupDHGroup1),
								string(network.DhGroupDHGroup14),
								string(network.DhGroupDHGroup2),
								string(network.DhGroupDHGroup2048),
								string(network.DhGroupDHGroup24),
								string(network.DhGroupECP256),
								string(network.DhGroupECP384),
								string(network.DhGroupNone),
							}, true),
						},

						"ike_encryption": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IkeEncryptionAES128),
								string(network.IkeEncryptionAES192),
								string(network.IkeEncryptionAES256),
								string(network.IkeEncryptionDES),
								string(network.IkeEncryptionDES3),
								string(network.IkeEncryptionGCMAES128),
								string(network.IkeEncryptionGCMAES256),
							}, true),
						},

						"ike_integrity": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IkeIntegrityGCMAES128),
								string(network.IkeIntegrityGCMAES256),
								string(network.IkeIntegrityMD5),
								string(network.IkeIntegritySHA1),
								string(network.IkeIntegritySHA256),
								string(network.IkeIntegritySHA384),
							}, true),
						},

						"ipsec_encryption": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
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
							}, true),
						},

						"ipsec_integrity": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IpsecIntegrityGCMAES128),
								string(network.IpsecIntegrityGCMAES192),
								string(network.IpsecIntegrityGCMAES256),
								string(network.IpsecIntegrityMD5),
								string(network.IpsecIntegritySHA1),
								string(network.IpsecIntegritySHA256),
							}, true),
						},

						"pfs_group": {
							Type:             pluginsdk.TypeString,
							Required:         true,
							DiffSuppressFunc: suppress.CaseDifference,
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
							}, true),
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway Connection creation.")

	name := d.Get("name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Virtual Network Gateway Connection %q (Resource Group %q): %s", name, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_network_gateway_connection", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	properties, err := getVirtualNetworkGatewayConnectionProperties(d)
	if err != nil {
		return err
	}

	connection := network.VirtualNetworkGatewayConnection{
		Name:     &name,
		Location: &location,
		Tags:     tags.Expand(t),
		VirtualNetworkGatewayConnectionPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, connection)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AzureRM Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if properties.SharedKey != nil && !d.IsNewResource() {
		future, err := client.SetSharedKey(ctx, resGroup, name, network.ConnectionSharedKey{
			Value: properties.SharedKey,
		})
		if err != nil {
			return fmt.Errorf("Updating Shared Key for Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
		}
		if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("Waiting for updating Shared Key for Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
		}
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Virtual Network Gateway Connection %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceVirtualNetworkGatewayConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway Connection %q: %+v", name, err)
	}

	conn := *resp.VirtualNetworkGatewayConnectionPropertiesFormat

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
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

	d.Set("connection_protocol", string(conn.ConnectionProtocol))

	if conn.ExpressRouteGatewayBypass != nil {
		d.Set("express_route_gateway_bypass", conn.ExpressRouteGatewayBypass)
	}

	if conn.IpsecPolicies != nil {
		ipsecPolicies := flattenVirtualNetworkGatewayConnectionIpsecPolicies(conn.IpsecPolicies)

		if err := d.Set("ipsec_policy", ipsecPolicies); err != nil {
			return fmt.Errorf("Error setting `ipsec_policy`: %+v", err)
		}
	}

	trafficSelectorPolicies := flattenVirtualNetworkGatewayConnectionTrafficSelectorPolicies(conn.TrafficSelectorPolicies)
	if err := d.Set("traffic_selector_policy", trafficSelectorPolicies); err != nil {
		return fmt.Errorf("Error setting `traffic_selector_policy`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceVirtualNetworkGatewayConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetGatewayConnectionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Deleting Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getVirtualNetworkGatewayConnectionProperties(d *pluginsdk.ResourceData) (*network.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := network.VirtualNetworkGatewayConnectionType(d.Get("type").(string))

	props := &network.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType:                 connectionType,
		EnableBgp:                      utils.Bool(d.Get("enable_bgp").(bool)),
		ExpressRouteGatewayBypass:      utils.Bool(d.Get("express_route_gateway_bypass").(bool)),
		UsePolicyBasedTrafficSelectors: utils.Bool(d.Get("use_policy_based_traffic_selectors").(bool)),
	}

	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)

		gwid, err := parse.VirtualNetworkGatewayID(virtualNetworkGatewayId)
		if err != nil {
			return nil, err
		}

		props.VirtualNetworkGateway1 = &network.VirtualNetworkGateway{
			ID:   &virtualNetworkGatewayId,
			Name: &gwid.Name,
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
			return nil, fmt.Errorf("Error Getting LocalNetworkGateway Name and Group:: %+v", err)
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

	return props, nil
}

func resourceGroupAndVirtualNetworkGatewayConnectionFromId(virtualNetworkGatewayConnectionId string) (string, string, error) {
	id, err := azure.ParseAzureResourceID(virtualNetworkGatewayConnectionId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["connections"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
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
