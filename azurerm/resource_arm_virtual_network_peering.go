package azurerm

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-04-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// TODO: switch over to using the main mutex
// peerMutex is used to prevent multiple Peering resources being created, updated
// or deleted at the same time
var peerMutex = &sync.Mutex{}

func resourceArmVirtualNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkPeeringCreateUpdate,
		Read:   resourceArmVirtualNetworkPeeringRead,
		Update: resourceArmVirtualNetworkPeeringCreateUpdate,
		Delete: resourceArmVirtualNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(time.Minute * 30),
			Update: schema.DefaultTimeout(time.Minute * 30),
			Delete: schema.DefaultTimeout(time.Minute * 30),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"virtual_network_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_virtual_network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"allow_virtual_network_access": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"allow_forwarded_traffic": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"allow_gateway_transit": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"use_remote_gateways": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceArmVirtualNetworkPeeringCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetPeeringsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network peering creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		// first check if there's one in this subscription requiring import
		resp, err := client.Get(ctx, resGroup, vnetName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for the existence of Peering %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
			}
		}

		if resp.ID != nil {
			return tf.ImportAsExistsError("azurerm_virtual_network_peering", *resp.ID)
		}
	}

	peer := network.VirtualNetworkPeering{
		Name: &name,
		VirtualNetworkPeeringPropertiesFormat: getVirtualNetworkPeeringProperties(d),
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	future, err := client.CreateOrUpdate(ctx, resGroup, vnetName, name, peer)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Virtual Network Peering %q (Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(tf.TimeoutForCreateUpdate(d)))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for completion of Virtual Network Peering %q (Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, vnetName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Virtual Network Peering %q (resource group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmVirtualNetworkPeeringRead(d, meta)
}

func resourceArmVirtualNetworkPeeringRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetPeeringsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["virtualNetworkPeerings"]

	resp, err := client.Get(ctx, resGroup, vnetName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure virtual network peering %q: %+v", name, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("virtual_network_name", vnetName)

	if props := resp.VirtualNetworkPeeringPropertiesFormat; props != nil {
		d.Set("allow_virtual_network_access", props.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", props.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", props.AllowGatewayTransit)
		d.Set("use_remote_gateways", props.UseRemoteGateways)

		if network := props.RemoteVirtualNetwork; network != nil {
			d.Set("remote_virtual_network_id", network.ID)
		}
	}

	return nil
}

func resourceArmVirtualNetworkPeeringDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetPeeringsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["virtualNetworkPeerings"]

	peerMutex.Lock()
	defer peerMutex.Unlock()

	future, err := client.Delete(ctx, resGroup, vnetName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Virtual Network Peering %q (Network %q / RG %q): %+v", name, vnetName, resGroup, err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutDelete))
	defer cancel()
	err = future.WaitForCompletionRef(waitCtx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Peering %q (Network %q / RG %q): %+v", name, vnetName, resGroup, err)
	}

	return err
}

func getVirtualNetworkPeeringProperties(d *schema.ResourceData) *network.VirtualNetworkPeeringPropertiesFormat {
	allowVirtualNetworkAccess := d.Get("allow_virtual_network_access").(bool)
	allowForwardedTraffic := d.Get("allow_forwarded_traffic").(bool)
	allowGatewayTransit := d.Get("allow_gateway_transit").(bool)
	useRemoteGateways := d.Get("use_remote_gateways").(bool)
	remoteVirtualNetworkID := d.Get("remote_virtual_network_id").(string)

	return &network.VirtualNetworkPeeringPropertiesFormat{
		AllowVirtualNetworkAccess: utils.Bool(allowVirtualNetworkAccess),
		AllowForwardedTraffic:     utils.Bool(allowForwardedTraffic),
		AllowGatewayTransit:       utils.Bool(allowGatewayTransit),
		UseRemoteGateways:         utils.Bool(useRemoteGateways),
		RemoteVirtualNetwork: &network.SubResource{
			ID: utils.String(remoteVirtualNetworkID),
		},
	}
}
