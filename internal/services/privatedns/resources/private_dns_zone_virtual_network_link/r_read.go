// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package private_dns_zone_virtual_network_link

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/virtualnetworklinks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	privatednsClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePrivateDnsZoneVirtualNetworkLinkRead(d *pluginsdk.ResourceData, meta any) error {
	client := privatednsClient.From(meta).VirtualNetworkLinksClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := virtualnetworklinks.ParseVirtualNetworkLinkID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}
	return resourcePrivateDnsZoneVirtualNetworkLinkFlatten(d, id, resp.Model)
}
