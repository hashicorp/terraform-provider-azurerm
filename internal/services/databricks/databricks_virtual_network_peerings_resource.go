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
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DatabricksVirtualNetworkPeeringsName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: vnetpeering.ValidateWorkspaceID,
			},

			"remote_address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,

				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.CIDR,
				},
			},

			"remote_virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.VirtualNetworkID,
			},

			"address_space_prefixes": {
				Type:     pluginsdk.TypeList,
				Computed: true,

				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"virtual_network_id": {
				Type:         pluginsdk.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: networkValidate.VirtualNetworkID,
			},

			"virtual_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"forwarded_traffic_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
				Optional: true,
			},

			"gateway_transit_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
				Optional: true,
			},

			"remote_gateways_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
				Optional: true,
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
	var id vnetpeering.VirtualNetworkPeeringId

	if d.IsNewResource() {
		// I need to include the workspace ID in the properties because I need the name
		// of the workspace to create the peerings ID
		workspaceId, err := vnetpeering.ParseWorkspaceID(d.Get("workspace_id").(string))
		if err != nil {
			return fmt.Errorf("unable to parse 'workspace_id': %+v", err)
		}

		id = vnetpeering.NewVirtualNetworkPeeringID(subscriptionId, d.Get("resource_group_name").(string), workspaceId.WorkspaceName, d.Get("name").(string))

		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Databricks %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_databricks_virtual_network_peering", id.ID())
		}
	} else {
		// For import and update, you need to parse the id else you will have an empty subscription,
		// resource group name and workspace name in the PUT call causing a 400 Bad Request...
		rawId, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
		if err != nil {
			return fmt.Errorf("unable to parse 'id': %+v", err)
		}

		id = *rawId
	}

	peer := vnetpeering.VirtualNetworkPeering{
		Name:       &id.VirtualNetworkPeeringName,
		Properties: expandDatabricksVirtualNetworkPeeringProperties(d),
	}

	// The RP always creates the same vNet ID for the Databricks internal vNet in the below format:
	// '/subscriptions/{subscription}/resourceGroups/{group1}/providers/Microsoft.Network/virtualNetworks/workers-vnet'
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/workers-vnet"
	virtualNetworkId := fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName)

	peer.Properties.DatabricksVirtualNetwork = &vnetpeering.VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork{
		Id: utils.String(virtualNetworkId),
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	if err := pluginsdk.Retry(300*time.Second, retryDatabricksVnetPeeringsClientCreateUpdate(d, id, peer, meta)); err != nil {
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

	// for the import scenario, once the peering is created I can generate the workspace ID from the peering ID
	// NOTE: this may not be right because I don't know if it is valid to create a peering with a
	// workspace across subscription boundaries?
	workspaceId := vnetpeering.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making Read request on Databricks %s: %+v", *id, err)
	}

	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("name", id.VirtualNetworkPeeringName)
	d.Set("workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		d.Set("virtual_network_access_enabled", model.Properties.AllowVirtualNetworkAccess)
		d.Set("forwarded_traffic_enabled", model.Properties.AllowForwardedTraffic)
		d.Set("gateway_transit_enabled", model.Properties.AllowGatewayTransit)
		d.Set("remote_gateways_enabled", model.Properties.UseRemoteGateways)

		if model.Properties.DatabricksAddressSpace != nil && model.Properties.DatabricksAddressSpace.AddressPrefixes != nil {
			d.Set("address_space_prefixes", model.Properties.DatabricksAddressSpace.AddressPrefixes)
		}

		if model.Properties.RemoteAddressSpace != nil && model.Properties.RemoteAddressSpace.AddressPrefixes != nil {
			d.Set("remote_address_space_prefixes", model.Properties.RemoteAddressSpace.AddressPrefixes)
		}

		if model.Properties.DatabricksVirtualNetwork != nil {
			if databricksNetwork := model.Properties.DatabricksVirtualNetwork.Id; databricksNetwork != nil {
				d.Set("virtual_network_id", databricksNetwork)
			}
		}

		if remoteNetwork := model.Properties.RemoteVirtualNetwork.Id; remoteNetwork != nil {
			d.Set("remote_virtual_network_id", remoteNetwork)
		}
	}

	return nil
}

