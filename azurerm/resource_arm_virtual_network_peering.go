package azurerm

import (
	"fmt"
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2017-09-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// peerMutex is used to prevent multiple Peering resources being created, updated
// or deleted at the same time
var peerMutex = &sync.Mutex{}

func resourceArmVirtualNetworkPeering() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmVirtualNetworkPeeringCreate,
		Read:   resourceArmVirtualNetworkPeeringRead,
		Update: resourceArmVirtualNetworkPeeringCreate,
		Delete: resourceArmVirtualNetworkPeeringDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceArmVirtualNetworkPeeringCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).vnetPeeringsClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network peering creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)

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

	err = future.WaitForCompletion(ctx, client.Client)
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

	peer := *resp.VirtualNetworkPeeringPropertiesFormat

	// update appropriate values
	d.Set("resource_group_name", resGroup)
	d.Set("name", resp.Name)
	d.Set("virtual_network_name", vnetName)
	d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
	d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
	d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
	d.Set("use_remote_gateways", peer.UseRemoteGateways)
	d.Set("remote_virtual_network_id", peer.RemoteVirtualNetwork.ID)

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

	err = future.WaitForCompletion(ctx, client.Client)
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
		AllowVirtualNetworkAccess: &allowVirtualNetworkAccess,
		AllowForwardedTraffic:     &allowForwardedTraffic,
		AllowGatewayTransit:       &allowGatewayTransit,
		UseRemoteGateways:         &useRemoteGateways,
		RemoteVirtualNetwork: &network.SubResource{
			ID: &remoteVirtualNetworkID,
		},
	}
}
