// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_dns_servers

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkDnsServersCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	vnetId, err := commonids.ParseVirtualNetworkID(d.Get("virtual_network_id").(string))
	if err != nil {
		return err
	}

	// This is a virtual resource so the last segment is hardcoded
	id := parse.NewVirtualNetworkDnsServersID(vnetId.SubscriptionId, vnetId.ResourceGroupName, vnetId.VirtualNetworkName, "default")

	locks.ByID(vnetId.ID())
	defer locks.UnlockByID(vnetId.ID())

	vnet, err := client.Get(ctx, *vnetId, virtualnetworks.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(vnet.HttpResponse) {
			return fmt.Errorf("%s could not be found: %+v", vnetId, err)
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
		vnet.Model.Properties.DhcpOptions = &virtualnetworks.DhcpOptions{}
	}

	vnet.Model.Properties.DhcpOptions.DnsServers = pluginsdk.ExpandStringSlice(d.Get("dns_servers").([]any))

	// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
	if err := client.CreateOrUpdateThenPoll(ctx, *vnetId, *vnet.Model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	timeout, _ := ctx.Deadline()

	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(virtualnetworks.ProvisioningStateUpdating)},
		Target:     []string{string(virtualnetworks.ProvisioningStateSucceeded)},
		Refresh:    helpers.VirtualNetworkProvisioningStateRefreshFunc(ctx, client, *vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	return resourceVirtualNetworkDnsServersRead(d, meta)
}
