package testdata

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/virtualwans"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualHubConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubConnectionCreateOrUpdate,
		Read:   resourceVirtualHubConnectionRead,
		Update: resourceVirtualHubConnectionCreateOrUpdate,
		Delete: resourceVirtualHubConnectionDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&virtualwans.HubVirtualNetworkConnectionId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&virtualwans.HubVirtualNetworkConnectionId{}),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtual_hub_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_virtual_network_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVirtualHubConnectionCreateOrUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtualHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	id := virtualwans.NewHubVirtualNetworkConnectionID(virtualHubId.SubscriptionId, virtualHubId.ResourceGroupName, virtualHubId.VirtualHubName, d.Get("name").(string))

	remoteVirtualNetworkId, err := commonids.ParseVirtualNetworkID(d.Get("remote_virtual_network_id").(string))
	if err != nil {
		return err
	}
	_ = remoteVirtualNetworkId

	connection := virtualwans.HubVirtualNetworkConnection{
		Name: pointer(id.HubVirtualNetworkConnectionName),
	}

	if err := client.HubVirtualNetworkConnectionsCreateOrUpdateThenPoll(ctx, id, connection); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubConnectionRead(d, meta)
}

func resourceVirtualHubConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.HubVirtualNetworkConnectionsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.HubVirtualNetworkConnectionName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("remote_virtual_network_id", props.RemoteVirtualNetwork.Id)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func resourceVirtualHubConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseHubVirtualNetworkConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.HubVirtualNetworkConnectionsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func pointer(s string) *string { return &s }
