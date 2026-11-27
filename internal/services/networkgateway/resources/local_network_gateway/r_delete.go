// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package local_network_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/localnetworkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkgatewayClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLocalNetworkGatewayDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkgatewayClient.From(meta).LocalNetworkGateways
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := localnetworkgateways.ParseLocalNetworkGatewayID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
