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

func resourceVirtualHubBgpConnection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualHubBgpConnectionCreate,
		Read:   resourceVirtualHubBgpConnectionRead,
		Update: resourceVirtualHubBgpConnectionUpdate,
		Delete: resourceVirtualHubBgpConnectionDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
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

			"virtual_hub_id": {
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

			"virtual_network_connection_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: virtualwans.ValidateHubVirtualNetworkConnectionID,
			},
		},
	}
}

func resourceVirtualHubBgpConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id := commonids.NewVirtualHubBGPConnectionID(virtHubId.SubscriptionId, virtHubId.ResourceGroupName, virtHubId.VirtualHubName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.VirtualHubBgpConnectionGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_bgp_connection", id.ID())
		}
	}

	parameters := virtualwans.BgpConnection{
		Name: pointer.To(d.Get("name").(string)),
		Properties: &virtualwans.BgpConnectionProperties{
			PeerAsn: pointer.To(int64(d.Get("peer_asn").(int))),
			PeerIP:  pointer.To(d.Get("peer_ip").(string)),
		},
	}

	if v, ok := d.GetOk("virtual_network_connection_id"); ok {
		parameters.Properties.HubVirtualNetworkConnection = &virtualwans.SubResource{
			Id: pointer.To(v.(string)),
		}
	}

	if err := client.VirtualHubBgpConnectionCreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubBgpConnectionRead(d, meta)
}

func resourceVirtualHubBgpConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := virtualwans.ParseVirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.VirtualHubName, virtualHubResourceName)

	id, err := commonids.ParseVirtualHubBGPConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.VirtualHubBgpConnectionGet(ctx, *id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	if d.HasChange("virtual_network_connection_id") {
		if v, ok := d.GetOk("virtual_network_connection_id"); ok {
			existing.Model.Properties.HubVirtualNetworkConnection = &virtualwans.SubResource{
				Id: pointer.To(v.(string)),
			}
		} else {
			existing.Model.Properties.HubVirtualNetworkConnection = &virtualwans.SubResource{
				Id: nil,
			}
		}
	}

	if err := client.VirtualHubBgpConnectionCreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubBgpConnectionRead(d, meta)
}

func resourceVirtualHubBgpConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
			log.Printf("[INFO] network %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ConnectionName)
	d.Set("virtual_hub_id", virtualwans.NewVirtualHubID(id.SubscriptionId, id.ResourceGroupName, id.HubName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("peer_asn", props.PeerAsn)
			d.Set("peer_ip", props.PeerIP)
			if v := props.HubVirtualNetworkConnection; v != nil {
				d.Set("virtual_network_connection_id", v.Id)
			}
		}
	}

	return nil
}

func resourceVirtualHubBgpConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseVirtualHubBGPConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.HubName, virtualHubResourceName)
	defer locks.UnlockByName(id.HubName, virtualHubResourceName)

	if err := client.VirtualHubBgpConnectionDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
