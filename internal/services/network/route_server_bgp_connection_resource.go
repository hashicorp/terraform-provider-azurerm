// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServerBgpConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceRouteServerBgpConnectionCreate,
		Read:   resourceRouteServerBgpConnectionRead,
		Delete: resourceRouteServerBgpConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseVirtualHubBGPConnectionID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"route_server_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: virtualwans.ValidateVirtualHubID,
			},

			"peer_asn": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"peer_ip": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
			},
		},
	}
}

func resourceRouteServerBgpConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	routerServerId, err := virtualwans.ParseVirtualHubID(d.Get("route_server_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(routerServerId.VirtualHubName, "azurerm_route_server")
	defer locks.UnlockByName(routerServerId.VirtualHubName, "azurerm_route_server")

	id := commonids.NewVirtualHubBGPConnectionID(routerServerId.SubscriptionId, routerServerId.ResourceGroupName, routerServerId.VirtualHubName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.VirtualHubBgpConnectionGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_route_server_bgp_connection", id.ID())
		}
	}
	parameters := virtualwans.BgpConnection{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.BgpConnectionProperties{
			PeerAsn: pointer.To(int64(d.Get("peer_asn").(int))),
			PeerIP:  pointer.To(d.Get("peer_ip").(string)),
		},
	}

	if err := client.VirtualHubBgpConnectionCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceRouteServerBgpConnectionRead(d, meta)
}

func resourceRouteServerBgpConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubBGPConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.VirtualHubBgpConnectionGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exists", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("route_server_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.HubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.PeerAsn != nil {
				d.Set("peer_asn", props.PeerAsn)
			}
			if props.PeerIP != nil {
				d.Set("peer_ip", props.PeerIP)
			}
		}
	}
	return nil
}

func resourceRouteServerBgpConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubBGPConnectionID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VirtualHubBgpConnectionDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s： %+v", id, err)
	}

	return nil
}
