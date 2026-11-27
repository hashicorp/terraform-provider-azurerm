// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network_gateway_nat_rule

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	networkgatewayClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/networkgateway/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualNetworkGatewayNatRuleDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkgatewayClient.From(meta).NetworkGateways
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := networkgateways.ParseVirtualNetworkGatewayNatRuleID(d.Id())
	if err != nil {
		return err
	}

	if err := client.VirtualNetworkGatewayNatRulesDeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}
