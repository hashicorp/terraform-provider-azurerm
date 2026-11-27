// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceRouteServerDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).VirtualWan.VirtualWANs
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualwans.ParseVirtualHubID(d.Id())
	if err != nil {
		return err
	}

	ipConfig, err := client.VirtualHubIPConfigurationList(ctx, *id)
	if err != nil {
		return fmt.Errorf("retrieving IP Config for %s: %+v", id, err)
	}

	if ipConfig.Model != nil {
		var ipName string
		for _, config := range *ipConfig.Model {
			if config.Name != nil {
				ipName = *config.Name
			}
		}

		ipConfigId := commonids.NewVirtualHubIPConfigurationID(id.SubscriptionId, id.ResourceGroupName, id.VirtualHubName, ipName)

		if err := deleteRouteServerIpConfiguration(ctx, client, ipConfigId); err != nil {
			return err
		}
	}

	if err := client.VirtualHubsDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
