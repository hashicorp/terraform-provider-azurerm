package databricks

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/databricks/2023-02-01/vnetpeering"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// peerMutex is used to prevent multiple Peering resources being created, updated
// or deleted at the same time
var peerMutex = &sync.Mutex{}

func resourceDatabricksVirtualNetworkPeering() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDatabricksVirtualNetworkPeeringCreateUpdate,
		Read:   resourceDatabricksVirtualNetworkPeeringRead,
		Update: resourceDatabricksVirtualNetworkPeeringCreateUpdate,
		Delete: resourceDatabricksVirtualNetworkPeeringDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := vnetpeering.ParseVirtualNetworkPeeringID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
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

			"remote_address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"remote_virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"use_remote_gateways": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDatabricksVirtualNetworkPeeringCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM databricks virtual network peering creation.")
	id := vnetpeering.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), d.Get("workspace_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_virtual_network_peering", id.ID())
		}
	}

	peer := vnetpeering.VirtualNetworkPeering{
		Name:       &id.VirtualNetworkPeeringName,
		Properties: expandDatabricksVirtualNetworkPeeringProperties(d),
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	if err := pluginsdk.Retry(300*time.Second, retryVnetPeeringsClientCreateUpdate(d, id, peer, meta)); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceDatabricksVirtualNetworkPeeringRead(d, meta)
}

func resourceDatabricksVirtualNetworkPeeringRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("name", id.VirtualNetworkPeeringName)

	if peer := resp.Model.Properties; peer != nil {
		d.Set("allow_virtual_network_access", peer.AllowVirtualNetworkAccess)
		d.Set("allow_forwarded_traffic", peer.AllowForwardedTraffic)
		d.Set("allow_gateway_transit", peer.AllowGatewayTransit)
		d.Set("use_remote_gateways", peer.UseRemoteGateways)
		d.Set("virtual_network_id", peer.RemoteVirtualNetwork.Id)
		if network := peer.RemoteVirtualNetwork.Id; network != nil {
			d.Set("remote_virtual_network_id", network)
		}
	}

	return nil
}

func resourceDatabricksVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DatabricksVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualNetworkName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
	}

	return err
}

func expandDatabricksVirtualNetworkPeeringProperties(d *pluginsdk.ResourceData) vnetpeering.VirtualNetworkPeeringPropertiesFormat {
	allowForwardedTraffic := d.Get("allow_forwarded_traffic").(bool)
	allowGatewayTransit := d.Get("allow_gateway_transit").(bool)
	allowVirtualNetworkAccess := d.Get("allow_virtual_network_access").(bool)
	databricksAddressSpace := utils.ExpandStringSlice(d.Get("address_space_prefixes").([]interface{}))
	databricksVirtualNetwork := d.Get("virtual_network_id").(string)
	remoteAddressSpace := utils.ExpandStringSlice(d.Get("remote_address_space_prefixes").([]interface{}))
	remoteVirtualNetwork := d.Get("remote_virtual_network_id").(string)
	useRemoteGateways := d.Get("use_remote_gateways").(bool)

	return vnetpeering.VirtualNetworkPeeringPropertiesFormat{
		AllowForwardedTraffic:     &allowForwardedTraffic,
		AllowGatewayTransit:       &allowGatewayTransit,
		AllowVirtualNetworkAccess: &allowVirtualNetworkAccess,
		DatabricksAddressSpace:    &vnetpeering.AddressSpace{databricksAddressSpace},
		DatabricksVirtualNetwork: &vnetpeering.VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork{
			Id: utils.String(databricksVirtualNetwork),
		},
		RemoteAddressSpace: &vnetpeering.AddressSpace{remoteAddressSpace},
		RemoteVirtualNetwork: vnetpeering.VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork{
			Id: utils.String(remoteVirtualNetwork),
		},
		UseRemoteGateways: &useRemoteGateways,
	}
}

func retryVnetPeeringsClientCreateUpdate(d *pluginsdk.ResourceData, id vnetpeering.VirtualNetworkPeeringId, peer vnetpeering.VirtualNetworkPeering, meta interface{}) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		vnetPeeringsClient := meta.(*clients.Client).DataBricks.VnetPeeringsClient
		ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		future, err := vnetPeeringsClient.CreateOrUpdate(ctx, id, peer)
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
