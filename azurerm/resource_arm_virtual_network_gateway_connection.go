package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmVirtualNetworkGatewayConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkGatewayConnectionCreateUpdate,
		Read:   resourceArmVirtualNetworkGatewayConnectionRead,
		Update: resourceArmVirtualNetworkGatewayConnectionCreateUpdate,
		Delete: resourceArmVirtualNetworkGatewayConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
					string(network.ExpressRoute),
					string(network.IPsec),
					string(network.Vnet2Vnet),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"virtual_network_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"authorization_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"express_route_circuit_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"peer_virtual_network_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"local_network_gateway_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enable_bgp": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"use_policy_based_traffic_selectors": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"routing_weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"ipsec_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dh_group": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.DHGroup1),
								string(network.DHGroup14),
								string(network.DHGroup2),
								string(network.DHGroup2048),
								string(network.DHGroup24),
								string(network.ECP256),
								string(network.ECP384),
								string(network.None),
							}, true),
						},
						"ike_encryption": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.AES128),
								string(network.AES192),
								string(network.AES256),
								string(network.DES),
								string(network.DES3),
							}, true),
						},
						"ike_integrity": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.MD5),
								string(network.SHA1),
								string(network.SHA256),
								string(network.SHA384),
							}, true),
						},
						"ipsec_encryption": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.PfsGroupECP256),
								string(network.PfsGroupECP384),
								string(network.PfsGroupNone),
								string(network.PfsGroupPFS1),
								string(network.PfsGroupPFS2),
								string(network.PfsGroupPFS2048),
								string(network.PfsGroupPFS24),
							}, true),
						},
						"sa_datasize": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(1024),
						},
						"sa_lifetime": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntAtLeast(300),
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkGatewayConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Virtual Network Gateway Connection creation.")

	name := d.Get("name").(string)
	location := d.Get("location").(string)
	resGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	properties, err := getArmVirtualNetworkGatewayConnectionProperties(d)
	if err != nil {
		return err
	}

	connection := network.VirtualNetworkGatewayConnection{
		Name:     &name,
		Location: &location,
		Tags:     expandTags(tags),
		VirtualNetworkGatewayConnectionPropertiesFormat: properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, name, connection)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating AzureRM Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Virtual Network Gateway Connection %q (resource group %q) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceArmVirtualNetworkGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

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
		d.Set("location", azureRMNormalizeLocation(*location))
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

	if conn.Peer != nil {
		d.Set("express_route_circuit_id", conn.Peer.ID)
	}

	if conn.VirtualNetworkGateway2 != nil {
		d.Set("peer_virtual_network_gateway_id", conn.VirtualNetworkGateway2.ID)
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

	if conn.IpsecPolicies != nil {
		ipsec_policies := flattenArmVirtualNetworkGatewayConnectionIpsecPolicies(conn.IpsecPolicies)
		d.Set("ipsec_policy", ipsec_policies)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualNetworkGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetGatewayConnectionsClient
	ctx := meta.(*ArmClient).StopContext

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error Deleting Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Gateway Connection %q (Resource Group %q): %+v", name, resGroup, err)
	}

	return nil
}