func resourceDatabricksVirtualNetworkPeeringDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DataBricks.VnetPeeringsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := vnetpeering.ParseVirtualNetworkPeeringID(d.Id())
	if err != nil {
		return err
	}

	// check to see if it exist first, this maybe a no op because when you delete a peering
	// it auto deletes the corresponding peering in the other resource
	existing, err := client.Get(ctx, *id)
	if err != nil && response.WasNotFound(existing.HttpResponse) {
		return nil
	}

	peerMutex.Lock()
	defer peerMutex.Unlock()

	// this is a workaround for the 'failed to unmarshal response body: StatusCode=0
	// -- Original Error: json: cannot unmarshal string into Go value of type map[string]interface {}'
	// error
	result, err := client.Delete(ctx, *id)
	if err != nil && result.HttpResponse != nil {
		return fmt.Errorf("deleting Databricks %s: %+v", *id, err)
	}

	return nil
}

func expandDatabricksVirtualNetworkPeeringProperties(d *pluginsdk.ResourceData) vnetpeering.VirtualNetworkPeeringPropertiesFormat {
	allowForwardedTraffic := d.Get("forwarded_traffic_enabled").(bool)
	allowGatewayTransit := d.Get("gateway_transit_enabled").(bool)
	allowVirtualNetworkAccess := d.Get("virtual_network_access_enabled").(bool)
	useRemoteGateways := d.Get("remote_gateways_enabled").(bool)
	remoteVirtualNetwork := d.Get("remote_virtual_network_id").(string)
	databricksAddressSpace := utils.ExpandStringSlice(d.Get("address_space_prefixes").([]interface{}))
	remoteAddressSpace := utils.ExpandStringSlice(d.Get("remote_address_space_prefixes").([]interface{}))

	return vnetpeering.VirtualNetworkPeeringPropertiesFormat{
		AllowForwardedTraffic:     &allowForwardedTraffic,
		AllowGatewayTransit:       &allowGatewayTransit,
		AllowVirtualNetworkAccess: &allowVirtualNetworkAccess,
		DatabricksAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: databricksAddressSpace,
		},
		RemoteAddressSpace: &vnetpeering.AddressSpace{
			AddressPrefixes: remoteAddressSpace,
		},
		RemoteVirtualNetwork: vnetpeering.VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork{
			Id: utils.String(remoteVirtualNetwork),
		},
		UseRemoteGateways: &useRemoteGateways,
	}
}

func retryDatabricksVnetPeeringsClientCreateUpdate(d *pluginsdk.ResourceData, id vnetpeering.VirtualNetworkPeeringId, peer vnetpeering.VirtualNetworkPeering, meta interface{}) func() *pluginsdk.RetryError {
	return func() *pluginsdk.RetryError {
		vnetPeeringsClient := meta.(*clients.Client).DataBricks.VnetPeeringsClient
		ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
		defer cancel()

		future, err := vnetPeeringsClient.CreateOrUpdate(ctx, id, peer)
		if err != nil {
			if utils.ResponseErrorIsRetryable(err) {
				return pluginsdk.RetryableError(err)
			} else if future.HttpResponse != nil && future.HttpResponse.StatusCode == 400 && strings.Contains(err.Error(), "ReferencedResourceNotProvisioned") {
				// Resource is not yet ready, this may be the case if the Vnet was just created or another peering was just initiated.
				return pluginsdk.RetryableError(err)
			}

			return pluginsdk.NonRetryableError(err)
		}

		if err = future.Poller.PollUntilDone(); err != nil {
			return pluginsdk.NonRetryableError(err)
		}

		return nil
	}
}
