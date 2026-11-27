// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package nat_gateway

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/natgateways"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceNatGatewayDelete(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).NatGateways
	ctx, cancel := timeouts.ForDelete(account.StopContext(meta), d)
	defer cancel()

	id, err := natgateways.ParseNatGatewayID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.NatGatewayName, natGatewayResourceName)
	defer locks.UnlockByName(id.NatGatewayName, natGatewayResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