func getArmVirtualNetworkGatewayConnectionProperties(d *schema.ResourceData) (*network.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := network.VirtualNetworkGatewayConnectionType(d.Get("type").(string))

	props := &network.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType: connectionType,
	}

	if v, ok := d.GetOk("virtual_network_gateway_id"); ok {
		virtualNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndVirtualNetworkGatewayFromId(virtualNetworkGatewayId)
		if err != nil {
			return nil, errwrap.Wrapf("Error Getting VirtualNetworkGateway Name and Group: {{err}}", err)
		}

		props.VirtualNetworkGateway1 = &network.VirtualNetworkGateway{
			ID:   &virtualNetworkGatewayId,
			Name: &name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("authorization_key"); ok {
		authorizationKey := v.(string)
		props.AuthorizationKey = &authorizationKey
	}

	if v, ok := d.GetOk("express_route_circuit_id"); ok {
		expressRouteCircuitId := v.(string)
		props.Peer = &network.SubResource{
			ID: &expressRouteCircuitId,
		}
	}

	if v, ok := d.GetOk("peer_virtual_network_gateway_id"); ok {
		peerVirtualNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndVirtualNetworkGatewayFromId(peerVirtualNetworkGatewayId)
		if err != nil {
			return nil, errwrap.Wrapf("Error Getting VirtualNetworkGateway Name and Group: {{err}}", err)
		}

		props.VirtualNetworkGateway2 = &network.VirtualNetworkGateway{
			ID:   &peerVirtualNetworkGatewayId,
			Name: &name,
			VirtualNetworkGatewayPropertiesFormat: &network.VirtualNetworkGatewayPropertiesFormat{
				IPConfigurations: &[]network.VirtualNetworkGatewayIPConfiguration{},
			},
		}
	}

	if v, ok := d.GetOk("local_network_gateway_id"); ok {
		localNetworkGatewayId := v.(string)
		_, name, err := resourceGroupAndLocalNetworkGatewayFromId(localNetworkGatewayId)
		if err != nil {
			return nil, errwrap.Wrapf("Error Getting LocalNetworkGateway Name and Group: {{err}}", err)
		}

		props.LocalNetworkGateway2 = &network.LocalNetworkGateway{
			ID:   &localNetworkGatewayId,
			Name: &name,
			LocalNetworkGatewayPropertiesFormat: &network.LocalNetworkGatewayPropertiesFormat{
				LocalNetworkAddressSpace: &network.AddressSpace{},
			},
		}
	}

	if v, ok := d.GetOk("enable_bgp"); ok {
		props.EnableBgp = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("use_policy_based_traffic_selectors"); ok {
		props.UsePolicyBasedTrafficSelectors = utils.Bool(v.(bool))
	}

	if v, ok := d.GetOk("routing_weight"); ok {
		routing_weight := int32(v.(int))
		props.RoutingWeight = &routing_weight
	}

	if v, ok := d.GetOk("shared_key"); ok {
		props.SharedKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("ipsec_policy"); ok {
		ipsec_policies := v.([]interface{})
		props.IpsecPolicies = expandArmVirtualNetworkGatewayConnectionIpsecPolicies(ipsec_policies)
		log.Printf("Ipsec policies: %+q", props.IpsecPolicies)
	}

	if props.ConnectionType == network.ExpressRoute {
		if props.Peer == nil || props.Peer.ID == nil {
			return nil, fmt.Errorf("`express_route_circuit_id` must be specified when `type` is set to `ExpressRoute")
		}
	}

	if props.ConnectionType == network.IPsec {
		if props.LocalNetworkGateway2 == nil || props.LocalNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` and `shared_key` must be specified when `type` is set to `IPsec")
		}

		if props.SharedKey == nil {
			return nil, fmt.Errorf("`local_network_gateway_id` and `shared_key` must be specified when `type` is set to `IPsec")
		}
	}

	if props.ConnectionType == network.Vnet2Vnet {
		if props.VirtualNetworkGateway2 == nil || props.VirtualNetworkGateway2.ID == nil {
			return nil, fmt.Errorf("`peer_virtual_network_gateway_id` and `shared_key` must be specified when `type` is set to `Vnet2Vnet")
		}
	}

	return props, nil
}

func resourceGroupAndVirtualNetworkGatewayConnectionFromId(virtualNetworkGatewayConnectionId string) (string, string, error) {
	id, err := parseAzureResourceID(virtualNetworkGatewayConnectionId)
	if err != nil {
		return "", "", err
	}
	name := id.Path["connections"]
	resGroup := id.ResourceGroup

	return resGroup, name, nil
}

func expandArmVirtualNetworkGatewayConnectionIpsecPolicies(ipsec_policies []interface{}) *[]network.IpsecPolicy {
	ipsecPolicies := make([]network.IpsecPolicy, 0, len(ipsec_policies))

	for _, d := range ipsec_policies {
		ipsec_policy := d.(map[string]interface{})
		ipsecPolicy := &network.IpsecPolicy{}

		if dh_group, ok := ipsec_policy["dh_group"].(string); ok && dh_group != "" {
			ipsecPolicy.DhGroup = network.DhGroup(dh_group)
		}

		if ike_encryption, ok := ipsec_policy["ike_encryption"].(string); ok && ike_encryption != "" {
			ipsecPolicy.IkeEncryption = network.IkeEncryption(ike_encryption)
		}

		if ike_integrity, ok := ipsec_policy["ike_integrity"].(string); ok && ike_integrity != "" {
			ipsecPolicy.IkeIntegrity = network.IkeIntegrity(ike_integrity)
		}

		if ipsec_encryption, ok := ipsec_policy["ipsec_encryption"].(string); ok && ipsec_encryption != "" {
			ipsecPolicy.IpsecEncryption = network.IpsecEncryption(ipsec_encryption)
		}

		if ipsec_integrity, ok := ipsec_policy["ipsec_integrity"].(string); ok && ipsec_integrity != "" {
			ipsecPolicy.IpsecIntegrity = network.IpsecIntegrity(ipsec_integrity)
		}

		if pfs_group, ok := ipsec_policy["pfs_group"].(string); ok && pfs_group != "" {
			ipsecPolicy.PfsGroup = network.PfsGroup(pfs_group)
		}

		if v, ok := ipsec_policy["sa_datasize"].(int); ok {
			sa_datasize := int32(v)
			ipsecPolicy.SaDataSizeKilobytes = &sa_datasize
		}

		if v, ok := ipsec_policy["sa_lifetime"].(int); ok {
			sa_lifetime := int32(v)
			ipsecPolicy.SaLifeTimeSeconds = &sa_lifetime
		}

		ipsecPolicies = append(ipsecPolicies, *ipsecPolicy)
	}

	return &ipsecPolicies
}

func flattenArmVirtualNetworkGatewayConnectionIpsecPolicies(ipsecPolicies *[]network.IpsecPolicy) []interface{} {
	ipsec_policies := make([]interface{}, 0, len(*ipsecPolicies))

	for _, ipsecPolicy := range *ipsecPolicies {
		ipsec_policy := make(map[string]interface{})

		ipsec_policy["dh_group"] = string(ipsecPolicy.DhGroup)
		ipsec_policy["ike_encryption"] = string(ipsecPolicy.IkeEncryption)
		ipsec_policy["ike_integrity"] = string(ipsecPolicy.IkeIntegrity)
		ipsec_policy["ipsec_encryption"] = string(ipsecPolicy.IpsecEncryption)
		ipsec_policy["ipsec_integrity"] = string(ipsecPolicy.IpsecIntegrity)
		ipsec_policy["pfs_group"] = string(ipsecPolicy.PfsGroup)

		if ipsecPolicy.SaDataSizeKilobytes != nil {
			ipsec_policy["sa_datasize"] = int(*ipsecPolicy.SaDataSizeKilobytes)
		}

		if ipsecPolicy.SaLifeTimeSeconds != nil {
			ipsec_policy["sa_lifetime"] = int(*ipsecPolicy.SaLifeTimeSeconds)
		}

		ipsec_policies = append(ipsec_policies, ipsec_policy)
	}

	return ipsec_policies
}
