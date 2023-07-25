// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
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
			_, err := parse.BgpConnectionID(id)
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
				ValidateFunc: validate.VirtualHubID,
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
				ValidateFunc: validate.HubVirtualNetworkConnectionID,
			},
		},
	}
}

func resourceVirtualHubBgpConnectionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubBgpConnectionClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := parse.VirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.Name, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.Name, virtualHubResourceName)

	id := parse.NewBgpConnectionID(virtHubId.SubscriptionId, virtHubId.ResourceGroup, virtHubId.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_virtual_hub_bgp_connection", id.ID())
		}
	}

	parameters := network.BgpConnection{
		Name: utils.String(d.Get("name").(string)),
		BgpConnectionProperties: &network.BgpConnectionProperties{
			PeerAsn: utils.Int64(int64(d.Get("peer_asn").(int))),
			PeerIP:  utils.String(d.Get("peer_ip").(string)),
		},
	}

	if v, ok := d.GetOk("virtual_network_connection_id"); ok {
		parameters.BgpConnectionProperties.HubVirtualNetworkConnection = &network.SubResource{
			ID: utils.String(v.(string)),
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creating/updating future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubBgpConnectionRead(d, meta)
}

func resourceVirtualHubBgpConnectionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubBgpConnectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	virtHubId, err := parse.VirtualHubID(d.Get("virtual_hub_id").(string))
	if err != nil {
		return err
	}

	locks.ByName(virtHubId.Name, virtualHubResourceName)
	defer locks.UnlockByName(virtHubId.Name, virtualHubResourceName)

	id, err := parse.BgpConnectionID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if d.HasChange("virtual_network_connection_id") {
		if v, ok := d.GetOk("virtual_network_connection_id"); ok {
			existing.BgpConnectionProperties.HubVirtualNetworkConnection = &network.SubResource{
				ID: utils.String(v.(string)),
			}
		} else {
			existing.BgpConnectionProperties.HubVirtualNetworkConnection = &network.SubResource{
				ID: nil,
			}
		}
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.VirtualHubName, id.Name, existing)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on updating future for %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualHubBgpConnectionRead(d, meta)
}

func resourceVirtualHubBgpConnectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubBgpConnectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BgpConnectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] network %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Virtual Hub Bgp Connection %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	d.Set("name", id.Name)
	d.Set("virtual_hub_id", parse.NewVirtualHubID(id.SubscriptionId, id.ResourceGroup, id.VirtualHubName).ID())

	if props := resp.BgpConnectionProperties; props != nil {
		d.Set("peer_asn", props.PeerAsn)
		d.Set("peer_ip", props.PeerIP)
		if v := props.HubVirtualNetworkConnection; v != nil {
			d.Set("virtual_network_connection_id", v.ID)
		}
	}

	return nil
}

func resourceVirtualHubBgpConnectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.VirtualHubBgpConnectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.BgpConnectionID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualHubName, virtualHubResourceName)
	defer locks.UnlockByName(id.VirtualHubName, virtualHubResourceName)

	future, err := client.Delete(ctx, id.ResourceGroup, id.VirtualHubName, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Virtual Hub Bgp Connection %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on deleting future for Virtual Hub Bgp Connection %q (Resource Group %q / Virtual Hub %q): %+v", id.Name, id.ResourceGroup, id.VirtualHubName, err)
	}

	return nil
}
