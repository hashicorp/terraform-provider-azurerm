// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server_bgp_connection

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServerBgpConnectionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
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

	return pluginsdk.SetResourceIdentityData(d, id)
}
