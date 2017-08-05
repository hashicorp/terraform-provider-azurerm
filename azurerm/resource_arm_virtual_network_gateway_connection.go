package azurerm

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"time"
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

			"type": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
				ValidateFunc: validation.StringInSlice([]string{
					string(network.ExpressRoute),
					string(network.IPsec),
					string(network.Vnet2Vnet),
				}, true),
				ForceNew: true,
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
				Default:  false,
			},

			"routing_weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},

			"shared_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"location": locationSchema(),

			"tags": tagsSchema(),
		},
	}
}

func resourceArmVirtualNetworkGatewayConnectionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	vnetGatewayConnectionsClient := client.vnetGatewayConnectionsClient

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

	_, error := vnetGatewayConnectionsClient.CreateOrUpdate(resGroup, name, connection, make(chan struct{}))
	err = <-error
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Waiting for AzureRM Virtual Network Gateway Connection %s to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Updating"},
		Target:  []string{"Succeeded"},
		Refresh: virtualNetworkGatewayConnectionStateRefreshFunc(client, resGroup, name, false),
		Timeout: 15 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AzureRM Virtual Network Gateway Connection %s to become available: %+v", name, err)
	}

	read, err := vnetGatewayConnectionsClient.Get(resGroup, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read AzureRM Virtual Network Gateway Connection %s (resource group %s) ID", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkGatewayConnectionRead(d, meta)
}

func resourceArmVirtualNetworkGatewayConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	vnetGatewayConnectionsClient := client.vnetGatewayConnectionsClient

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := vnetGatewayConnectionsClient.Get(resGroup, name)
	if err != nil {
		if responseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway Connection %s: %+v", name, err)
	}

	conn := *resp.VirtualNetworkGatewayConnectionPropertiesFormat

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)

	if string(conn.ConnectionType) != "" {
		d.Set("type", string(conn.ConnectionType))
	}

	if conn.VirtualNetworkGateway1 != nil {
		d.Set("virtual_network_gateway_id", conn.VirtualNetworkGateway1.ID)
	}

	d.Set("authorization_key", conn.AuthorizationKey)

	if conn.Peer != nil {
		d.Set("express_route_circuit_id", conn.Peer.ID)
	}

	if conn.VirtualNetworkGateway2 != nil {
		d.Set("peer_virtual_network_gateway_id", conn.VirtualNetworkGateway2.ID)
	}

	if conn.LocalNetworkGateway2 != nil {
		d.Set("local_network_gateway_id", conn.LocalNetworkGateway2.ID)
	}

	d.Set("enable_bgp", conn.EnableBgp)
	d.Set("routing_weight", conn.RoutingWeight)
	d.Set("shared_key", conn.SharedKey)

	d.Set("connection_status", string(conn.ConnectionStatus))
	d.Set("egress_bytes_transferred", conn.EgressBytesTransferred)
	d.Set("ingress_bytes_transferred", conn.EgressBytesTransferred)

	d.Set("location", azureRMNormalizeLocation(*resp.Location))

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmVirtualNetworkGatewayConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient)
	vnetGatewayConnectionsClient := client.vnetGatewayConnectionsClient

	resGroup, name, err := resourceGroupAndVirtualNetworkGatewayConnectionFromId(d.Id())
	if err != nil {
		return err
	}

	_, error := vnetGatewayConnectionsClient.Delete(resGroup, name, make(chan struct{}))
	err = <-error
	if err != nil {
		return errwrap.Wrapf("Error Deleting VirtualNetworkGatewayConnection {{err}}", err)
	}

	log.Printf("[DEBUG] Waiting for AzureRM Virtual Network Gateway Connection %s to be removed", name)
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Accepted", "Deleting"},
		Target:  []string{"NotFound"},
		Refresh: virtualNetworkGatewayConnectionStateRefreshFunc(client, resGroup, name, true),
		Timeout: 15 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for AzureRM Virtual Network Gateway Connection %s to be removed: %+v", name, err)
	}

	d.SetId("")
	return nil
}

func virtualNetworkGatewayConnectionStateRefreshFunc(client *ArmClient, resourceGroupName string, virtualNetworkGatewayConnection string, withNotFound bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.vnetGatewayConnectionsClient.Get(resourceGroupName, virtualNetworkGatewayConnection)
		if err != nil {
			if withNotFound && responseWasNotFound(resp.Response) {
				return resp, "NotFound", nil
			}
			return nil, "", fmt.Errorf("Error making Read request on AzureRM Virtual Network Gateway Connection %s: %+v", virtualNetworkGatewayConnection, err)
		}

		return resp, *resp.VirtualNetworkGatewayConnectionPropertiesFormat.ProvisioningState, nil
	}
}

func getArmVirtualNetworkGatewayConnectionProperties(d *schema.ResourceData) (*network.VirtualNetworkGatewayConnectionPropertiesFormat, error) {
	connectionType := network.VirtualNetworkGatewayConnectionType(d.Get("type").(string))
	enableBgp := d.Get("enable_bgp").(bool)

	props := &network.VirtualNetworkGatewayConnectionPropertiesFormat{
		ConnectionType: connectionType,
		EnableBgp:      &enableBgp,
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

	if v, ok := d.GetOk("routing_weight"); ok {
		routingWeight := int32(v.(int))
		props.RoutingWeight = &routingWeight
	}

	if v, ok := d.GetOk("shared_key"); ok {
		sharedKey := v.(string)
		props.SharedKey = &sharedKey
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
