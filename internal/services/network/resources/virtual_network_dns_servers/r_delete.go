// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_dns_servers

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkDnsServersDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.VirtualNetworkDnsServersID(d.Id())
	if err != nil {
		return err
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroup, id.VirtualNetworkName)

	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	vnet, err := client.Get(ctx, vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vnet.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing %s from state", vnetId.ID(), id)
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", vnetId, err)
	}

	if vnet.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", vnetId)
	}
	if vnet.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", vnetId)
	}

	if vnet.Model.Properties.DhcpOptions == nil {
		log.Printf("[INFO] dhcpOptions for %s was nil, dnsServers already deleted - removing %s from state", vnetId.ID(), id)
		return nil
	}

	vnet.Model.Properties.DhcpOptions.DnsServers = pluginsdk.ExpandStringSlice(make([]any, 0))

	if err := client.CreateOrUpdateThenPoll(ctx, vnetId, *vnet.Model); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
