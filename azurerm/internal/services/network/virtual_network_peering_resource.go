package network

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

// peerMutex is used to prevent multiple Peering resources being created, updated
// or deleted at the same time
var peerMutex = &sync.Mutex{}

func resourceVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualNetworkPeeringCreateUpdate,
		Read:   resourceVirtualNetworkPeeringRead,
		Update: resourceVirtualNetworkPeeringCreateUpdate,
		Delete: resourceVirtualNetworkPeeringDelete,
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_network_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"allow_virtual_network_access": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"allow_forwarded_traffic": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"allow_gateway_transit": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"use_remote_gateways": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceVirtualNetworkPeeringCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM virtual network peering creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resGroup, vnetName, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Peering %q (Virtual Network %q / Resource Group %q): %s", name, vnetName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_virtual_network_peering", *existing.ID)
		}
	}

	peer := network.VirtualNetworkPeering{
		Name:                                  &name,
		VirtualNetworkPeeringPropertiesFormat: getVirtualNetworkPeeringProperties(d),
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	if err := pluginsdk.Retry(300*time.Second, retryVnetPeeringsClientCreateUpdate(d, resGroup, vnetName, name, peer, meta)); err != nil {
		return err
	}

	read, err := client.Get(ctx, resGroup, vnetName, name)
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Virtual Network Peering %q (resource group %q)", name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceVirtualNetworkPeeringRead(d, meta)
}

func resourceVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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

	// update appropriate values
	d.Set("resource_group_name", resGroup)
	d.Set("name", resp.Name)
	d.Set("virtual_network_name", vnetName)

	if peer := resp.VirtualNetworkPeeringPropertiesFormat; peer != nil {
		d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
		d.Set("use_remote_gateways", peer.UseRemoteGateways)
		if network := peer.RemoteVirtualNetwork; network != nil {
			d.Set("remote_virtual_network_id", network.ID)
		}
	}

	return nil
}

func resourceVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VnetPeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
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

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Virtual Network Peering %q (Network %q / RG %q): %+v", name, vnetName, resGroup, err)
	}

	return err
}

func getVirtualNetworkPeeringProperties(d *pluginsdk.ResourceData) *network.VirtualNetworkPeeringPropertiesFormat {
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

//lintignore:R006
func retryVnetPeeringsClientCreateUpdate(d *pluginsdk.ResourceData, resGroup string, vnetName string, name string, peer network.VirtualNetworkPeering, meta interface{}) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		vnetPeeringsClient := meta.(*clients.Client).Network.VnetPeeringsClient
		ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		future, err := vnetPeeringsClient.CreateOrUpdate(ctx, resGroup, vnetName, name, peer)
		if err != nil {
			if utils.ResponseErrorIsRetryable(err) {
				return pluginsdk.RetryableError(err)
			} else if future.Response().StatusCode == 400 && strings.Contains(err.Error(), "ReferencedResourceNotProvisioned") {
				// Resource is not yet ready, this may be the case if the Vnet was just created or another peering was just initiated.
				return pluginsdk.RetryableError(err)
			}

			return pluginsdk.NonRetryableError(err)
		}

		if err = future.WaitForCompletionRef(ctx, vnetPeeringsClient.Client); err != nil {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	}
}
